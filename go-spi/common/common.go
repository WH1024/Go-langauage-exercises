package common

import (
	"errors"
	"os"
	"regexp"

	geojson "github.com/paulmach/go.geojson"
	"github.com/xiaolingis/geojsonvt"
)

var (
	dir, _   = os.Getwd()
	RootPath string // 暴露给外部使用的根路径
	features []*geojson.Feature
	tileMap  map[string]map[string]*geojsonvt.Tile
	// 请求网格数据的时间段
	startTime float64
	endTime   float64
)

func init() {
	rePath := `(.*?)go-spi`
	re := regexp.MustCompile(rePath)
	RootPath = re.FindAllStringSubmatch(dir, -1)[0][0]
}

func GetTime() (float64, float64) {
	return startTime, endTime
}

func SetTime(start, end float64) {
	startTime = start
	endTime = end
}

func GetFeatures() ([]*geojson.Feature, error) {
	if len(features) < 1 {
		return nil, errors.New("features is nil")
	}
	return features, nil
}

func SetFeatures(f []*geojson.Feature) {
	features = f
}

// 存储时间段对应的瓦片
func SetTile(account, timeKey string, tile *geojsonvt.Tile) {
	tileMap[account][timeKey] = tile
}

func GetTile(account, timeKey string) (tile *geojsonvt.Tile, ok bool) {
	tile, ok = tileMap[account][timeKey]
	return
}
