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

// a GeoGrid is a GeoRect split into cells
// of width CellWdith and height CellHeight.
// This can be used to map GeoPoints to cells
// CellWidth and CellHeight are in meters
type GeoGrid struct {
	Boundary GeoRect
	CellWidth, CellHeight float64
}

func (r GeoRect) MinLat() float64 {
	return r.SouthEast.LatitudeInDegrees
}

func (r GeoRect) MaxLat() float64 {
	return r.NorthWest.LatitudeInDegrees
}

func (r GeoRect) MinLon() float64 {
	return r.NorthWest.LongitudeInDegrees
}

func (r GeoRect) MaxLon() float64 {
	return r.SouthEast.LongitudeInDegrees
}

// accurate assuming these points no further than a few miles apart
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

func PathBoundary(paths ...[]GeoPoint) GeoRect {

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

func MakeGrid(cellWidth, cellHeight float64, boundary GeoRect) *GeoGrid {
	grid := GeoGrid{
		Boundary: boundary,
		CellWidth: cellWidth,
		CellHeight: cellHeight,
	}
	return &grid
}

func (g *GeoGrid) MapPoint(point GeoPoint) cellCoord {

	// meridian distances don't vary with longitude
	gridNorthSouthDistance := GeoDistance(GeoPoint{g.Boundary.MinLat(), g.Boundary.MinLon()},
		GeoPoint{g.Boundary.MaxLat(), g.Boundary.MinLon()})

	pointOffset := GeoDistance(GeoPoint{point.LatitudeInDegrees, g.Boundary.MinLon()},
		GeoPoint{g.Boundary.MaxLat(), g.Boundary.MinLon()})

	latIndex := cellSearch(0, gridNorthSouthDistance, g.CellHeight, pointOffset, 0)

	// distances along a parallel do vary with latitude
	// e.g. a degree of latitude is a longer distance at the equator
	// than further north / south. So the lat of the point for correctness.
	// The grid in fact has more cells per row closer to the equator
	gridEastWestDistance := GeoDistance(GeoPoint{point.LatitudeInDegrees, g.Boundary.MinLon()},
		GeoPoint{point.LatitudeInDegrees, g.Boundary.MaxLon()})

	pointOffset = GeoDistance(GeoPoint{g.Boundary.MinLat(), point.LongitudeInDegrees},
		GeoPoint{g.Boundary.MinLat(), g.Boundary.MinLon()})

	lonIndex := cellSearch(0, gridEastWestDistance, g.CellWidth, pointOffset, 0)
	
	return cellCoord(fmt.Sprintf("%d_%d", latIndex, lonIndex))
}

func (g *GeoGrid) MapPath(path []GeoPoint) []cellCoord {
	coordList := make([]cellCoord, len(path))
	for i, point := range(path) {
		coordList[i] = g.MapPoint(point)
	}
	return coordList
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

func PathLengthInMeters(path []GeoPoint) float64 {

	pathLength := 0.0
	for i := 1; i < len(path); i++ {
		pathLength += GeoDistance(path[i], path[i-1])
	}
	return pathLength
}

func (r1 GeoRect) Overlaps(r2 GeoRect) bool {
	// longitude always decreases to the west, until the anti-meridian
	// If data is from Taveuni or somewhere like that this will just be wrong.
	if r1.MinLon() >= r2.MaxLon() || r2.MinLon() >= r1.MaxLon() {
		return false
	}

	// Latitude can get weird around the poles, but this should still work
	// More likely is th rectangles being totally misrepresented
	if r1.MinLat() >= r2.MaxLat() || r2.MinLat() >= r1.MaxLat() {
		return false
	}

	return true	
}

func MergeGeoRect(rlist ...GeoRect) GeoRect {

	merged := GeoRect{
		NorthWest: GeoPoint{-91, 181},
		SouthEast: GeoPoint{91, -181},
	}

	for _, r := range(rlist) {

		if r.MinLat() < merged.MinLat() {
			merged.SouthEast.LatitudeInDegrees = r.SouthEast.LatitudeInDegrees
		}
		if r.MaxLat() > merged.MaxLat() {
			merged.NorthWest.LatitudeInDegrees = r.NorthWest.LatitudeInDegrees
		}
		if r.MinLon() < merged.MinLon() {
			merged.NorthWest.LongitudeInDegrees = r.NorthWest.LongitudeInDegrees
		}
		if r.MaxLon() > merged.MaxLon() {
			merged.SouthEast.LongitudeInDegrees = r.SouthEast.LongitudeInDegrees
		}
	}

	return merged
}

func PathSimilarity(cellWidth, cellHeight float64, path1, path2 []GeoPoint) float64 {


	boundary1 := PathBoundary(path1)
	boundary2 := PathBoundary(path2)

	if !boundary1.Overlaps(boundary2) {
		return 0.0
	}
	sharedBoundary := MergeGeoRect(boundary1, boundary2)

	grid := MakeGrid(cellWidth, cellHeight, sharedBoundary)

	coords1 := grid.MapPath(path1)
	coords2 := grid.MapPath(path2)

	// find the size of the difference between the two sets of coordinates
	inCoords1 := make(map[cellCoord]bool)
	inCoords2 := make(map[cellCoord]bool)
	for _, coord := range(coords1) {
		inCoords1[coord] = true
	}
	for _, coord := range(coords2) {
		inCoords2[coord] = true		
	}

	diffCount := 0
	for coord := range(inCoords1) {
		if !inCoords2[coord] {
			diffCount++
		}
	}
	for coord := range(inCoords2) {
		if !inCoords1[coord] {
			diffCount++
		}
	}
	
	totalPathLength := PathLengthInMeters(path1) + PathLengthInMeters(path2)
	return 1 - (float64(diffCount) * cellWidth / totalPathLength)
}
