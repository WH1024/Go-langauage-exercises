package controller

import (
	"encoding/json"
	"fmt"
	"go-spi/db"
	"go-spi/utils"
	"log"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/kataras/iris/v12"
	m "github.com/murphy214/mercantile"
	geojson "github.com/paulmach/go.geojson"
	"github.com/xiaolingis/geojsonvt"
)

// MeshQuery
// 拿到前端请求时间段，存到common里，后面直接拿；返回前端一个由时间段生成的md5
func MeshQuery(ctx iris.Context) {
	timeMD5 := ctx.URLParam("timeKey")
	if timeMD5 == "" {
		// 获取时间段
		// startTime, err := ctx.URLParamFloat64("startTime")
		// endTime, err := ctx.URLParamFloat64("endTime")
		var (
			err                error
			startTime, endTime int64
			wg                 sync.WaitGroup
			// rw                 sync.RWMutex
		)
		startTime, err = ctx.URLParamInt64("startTime")
		endTime, err = ctx.URLParamInt64("endTime")
		fmt.Println("时间段: ", startTime, "-", endTime)
		if err != nil || startTime >= endTime {
			ctx.StopWithStatus(iris.StatusBadRequest)
			return
		}
		startTime = startTime - 3600*8 // utc时间-8个时区=东八区
		endTime = endTime - 3600*8
		// md5Str := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%d-%d", startTime, endTime))))
		// // 设置数据返回类型
		// ctx.Header("Content-Type", "application/x-protobuf")
		// 从数据库查
		clickhousePtr := db.GetDB()
		defer clickhousePtr.Close()
		// // 存入pika
		// redisStart := time.Now()
		redisPtr := db.GetRedisDbInstance()
		defer redisPtr.Close()
		var d []*utils.Data
		// // err = clickhousePtr.Select(&d, fmt.Sprintf("select avg(lat) as lat,avg(lng) as lng,avg(ad2) as ad2,avg(ad3) as ad3,avg(ad4) as ad4,avg(ad7) as ad7,avg(ad16) as ad16 from AirData where ts between %f and %f group by geoid", startTime, endTime))
		sql := fmt.Sprintf("select avg(lat) as lat,avg(lng) as lng,avg(ad_2) as ad2,avg(ad_3) as ad3,avg(ad_4) as ad4,avg(ad_7) as ad7,avg(ad_16) as ad16 from AirData where ts between '%d' and '%d' group by geoid", startTime, endTime)
		// sql := fmt.Sprintf("select lat,lng,ad2,ad3,ad4,ad7,ad16 FROM AirData_view_1 where ts between '%d' and '%d'", startTime, endTime)
		fmt.Println("sql语句: ", sql)
		selectStart := time.Now()
		err = clickhousePtr.Select(&d, sql)
		if err != nil {
			log.Printf("数据库查询异常: %v", err)
			ctx.StopWithStatus(iris.StatusInternalServerError)
			return
		}
		fmt.Println("数据库查询耗时: ", time.Since(selectStart))
		createTime := time.Now()
		fcs := utils.CreateGeojson(d)
		fmt.Println("创建瓦片耗时: ", time.Since(createTime))
		rdPtr := db.GetRedisDbInstance()
		defer rdPtr.Close()
		tileTime := time.Now()
		// 方法
		var dataMap sync.Map
		splitFunc := func(features []*geojson.Feature, tileID m.TileID, _wg *sync.WaitGroup) {
			// options := geojsonvt.NewConfig()
			// tile := geojsonvt.TileFromGeoJSON(features, tileID, options)
			// tile.Options.LayerName = "points"
			id := utils.ToID(uint8(tileID.Z), uint32(tileID.X), uint32(tileID.Y))
			data, err := json.Marshal(features)
			if err == nil {
				dataMap.Store(id, data)
			}
			// dataMap.Store(id, tile.Marshal())
			_wg.Done()
		}

		for xyz, fc := range fcs {
			wg.Add(1)
			go splitFunc(fc.Features, xyz, &wg)
		}
		wg.Wait()
		fmt.Println("瓦片切割耗时: ", time.Since(tileTime))
		sendTime := time.Now()
		dataMap.Range(func(key, value interface{}) bool {
			rdPtr.Send("HSET", "19kefei", key, value)
			return true
		})
		fmt.Println("send花费时间: ", time.Since(sendTime))
		redisTime := time.Now()
		err = rdPtr.Flush()
		fmt.Println("存入pika耗时: ", time.Since(redisTime))
		if err != nil {
			fmt.Println("redis存入失败: ", err.Error())
			ctx.StopWithError(iris.StatusInternalServerError, err)
			return
		}

		ctx.StopWithJSON(iris.StatusOK, iris.Map{"success": true})
	}
}

func Tiles(ctx iris.Context) {
	// var wg sync.WaitGroup
	tileBase := "/tiles/"
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Content-Type", "application/x-protobuf")
	// ctx.Request().Header.Add("Access-Control-Allow-Origin", "*")
	tilePart := ctx.Request().URL.Path[len(tileBase):]
	log.Printf("url: %s", tilePart)
	xyz, err := utils.TilePathToXYZ(tilePart)
	if err != nil {
		fmt.Printf("TilePathToXYZ Error: %v", err)
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}
	// 层级过滤
	if xyz.Z < 9 || xyz.Z > 18 {
		ctx.StopWithStatus(iris.StatusBadRequest)
		return
	}

	// 头部信息
	// refferer := ctx.GetReferrer()
	// switch refferer.Type {
	// case iris.ReferrerSearch:
	// 	fmt.Println("ReferrerSearch: ", iris.ReferrerSearch)
	// case iris.ReferrerSocial:
	// 	fmt.Println("ReferrerSocial: ", iris.ReferrerSocial)
	// case iris.ReferrerIndirect:
	// 	fmt.Println("ReferrerIndirect: ", iris.ReferrerIndirect)
	// case iris.ReferrerEmail:
	// 	fmt.Println("ReferrerEmail: ", iris.ReferrerEmail)
	// }

	// fmt.Printf("对象: %#v\n", refferer)

	// header1 := ctx.Request().Header.Get("Last-Modified")
	// header := ctx.Request().Header.Get("If-Modified-Since")
	// header2 := ctx.Request().Header.Get("If-None-Match")
	// if len(header) > 0 || len(header2) > 0 {
	// 	fmt.Println("有缓存...", header, header2)
	// 	ctx.StopWithStatus(iris.StatusNotModified)
	// 	return
	// }
	targetID := utils.ToID(uint8(xyz.Z), xyz.X, xyz.Y)
	rdPtr := db.GetRedisDbInstance()
	defer rdPtr.Close()

	queryTime := time.Now()
	data, err := redis.Bytes(rdPtr.Do("HGET", "19kefei", targetID))
	if err == nil {
		// ctx.Header("Last-Modified", time.Now().Format(iris.DefaultConfiguration().TimeFormat))
		// ctx.Header("ETag", targetID)
		// ctx.Write(data)
		var features []*geojson.Feature
		err = json.Unmarshal(data, &features)
		if err == nil {
			options := geojsonvt.NewConfig()
			tile := geojsonvt.TileFromGeoJSON(features, m.TileID{X: int64(xyz.X), Y: int64(xyz.Y), Z: uint64(xyz.Z)}, options)
			tile.Options.LayerName = "points"
			// ctx.Header("Last-Modified", time.Now().Format(iris.DefaultConfiguration().TimeFormat))
			// ctx.Header("ETag", targetID)
			ctx.Write(tile.Marshal())
			ctx.StopWithStatus(iris.StatusOK)
			fmt.Println("有数据，返回，耗时: ", time.Since(queryTime))
			return
		}
		fmt.Println("有数据，返回，耗时: ", time.Since(queryTime))
		return
	}
	// fmt.Println("走道最后没数据")
	ctx.StopWithStatus(iris.StatusInternalServerError)
}
