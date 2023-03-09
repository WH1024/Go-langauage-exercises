package utils

import (
	"math"

	m "github.com/murphy214/mercantile"
	"github.com/xiaolingis/geojsonvt"
)

func TileFromGeoJSON(features []geojsonvt.Feature, tileid m.TileID, options geojsonvt.Config) geojsonvt.Tile {
	tile := geojsonvt.NewTile()
	tile.TileID = m.Parent(tileid) // creating from parent first
	for _, feature := range features {

		tile.NumFeatures++

		minX := feature.MinX
		minY := feature.MinY
		maxX := feature.MaxX
		maxY := feature.MaxY
		if minX < tile.MinX {
			tile.MinX = minX
		}
		if minY < tile.MinY {
			tile.MinY = minY
		}
		if maxX > tile.MaxX {
			tile.MaxX = maxX
		}
		if maxY > tile.MaxY {
			tile.MaxY = maxY
		}
	}
	tile.Source = features
	tile.Options = options
	return tile.SplitTileChildren()[tileid]
}

func wrap(features []geojsonvt.Feature, options geojsonvt.Config) []geojsonvt.Feature {
	buffer := options.Buffer / options.Extent
	merged := features

	left := clip(features, 1, float64(-1.0-buffer), float64(buffer), 0, -1, 2, options)     // left world copy
	right := clip(features, 1, float64(1.0-buffer), float64(2.0+buffer), 0, -1, 2, options) // right world copy

	if len(left) > 0 || len(right) > 0 {
		merged = clip(features, 1, float64(-buffer), float64(1.0+buffer), 0, -1, 2, options) // center world copy
		if len(left) > 0 {
			merged = append(merged, shiftFeatureCoords(left, float64(1))...)
		}
		if len(right) > 0 {
			merged = append(merged, shiftFeatureCoords(right, float64(-1))...)
		}
	}
	return merged
}

func clip(features []geojsonvt.Feature, scale int, k1 float64, k2 float64, axis int, minAll float64, maxAll float64, options geojsonvt.Config) []geojsonvt.Feature {
	k1 = k1 / float64(scale)
	k2 = k2 / float64(scale)
	if minAll >= k1 && maxAll < k2 {
		return features // trivial accept
	} else if maxAll < k1 || minAll >= k2 {
		return []geojsonvt.Feature{}
	}

	clipped := []geojsonvt.Feature{}
	var min, max float64
	for _, feature := range features {
		if axis == 0 {
			min = feature.MinX
			max = feature.MaxX
		} else if axis == 1 {
			min = feature.MinY
			max = feature.MaxY

		}
		boolval := true
		if min >= k1 && max < k2 { // trivial accept
			//clipped = append(clipped, feature)
			//boolval = false
		} else if max < k1 || min >= k2 { // trivial reject
			//boolval = false
		}

		if boolval {
			clipgeom, _ := feature.Geometry.Clip(k1, k2, axis, feature.Type == "MultiPolygon" || feature.Type == "Polygon")
			boolval2 := (len(clipgeom.LineString) == 0) && (len(clipgeom.Point) == 0) && (len(clipgeom.Polygon) == 0) && (len(clipgeom.MultiLineString) == 0) && (len(clipgeom.MultiPoint) == 0) && (len(clipgeom.MultiPolygon) == 0)
			if !boolval2 {
				clipped = append(clipped, CreateFeature(feature.ID, clipgeom, feature.Tags))
			}
		}
	}
	return clipped
}

// creates a feature
func CreateFeature(id interface{}, geom geojsonvt.Geometry, tags map[string]interface{}) geojsonvt.Feature {
	feature := geojsonvt.Feature{
		ID:       id,
		Geometry: geom,
		Type:     geom.Type,
		Tags:     tags,
	}
	calcBBox(&feature)
	return feature
}

func calcBBox(feature *geojsonvt.Feature) {
	switch feature.Type {

	case "Point":
		calcLineBBox(feature, feature.Geometry.Point)
	case "MultiPoint":
		calcLineBBox(feature, feature.Geometry.MultiPoint)

	case "LineString":
		calcLineBBox(feature, feature.Geometry.LineString)

	case "MultiLineString":
		for i := range feature.Geometry.MultiLineString {
			calcLineBBox(feature, feature.Geometry.MultiLineString[i])
		}
	case "Polygon":
		for i := range feature.Geometry.Polygon {
			calcLineBBox(feature, feature.Geometry.Polygon[i])
		}
	case "MultiPolygon":
		for ii := range feature.Geometry.MultiPolygon {
			for i := range feature.Geometry.MultiPolygon[ii] {
				calcLineBBox(feature, feature.Geometry.MultiPolygon[ii][i])
			}
		}
	}
}

func calcLineBBox(feature *geojsonvt.Feature, line []float64) {
	for i := 0; i < len(line); i += 3 {
		feature.MinX = math.Min(feature.MinX, line[i])
		feature.MinY = math.Min(feature.MinY, line[i+1])
		feature.MaxX = math.Max(feature.MaxX, line[i])
		feature.MaxY = math.Max(feature.MaxY, line[i+1])
	}
}

func shiftFeatureCoords(features []geojsonvt.Feature, offset float64) []geojsonvt.Feature {
	for featurepos, feature := range features {
		switch feature.Type {
		case "Point":
			feature.Geometry.Point = shiftCoords(feature.Geometry.Point, offset)
		case "MultiPoint":
			feature.Geometry.MultiPoint = shiftCoords(feature.Geometry.MultiPoint, offset)
		case "LineString":
			feature.Geometry.LineString = shiftCoords(feature.Geometry.LineString, offset)
		case "MultiLineString":
			for pos := range feature.Geometry.MultiLineString {
				feature.Geometry.MultiLineString[pos] = shiftCoords(feature.Geometry.MultiLineString[pos], offset)
			}

		case "Polygon":
			for pos := range feature.Geometry.Polygon {
				feature.Geometry.Polygon[pos] = shiftCoords(feature.Geometry.Polygon[pos], offset)
			}
		case "MultiPolygon":
			for i := range feature.Geometry.MultiPolygon {
				for pos := range feature.Geometry.MultiPolygon[i] {
					feature.Geometry.MultiPolygon[i][pos] = shiftCoords(feature.Geometry.MultiPolygon[i][pos], offset)
				}
			}
		}
		features[featurepos] = feature
	}
	return features
}

// shifts each features coordinates over
func shiftCoords(points []float64, offset float64) []float64 {
	for i := 0; i < len(points); i += 3 {
		points[i] = points[i] + offset
	}
	return points
}
