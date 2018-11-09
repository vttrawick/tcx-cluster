package geo

import (
	"math"
	"fmt"
)

const EarthRadiusInMeters float64 = 6371.0088 * 1000
// a change in one degree of latitude along a meridian is about 111.32km
const LatitudeDegreesToMeters float64 = 111.32 * 1000

// a coordinate on the earth represented by lat / lon degrees
type GeoPoint struct {
	LatitudeInDegrees, LongitudeInDegrees float64
}

// the minimum area bounded by the four lat / lon lines
// all values in degrees
type GeoRect struct {
	NorthWest, SouthEast GeoPoint
}

type cellCoord string

// a geo grid is some GeoRect split into cells
// of width CellWdith and height CellHeight.
// The grid can contain points, which are represented as a map
// from grid coordinates to a slice of GeoPoints
// CellWidth and CellHeight are in meters
type GeoGrid struct {
	Boundary GeoRect
	CellWidth, CellHeight float64
	Cells map[cellCoord][]GeoPoint
}

// assuming these points are only a few miles away (up to 30)
// always returns a positive real number
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

	minLat := float64(91)
	maxLat := float64(-91)
	minLon := float64(181)
	maxLon := float64(-181)

	for _, path := range(paths) {
		for _, pt := range(path) {
			if pt.LatitudeInDegrees < minLat {
				minLat = pt.LatitudeInDegrees
			}
			if pt.LatitudeInDegrees > maxLat {
				maxLat = pt.LatitudeInDegrees
			}
			if pt.LongitudeInDegrees < minLon {
				minLon = pt.LongitudeInDegrees
			}
			if pt.LongitudeInDegrees > maxLon {
				maxLon = pt.LongitudeInDegrees
			}
		}
	}
	return GeoRect{
		NorthWest: GeoPoint{maxLat, minLon},
		SouthEast: GeoPoint{minLat, maxLon},
	}
}

// map each point in a set of paths to a coordinate on a GeoGrid
func PathToGeoGrid(cellWidth, cellHeight float64, paths ...[]GeoPoint) *GeoGrid {
	
	boundary := GridBoundary(paths...)

	cells := make(map[cellCoord][]GeoPoint)

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

	minLat := grid.Boundary.SouthEast.LatitudeInDegrees
	maxLat := grid.Boundary.NorthWest.LatitudeInDegrees
	minLon := grid.Boundary.NorthWest.LongitudeInDegrees
	maxLon := grid.Boundary.SouthEast.LongitudeInDegrees

	// meridian distances don't vary with latitude / longitude
	meridianDistance := GeoDistance(GeoPoint{minLat, minLon},
		GeoPoint{maxLat, minLon})

	pointOffset := GeoDistance(GeoPoint{point.LatitudeInDegrees, minLon},
		GeoPoint{maxLat, minLon})

	latIndex := cellSearch(0, meridianDistance, grid.CellHeight, pointOffset, 0)

	// however, distances along a parallel vary with latitude
	parallelDistance := GeoDistance(GeoPoint{point.LatitudeInDegrees, minLon},
		GeoPoint{point.LatitudeInDegrees, maxLon})

	pointOffset = GeoDistance(GeoPoint{minLat, point.LongitudeInDegrees},
		GeoPoint{minLat, maxLon})

	lonIndex := cellSearch(0, parallelDistance, grid.CellWidth, pointOffset, 0)
	
	return cellCoord(fmt.Sprintf("%d_%d", latIndex, lonIndex))
}

// binary search through a given range to find the index of the cell
func cellSearch(min, max, cellSize, loc float64, offset int) int {

	if max - min <= cellSize {
		return offset
	}

	// the "right-dividing box index" or the index to the cell
	// to the right of the midpoint line of the current range
	rdbi := int(math.Floor(math.Ceil((max - min) / cellSize) / 2))
	midpoint := cellSize * float64(rdbi + offset)

	if loc == midpoint {
		return rdbi - 1
	} else if loc > midpoint {
		min = midpoint
		offset += rdbi
		return cellSearch(min, max, cellSize, loc, offset)
	} else {
		max = midpoint
		return cellSearch(min, max, cellSize, loc, offset)
	}
}
