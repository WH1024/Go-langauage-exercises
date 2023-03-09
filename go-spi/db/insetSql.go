package db

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"go-spi/common"
	"go-spi/utils"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mmcloughlin/geohash"
)

func InsertData() {
	var wg sync.WaitGroup
	start := time.Now()
	filenames := make([]string, 0, 20)
	for i := 15; i <= 30; i++ {
		filenames = append(filenames, fmt.Sprintf("2022-06-%d.csv", i))
	}
	for i := 1; i <= 4; i++ {
		filenames = append(filenames, fmt.Sprintf("2022-07-0%d.csv", i))
	}

	dbPtr := GetDB() // 拿到数据库实例
	defer dbPtr.Close()
	for _, filename := range filenames {
		wg.Add(1)
		go func(file string) {
			content, err := ioutil.ReadFile(common.RootPath + "/db/data/" + file)
			// content, err := ioutil.ReadFile(common.RootPath + "/db/test.csv")
			if err != nil {
				log.Fatal(err)
			}
			r := csv.NewReader(strings.NewReader(string(content[:])))
			records, err := r.ReadAll()
			if err != nil {
				log.Fatal(err)
			}

			num := 0
			var b bytes.Buffer
			for _, record := range records[1:] {
				ts := strings.Split(record[0], ".")[0] + "'" // 时间取.符号前面的，后面补一个单引号"'"
				lat_, _ := strconv.ParseFloat(record[3], 64)
				lng_, _ := strconv.ParseFloat(record[2], 64)
				// 这里是拿到数据点后，通过其经纬度计算其生成的网格经纬度的左上角的点
				_, lat_1, lng_1 := utils.GetMeshPosition(lat_, lng_)
				geoid := geohash.Encode(lat_1, lng_1)

				lat := strconv.FormatFloat(lat_1, 'f', 6, 64) //record[3]
				lng := strconv.FormatFloat(lng_1, 'f', 6, 64) //record[2]
				ad_2 := record[4]
				ad_3 := record[5]
				ad_4 := record[6]
				ad_5 := record[7]
				ad_6 := record[8]
				// 循环生成批量插入语句
				if num == 0 {
					b.WriteString(fmt.Sprintf("INSERT INTO AirData (ts,geoid,lat,lng,ad_2,ad_3,ad_4,ad_7,ad_16) VALUES (%s,'%s',%s,%s,%s,%s,%s,%s,%s)", ts, geoid, lat, lng, ad_2, ad_3, ad_4, ad_5, ad_6))
				} else {
					b.WriteString(fmt.Sprintf(",(%s,'%s',%s,%s,%s,%s,%s,%s,%s)", ts, geoid, lat, lng, ad_2, ad_3, ad_4, ad_5, ad_6))
				}
				num++
				//fmt.Println(fmt.Sprintf("INSERT INTO AirData (ts,geoid,lat,lng,ad_2,ad_3,ad_4,ad_7,ad_16) VALUES (%s,%d,%s,%s,%s,%s,%s,%s,%s)", ts, geoid, lat, lng, ad_2, ad_3, ad_4, ad_5, ad_6))
				//_, err := dbPtr.Exec(fmt.Sprintf("INSERT INTO AirData (ts,geoid,lat,lng,ad_2,ad_3,ad_4,ad_7,ad_16) VALUES (%s,%d,%s,%s,%s,%s,%s,%s,%s)", ts, geoid, lat, lng, ad_2, ad_3, ad_4, ad_5, ad_6))
				if num >= 1000 {
					// fmt.Println("send")
					_, err := dbPtr.Exec(b.String())
					if err != nil {
						log.Println("数据库插入失败: ", err.Error())
					}
					time.Sleep(time.Millisecond * 500) // 防止插入过快导致io timeout
					b.Reset()
					num = 0
				}
			}
			if b.Len() > 0 {
				_, err := dbPtr.Exec(b.String())
				if err != nil {
					log.Println("数据库插入失败: ", err.Error())
				}
				b.Reset()
			}
			wg.Done()
		}(filename)
	}
	wg.Wait()
	fmt.Println("耗时: ", time.Since(start))
}
