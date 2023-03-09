package utils

import (
	"math"

	m "github.com/murphy214/mercantile"
	geojson "github.com/paulmach/go.geojson"
)

// CreateGeojson
// 生成geojson数据
func CreateGeojson(d []*Data) map[m.TileID]*geojson.FeatureCollection {
	// 存放不同层级的瓦片
	mtile := make(map[m.TileID]*geojson.FeatureCollection)
	mesh_width := 50.0 * math.Abs(113.545749-115.231401) / (GetDistance(LngLat{Lng: 113.545749, Lat: 22.966396}, LngLat{Lng: 115.231401, Lat: 22.966396}) * 1000.0)
	mesh_height := 50.0 * math.Abs(22.966396-22.533389) / (GetDistance(LngLat{Lng: 113.545749, Lat: 22.966396}, LngLat{Lng: 113.545749, Lat: 22.533389}) * 1000.0)
	zoom := uint32(9)
	for ; zoom < 19; zoom++ {
		if zoom%2 == 0 {
			continue
		}
		for i := 0; i < len(d); i++ {
			// d[i].Lat, d[i].Lng = geohash.Decode(d[i].Geoid)
			var leftTop LngLat
			leftTop.Lat = 22.533389 - mesh_height*float64(int((22.533389-d[i].Lat)/mesh_height)) //22.966396 - mesh_height*float64(int((22.966396-d[i].Lat)/mesh_height))
			leftTop.Lng = 113.545749 + mesh_width*float64(int((d[i].Lng-113.545749)/mesh_width))
			// ad := -1.0
			// if d[i].Ad2 != nil {
			// 	ad = *(d[i].Ad2)
			// }
			x_0, y_0 := FromPixelToTileXY(FromLatLngToPixel(leftTop.Lat, leftTop.Lng, int(zoom)))
			if true {
				geometry := geojson.NewPolygonGeometry([][][]float64{{
					{leftTop.Lng, leftTop.Lat},
					{leftTop.Lng + mesh_width, leftTop.Lat},
					{leftTop.Lng + mesh_width, leftTop.Lat - mesh_height},
					{leftTop.Lng, leftTop.Lat - mesh_height},
					{leftTop.Lng, leftTop.Lat},
				}})
				f := geojson.NewFeature(geometry)
				f.Properties = map[string]interface{}{
					"ad_2":  d[i].Ad2,
					"ad_3":  d[i].Ad3,
					"ad_4":  d[i].Ad4,
					"ad_7":  d[i].Ad7,
					"ad_16": d[i].Ad16,
				}
				// id := m.TileID{X: uint32(x_0), Y: uint32(y_0), Z: zoom}
				id := m.TileID{X: int64(x_0), Y: int64(y_0), Z: uint64(zoom)}
				_, ok := mtile[id]
				if !ok {
					mtile[id] = geojson.NewFeatureCollection()
				}
				mtile[id].Features = append(mtile[id].Features, f)
			}

			minLat, minLon, maxLat, maxLon := TileBounds(x_0, y_0, int(zoom))
			minLat = minLat + mesh_height*0.2
			minLon = minLon + mesh_width*0.2
			maxLat = maxLat - mesh_height*0.2
			maxLon = maxLon - mesh_width*0.2
			if leftTop.Lng > maxLon || leftTop.Lng < minLon || leftTop.Lat > maxLat || leftTop.Lat < minLat ||
				(leftTop.Lng+mesh_width) > maxLon || (leftTop.Lng+mesh_width) < minLon || (leftTop.Lat-mesh_height) > maxLat || (leftTop.Lat-mesh_height) < minLat {
				//if true {
				x_1, y_1 := FromPixelToTileXY(FromLatLngToPixel(leftTop.Lat, leftTop.Lng+mesh_width, int(zoom)))
				x_2, y_2 := FromPixelToTileXY(FromLatLngToPixel(leftTop.Lat-mesh_height, leftTop.Lng+mesh_width, int(zoom)))
				x_3, y_3 := FromPixelToTileXY(FromLatLngToPixel(leftTop.Lat-mesh_height, leftTop.Lng, int(zoom)))
				{
					geometry := &geojson.Geometry{
						Type: "Polygon",
						Polygon: [][][]float64{{
							{leftTop.Lng, leftTop.Lat},
							{leftTop.Lng + mesh_width, leftTop.Lat},
							{leftTop.Lng + mesh_width, leftTop.Lat - mesh_height},
							{leftTop.Lng, leftTop.Lat - mesh_height},
							{leftTop.Lng, leftTop.Lat},
						}},
					}
					f := geojson.NewFeature(geometry)
					f.Properties = map[string]interface{}{
						"ad_2":  d[i].Ad2,
						"ad_3":  d[i].Ad3,
						"ad_4":  d[i].Ad4,
						"ad_7":  d[i].Ad7,
						"ad_16": d[i].Ad16,
					}
					// id := m.TileID{X: uint32(x_1), Y: uint32(y_1), Z: zoom}
					id := m.TileID{X: int64(x_1), Y: int64(y_1), Z: uint64(zoom)}
					_, ok := mtile[id]
					if !ok {
						mtile[id] = geojson.NewFeatureCollection()
					}
					mtile[id].Features = append(mtile[id].Features, f)
				}
				{
					geometry := &geojson.Geometry{
						Type: "Polygon",
						Polygon: [][][]float64{{
							{leftTop.Lng, leftTop.Lat},
							{leftTop.Lng + mesh_width, leftTop.Lat},
							{leftTop.Lng + mesh_width, leftTop.Lat - mesh_height},
							{leftTop.Lng, leftTop.Lat - mesh_height},
							{leftTop.Lng, leftTop.Lat},
						}},
					}
					f := geojson.NewFeature(geometry)
					f.Properties = map[string]interface{}{
						"ad_2":  d[i].Ad2,
						"ad_3":  d[i].Ad3,
						"ad_4":  d[i].Ad4,
						"ad_7":  d[i].Ad7,
						"ad_16": d[i].Ad16,
					}
					// id := m.TileID{X: uint32(x_2), Y: uint32(y_2), Z: zoom}
					id := m.TileID{X: int64(x_2), Y: int64(y_2), Z: uint64(zoom)}
					_, ok := mtile[id]
					if !ok {
						mtile[id] = geojson.NewFeatureCollection()
					}
					mtile[id].Features = append(mtile[id].Features, f)
				}
				{
					geometry := &geojson.Geometry{
						Type: "Polygon",
						Polygon: [][][]float64{{
							{leftTop.Lng, leftTop.Lat},
							{leftTop.Lng + mesh_width, leftTop.Lat},
							{leftTop.Lng + mesh_width, leftTop.Lat - mesh_height},
							{leftTop.Lng, leftTop.Lat - mesh_height},
							{leftTop.Lng, leftTop.Lat},
						}},
					}
					f := geojson.NewFeature(geometry)
					f.Properties = map[string]interface{}{
						"ad_2":  d[i].Ad2,
						"ad_3":  d[i].Ad3,
						"ad_4":  d[i].Ad4,
						"ad_7":  d[i].Ad7,
						"ad_16": d[i].Ad16,
					}
					// id := m.TileID{X: uint32(x_3), Y: uint32(y_3), Z: zoom}
					id := m.TileID{X: int64(x_3), Y: int64(y_3), Z: uint64(zoom)}
					_, ok := mtile[id]
					if !ok {
						mtile[id] = geojson.NewFeatureCollection()
					}
					mtile[id].Features = append(mtile[id].Features, f)
				}
			}
		}
	}
	return mtile
}

// func CreateVGeojsonVt()
