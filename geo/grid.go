package geo

import (
	"math"
)

const EarthRadiusInMeters float64 = 6371.0088 * 1000

// a coordinate on the earth represented by lat / lon degrees
type GeoPoint struct {
	LatitudeInDegrees, LongitudeInDegrees float64
}

// the minimum area bounded by the four lat / lon lines
// all values in degrees
type GeoRect struct {
	MinLat, MaxLat, MinLon, MaxLon float64
}

type cellCoord string

// a geo grid is some GeoRect split into cells
// of width CellWdith and height CellHeight.
// The grid can contain points, which are represented as a map
// from grid coordinates to a slice of GeoPoints
// CellWidth and CellHeight are in meters
type GeoGrid struct {
	Boundary GeoRect
	CellWidth, CellHeight int
	Cells map[cellCoord][]GeoPoint
}

// assuming these points are only a few miles away (up to 30)
func GeoDistance(p1, p2 GeoPoint) float64 {

	dLat := p1.LatitudeInDegrees - p2.LatitudeInDegrees
	dLat = math.Pi * dLat / 180

	dLon := p1.LongitudeInDegrees - p2.LongitudeInDegrees
	dLon = math.Pi * dLon / 180

	// take the average latitude and use that to shorten the latitude circle
	avgLat := (p1.LatitudeInDegrees + p2.LatitudeInDegrees) / 2
	avgLat = math.Pi * avgLat / 180
	dLon = math.Cos(avgLat) * dLon
	
	return EarthRadiusInMeters * math.Sqrt((dLat * dLat) + (dLon * dLon))
}

func GridBoundary(paths ...[]GeoPoint) GeoRect {

	r := GeoRect{
		MinLat: 91,
		MaxLat: -91,
		MinLon: 181,
		MaxLon: -181,
	}
	for _, path := range(paths) {
		for _, pt := range(path) {
			if pt.LatitudeInDegrees < r.MinLat {
				r.MinLat = pt.LatitudeInDegrees
			}
			if pt.LatitudeInDegrees > r.MaxLat {
				r.MaxLat = pt.LatitudeInDegrees
			}
			if pt.LongitudeInDegrees < r.MinLon {
				r.MinLon = pt.LongitudeInDegrees
			}
			if pt.LongitudeInDegrees > r.MaxLon {
				r.MaxLon = pt.LongitudeInDegrees
			}
		}
	}
	return r
}

// draws a GeoGrid onto a set of paths
func PathToGeoGrid(cellWidth, cellHeight int, paths ...[]GeoPoint) *GeoGrid {
	
	boundary := GridBoundary(paths...)

	var cells map[cellCoord][]GeoPoint

	grid := GeoGrid{
		Boundary: boundary,
		CellWidth: cellWidth,
		CellHeight: cellHeight,
		Cells: cells,
	}

	for _, path := range paths {
		for _, point := range path {
			coord := cellForPoint(point, &grid)
			if grid.Cells[coord] == nil {
				grid.Cells[coord] = make([]GeoPoint, 1)
			}
			grid.Cells[coord] = append(grid.Cells[coord], point)
		}
	}
	return &grid
}

func cellForPoint(point GeoPoint, grid *GeoGrid) cellCoord {
	// do lat and long separately
	// latDistance := grid.Boundary.MaxLat - grid.Boundary.MinLat

	// what's the difference in meters between latDistance at minLon and latDistance at maxLon
	// over how many rows would you need to add another cell?
	// the grid has a ragged edge

	// distanceInMeters := GeoDistance(GeoPoint{ }, GeoPoint{ })

	return cellCoord("0_0")
}
