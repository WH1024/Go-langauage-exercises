package utils

import (
	"errors"
	"math"
	"regexp"
	"strconv"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

type LngLat struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

// TileID represents the id of the tile.
type TileID struct {
	X       uint32
	Y       uint32
	Z       uint32
	Account string
	Key     string
}

type Data struct {
	Geoid string   `db:"geoid"`
	Lat   float64  `db:"lat" json:"lat"`
	Lng   float64  `db:"lng" json:"lng"`
	Ad2   *float64 //sql.NullFloat64
	Ad3   *float64
	Ad4   *float64
	Ad7   *float64
	Ad16  *float64
}

const (
	EarthRadiusKm = 6378.137 // WGS-84
	_rangeMeters  = 50.0     //m
)

const (
	gMinLat   = -85.05112878
	gMaxLat   = 85.05112878
	gMinLon   = -180.0
	gMaxLon   = 180.0
	gTileSize = 256
)

// Fabs
// 取绝对值
func Fabs(v float64) float64 {
	if v > 0 {
		return v
	} else {
		return -v
	}
}

// CreateTile
// 生成瓦片的protobuf
func CreateTile(d []*Data) map[TileID]*geojson.FeatureCollection {
	mtile := make(map[TileID]*geojson.FeatureCollection)
	mesh_width := 50.0 * math.Abs(113.545749-115.231401) / (GetDistance(LngLat{Lng: 113.545749, Lat: 22.966396}, LngLat{Lng: 115.231401, Lat: 22.966396}) * 1000.0)
	mesh_height := 50.0 * math.Abs(22.966396-22.533389) / (GetDistance(LngLat{Lng: 113.545749, Lat: 22.966396}, LngLat{Lng: 113.545749, Lat: 22.533389}) * 1000.0)
	count := 0
	// 生成[8, 18)层级的网格
	zoom := uint32(8)
	for ; zoom < 18; zoom++ {
		if zoom%2 == 0 {
			continue
		}
		for i := 0; i < len(d); i++ {
			var leftTop LngLat
			leftTop.Lat = 22.533389 - mesh_height*float64(int((22.533389-d[i].Lat)/mesh_height)) //22.966396 - mesh_height*float64(int((22.966396-d[i].Lat)/mesh_height))
			leftTop.Lng = 113.545749 + mesh_width*float64(int((d[i].Lng-113.545749)/mesh_width))
			ad := -1.0
			if d[i].Ad2 != nil {
				ad = *(d[i].Ad2)
			}
			x_0, y_0 := FromPixelToTileXY(FromLatLngToPixel(leftTop.Lat, leftTop.Lng, int(zoom)))
			outerRing := orb.Ring{
				{leftTop.Lng, leftTop.Lat},
				{leftTop.Lng + mesh_width, leftTop.Lat},
				{leftTop.Lng + mesh_width, leftTop.Lat - mesh_height},
				{leftTop.Lng, leftTop.Lat - mesh_height},
				{leftTop.Lng, leftTop.Lat},
			}
			outerPolygon := orb.Polygon{outerRing}
			//tile := maptile.New(uint32(x_0), uint32(y_0), maptile.Zoom(zoom))
			geometryCollection := orb.Collection{outerPolygon}
			f := geojson.NewFeature(geometryCollection)
			f.Properties = geojson.Properties{
				"ad": ad, //float64(120),
			}
			//id := ToID(uint8(zoom), uint32(x_0), uint32(y_0))
			id := TileID{X: uint32(x_0), Y: uint32(y_0), Z: zoom}
			_, ok := mtile[id]
			if !ok {
				mtile[id] = geojson.NewFeatureCollection()
			}
			mtile[id].Append(f)

			minLat, minLon, maxLat, maxLon := TileBounds(x_0, y_0, int(zoom))
			minLat = minLat + mesh_height*0.2
			minLon = minLon + mesh_width*0.2
			maxLat = maxLat - mesh_height*0.2
			maxLon = maxLon - mesh_width*0.2
			if leftTop.Lng > maxLon || leftTop.Lng < minLon || leftTop.Lat > maxLat || leftTop.Lat < minLat ||
				(leftTop.Lng+mesh_width) > maxLon || (leftTop.Lng+mesh_width) < minLon || (leftTop.Lat-mesh_height) > maxLat || (leftTop.Lat-mesh_height) < minLat {
				//if true {
				count++
				x_1, y_1 := FromPixelToTileXY(FromLatLngToPixel(leftTop.Lat, leftTop.Lng+mesh_width, int(zoom)))
				x_2, y_2 := FromPixelToTileXY(FromLatLngToPixel(leftTop.Lat-mesh_height, leftTop.Lng+mesh_width, int(zoom)))
				x_3, y_3 := FromPixelToTileXY(FromLatLngToPixel(leftTop.Lat-mesh_height, leftTop.Lng, int(zoom)))
				{
					outerRing := orb.Ring{
						{leftTop.Lng, leftTop.Lat},
						{leftTop.Lng + mesh_width, leftTop.Lat},
						{leftTop.Lng + mesh_width, leftTop.Lat - mesh_height},
						{leftTop.Lng, leftTop.Lat - mesh_height},
						{leftTop.Lng, leftTop.Lat},
					}
					outerPolygon := orb.Polygon{outerRing}
					//tile := maptile.New(uint32(x_0), uint32(y_0), maptile.Zoom(zoom))
					geometryCollection := orb.Collection{outerPolygon}
					f := geojson.NewFeature(geometryCollection)
					f.Properties = geojson.Properties{
						"ad": ad, //float64(120),
					}
					//id := ToID(uint8(zoom), uint32(x_1), uint32(y_1))
					id := TileID{X: uint32(x_1), Y: uint32(y_1), Z: zoom}
					_, ok := mtile[id]
					if !ok {
						mtile[id] = geojson.NewFeatureCollection()
					}
					mtile[id].Append(f)
				}
				{
					outerRing := orb.Ring{
						{leftTop.Lng, leftTop.Lat},
						{leftTop.Lng + mesh_width, leftTop.Lat},
						{leftTop.Lng + mesh_width, leftTop.Lat - mesh_height},
						{leftTop.Lng, leftTop.Lat - mesh_height},
						{leftTop.Lng, leftTop.Lat},
					}
					outerPolygon := orb.Polygon{outerRing}
					//tile := maptile.New(uint32(x_0), uint32(y_0), maptile.Zoom(zoom))
					geometryCollection := orb.Collection{outerPolygon}
					f := geojson.NewFeature(geometryCollection)
					f.Properties = geojson.Properties{
						"ad": ad, //float64(120),
					}
					//id := ToID(uint8(zoom), uint32(x_2), uint32(y_2))
					id := TileID{X: uint32(x_2), Y: uint32(y_2), Z: zoom}
					_, ok := mtile[id]
					if !ok {
						mtile[id] = geojson.NewFeatureCollection()
					}
					mtile[id].Append(f)
				}
				{
					outerRing := orb.Ring{
						{leftTop.Lng, leftTop.Lat},
						{leftTop.Lng + mesh_width, leftTop.Lat},
						{leftTop.Lng + mesh_width, leftTop.Lat - mesh_height},
						{leftTop.Lng, leftTop.Lat - mesh_height},
						{leftTop.Lng, leftTop.Lat},
					}
					outerPolygon := orb.Polygon{outerRing}
					//tile := maptile.New(uint32(x_0), uint32(y_0), maptile.Zoom(zoom))
					geometryCollection := orb.Collection{outerPolygon}
					f := geojson.NewFeature(geometryCollection)
					f.Properties = geojson.Properties{
						"ad": ad, //float64(120),
					}
					//id := ToID(uint8(zoom), uint32(x_3), uint32(y_3))
					id := TileID{X: uint32(x_3), Y: uint32(y_3), Z: zoom}
					_, ok := mtile[id]
					if !ok {
						mtile[id] = geojson.NewFeatureCollection()
					}
					mtile[id].Append(f)
				}
			}
		}
	}
	return mtile
}

// mercatorProjectionFromLatLng
func mercatorProjectionFromLatLng(lat_rad, lng_rad float64) (float64, float64) {
	a := math.Tan(math.Pi/4.0 + Fabs(lat_rad)/2)
	if a < 0 {
		return 0, 0
	}
	x_meter := lng_rad * EarthRadiusKm * 1000.0
	k := -1
	if lng_rad < 0 {
		k = -1
	} else {
		k = 1
	}
	y_meter := math.Log(a) * EarthRadiusKm * 1000.0 * float64(k)
	return x_meter, y_meter
}

func inverseMercatorProjectionFromXY(x_meter, y_meter float64) (float64, float64) {
	lng_rad := x_meter / 1000.0 / EarthRadiusKm
	k := -1
	if y_meter < 0 {
		k = -1
	} else {
		k = 1
	}
	lat_rad := (2*math.Atan(math.Exp(Fabs(y_meter/1000.0/EarthRadiusKm))) - (math.Pi / 2.0)) * float64(k)
	return lat_rad, lng_rad
}

func GetMeshPosition(lat, lng float64) (bool, float64, float64) {
	if lng != 0 && lat != 0 && lng > -180 && lat > -90 {
		var x_meter float64
		var y_meter float64
		var lat_rad float64
		var lng_rad float64

		x_meter, y_meter = mercatorProjectionFromLatLng(lat*math.Pi/180.0, lng*math.Pi/180.0)
		k0 := x_meter + 1
		if x_meter == 0 {
			k0 = x_meter + 1
		} else {
			k0 = x_meter
		}
		k1 := -1
		if lng < 0 {
			k1 = -1
		} else {
			k1 = 1
		}
		k2 := y_meter + 1
		if y_meter == 0 {
			k2 = y_meter + 1
		} else {
			k2 = y_meter
		}
		k3 := -1
		if lat < 0 {
			k3 = -1
		} else {
			k3 = 1
		}
		num_x := math.Ceil(Fabs(k0)/_rangeMeters) * float64(k1)
		num_y := math.Ceil(Fabs(k2)/_rangeMeters) * float64(k3)

		x_meter = float64(num_x) * _rangeMeters
		y_meter = float64(num_y) * _rangeMeters

		lat_rad, lng_rad = inverseMercatorProjectionFromXY(x_meter, y_meter)

		lat_ := lat_rad * 180.0 / math.Pi
		lng_ := lng_rad * 180.0 / math.Pi

		return true, lat_, lng_
	}

	return false, 0, 0
}

func Clip(n, minValue, maxValue float64) float64 {
	return math.Min(math.Max(n, minValue), maxValue)
}

func GetTileMatrixMaxXY(zoom int) (int, int) {
	xy := (1 << zoom)
	return xy - 1, xy - 1
}
func GetTileMatrixMinXY(zoom int) (int, int) {
	return 0, 0
}

func GetTileMatrixSizeXY(zoom int) (int, int) {
	sMin_x, sMin_y := GetTileMatrixMinXY(zoom)
	sMax_x, sMax_y := GetTileMatrixMaxXY(zoom)

	return sMax_x - sMin_x + 1, sMax_y - sMin_y + 1
}

func GetTileMatrixSizePixel(zoom int) (int, int) {
	s_x, s_y := GetTileMatrixSizeXY(zoom)
	return s_x * 256, s_y * 256
}

func FromLatLngToPixel(lat, lng float64, zoom int) (int, int) {
	lat = Clip(lat, -85.05112878, 85.05112878)
	lng = Clip(lng, -177, 177)

	x := (lng + 180) / 360
	sinLatitude := math.Sin(lat * math.Pi / 180)
	y := 0.5 - math.Log((1+sinLatitude)/(1-sinLatitude))/(4*math.Pi)

	s_x, s_y := GetTileMatrixSizePixel(zoom)
	mapSizeX := int(s_x)
	mapSizeY := int(s_y)

	ret_x := (int)(Clip(x*float64(mapSizeX)+0.5, 0, float64(mapSizeX)-1))
	ret_y := (int)(Clip(y*float64(mapSizeY)+0.5, 0, float64(mapSizeY)-1))

	return ret_x, ret_y
}

func FromPixelToTileXY(p_x, p_y int) (int, int) {
	return (int)(p_x / 256), (int)(p_y / 256)
}

func lngLatToTileXY(ll LngLat, tile TileID) (float64, float64) {
	totalTilesX := math.Pow(2, float64(tile.Z))
	totalTilesY := math.Pow(2, float64(tile.Z))
	lambda := (ll.Lng + 180) / 180 * math.Pi
	// phi: [-pi/2, pi/2]
	phi := ll.Lat / 180 * math.Pi
	tileX := lambda / (2 * math.Pi) * totalTilesX
	// [-1.4844, 1.4844] -> [1, 0]  * totalTilesY
	tileY := (math.Log(math.Tan(math.Pi/4-phi/2))/math.Pi/2 + 0.5) * totalTilesY
	return tileX - float64(tile.X), tileY - float64(tile.Y)
}

func ToID(z uint8, x uint32, y uint32) string {
	id := (((uint64(1)<<uint64(z))*uint64(y) + uint64(x)) * uint64(32)) + uint64(z)
	return strconv.FormatUint(id, 10)
}

func cmdEnc(id uint32, count uint32) uint32 {
	return (id & 0x7) | (count << 3)
}

func moveTo(count uint32) uint32 {
	return cmdEnc(1, count)
}

func lineTo(count uint32) uint32 {
	return cmdEnc(2, count)
}

func closePath(count uint32) uint32 {
	return cmdEnc(7, count)
}

func paramEnc(value int32) int32 {
	return (value << 1) ^ (value >> 31)
}

// Takes a string of the form `<z>/<x>/<y>` (for example, 1/2/3) and returns
// the individual uint32 values for x, y, and z if there was no error.
// Otherwise, err is set to a non `nil` value and x, y, z are set to 0.
func TilePathToXYZ(path string) (TileID, error) {
	xyzReg := regexp.MustCompile("(?P<z>[0-9]+)/(?P<x>[0-9]+)/(?P<y>[0-9]+)/(?P<account>[0-9A-Za-z]+)/(?P<key>[0-9A-Za-z]+)")
	matches := xyzReg.FindStringSubmatch(path)
	if len(matches) == 0 {
		return TileID{}, errors.New("unable to parse path as tile")
	}
	x, err := strconv.ParseUint(matches[2], 10, 32)
	if err != nil {
		return TileID{}, err
	}
	y, err := strconv.ParseUint(matches[3], 10, 32)
	if err != nil {
		return TileID{}, err
	}
	z, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		return TileID{}, err
	}
	// 账号名
	account := matches[4]
	if account == "" {
		return TileID{}, errors.New("Account is nil")
	}
	// 时间段生成的md5
	key := matches[5]
	if key == "" {
		return TileID{}, errors.New("Time key is nil")
	}
	// fmt.Println("拿到的账号名和时间段key: ", account, key)
	return TileID{X: uint32(x), Y: uint32(y), Z: uint32(z), Account: account, Key: key}, nil
}

func meterToDegree(length, lat float64) float64 {
	// return Math.pow((Math.pow(len, -1)) * 111319.49079327358, -1);
	degree := (2 * math.Pi * math.Cos((math.Pi/180.0)*lat) * 6378137.0) / 360.0
	return (1 / degree) * length * math.Cos((lat*math.Pi)/180.0)
}

func GetDistance(p1, p2 LngLat) float64 { //km
	dLat1InRad := p1.Lat * (math.Pi / 180.0)
	dLong1InRad := p1.Lng * (math.Pi / 180.0)
	dLat2InRad := p2.Lat * (math.Pi / 180.0)
	dLong2InRad := p2.Lng * (math.Pi / 180.0)
	dLongitude := dLong2InRad - dLong1InRad
	dLatitude := dLat2InRad - dLat1InRad
	a := math.Pow(math.Sin(dLatitude/2.0), 2.0) + math.Cos(dLat1InRad)*math.Cos(dLat2InRad)*math.Pow(math.Sin(dLongitude/2.0), 2.0)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	dDistance := EarthRadiusKm * c
	return dDistance
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func tileXYToPixelXY(tileX, tileY int) (pixelX, pixelY int) {
	return tileX << 8, tileY << 8
}

func gMapSize(levelOfDetail int) uint64 {
	return gTileSize << levelOfDetail
}

func pixelXYToLatLon(pixelX, pixelY, levelOfDetail int) (lat, lon float64) {
	mapSize := float64(gMapSize(levelOfDetail))
	x := (clamp(float64(pixelX), 0, mapSize-1) / mapSize) - 0.5
	y := 0.5 - (clamp(float64(pixelY), 0, mapSize-1) / mapSize)
	lat = 90 - 360*math.Atan(math.Exp(-y*2*math.Pi))/math.Pi
	lon = 360 * x
	return
}

// TileBounds returns the lat/lon bounds around a tile.
func TileBounds(tileX, tileY, tileZ int,
) (minLat, minLon, maxLat, maxLon float64) {
	levelOfDetail := tileZ
	size := int(1 << levelOfDetail)
	pixelX, pixelY := tileXYToPixelXY(tileX, tileY)
	maxLat, minLon = pixelXYToLatLon(pixelX, pixelY, levelOfDetail)
	pixelX, pixelY = tileXYToPixelXY(tileX+1, tileY+1)
	minLat, maxLon = pixelXYToLatLon(pixelX, pixelY, levelOfDetail)
	if size == 0 || tileX%size == 0 {
		minLon = gMinLon
	}
	if size == 0 || tileX%size == size-1 {
		maxLon = gMaxLon
	}
	if tileY <= 0 {
		maxLat = gMaxLat
	}
	if tileY >= size-1 {
		minLat = gMinLat
	}
	return
}

// 两点确定距离
// func GetDistance(p1, p2 LngLat) float64 { //km
// 	dLat1InRad := p1.Lat * (math.Pi / 180.0)
// 	dLong1InRad := p1.Lng * (math.Pi / 180.0)
// 	dLat2InRad := p2.Lat * (math.Pi / 180.0)
// 	dLong2InRad := p2.Lng * (math.Pi / 180.0)
// 	dLongitude := dLong2InRad - dLong1InRad
// 	dLatitude := dLat2InRad - dLat1InRad
// 	a := math.Pow(math.Sin(dLatitude/2.0), 2.0) + math.Cos(dLat1InRad)*math.Cos(dLat2InRad)*math.Pow(math.Sin(dLongitude/2.0), 2.0)
// 	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
// 	dDistance := EarthRadiusKm * c
// 	return dDistance
// }

// GetMeshPosition
// 生成网格左上角点的经纬度
// func GetMeshPosition(lat, lng float64) (float64, float64) {
// 	// 网格宽度
// 	mesh_width := 50.0 * math.Abs(113.545749-115.231401) / (GetDistance(LngLat{Lng: 113.545749, Lat: 22.966396}, LngLat{Lng: 115.231401, Lat: 22.966396}) * 1000.0)
// 	// 网格高度
// 	mesh_height := 50.0 * math.Abs(22.966396-22.533389) / (GetDistance(LngLat{Lng: 113.545749, Lat: 22.966396}, LngLat{Lng: 113.545749, Lat: 22.533389}) * 1000.0)
// 	// 数据点生成网格的左上角点的经纬度
// 	lat_1 := 22.966396 - mesh_height*float64(int((22.966396-lat)/mesh_height))
// 	lng_1 := 113.545749 + mesh_width*float64(int((lng-113.545749)/mesh_width))
// 	// lat_1 = FloatRound(lat_1, 6)
// 	// lng_1 = FloatRound(lng_1, 6)
// 	return lat_1, lng_1
// }

// func FloatRound(f float64, n int) float64 {
// 	format := "%." + strconv.Itoa(n) + "f"
// 	res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
// 	return res
// }
