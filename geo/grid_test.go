package geo

import(
	"testing"
	"reflect"
	"math"
	"strings"
	"strconv"
	"fmt"
)

var path1 = []GeoPoint{
	GeoPoint{42.365592, -71.103875},
	GeoPoint{42.364776, -71.110749},
	GeoPoint{42.364237, -71.116022},
	GeoPoint{42.361439, -71.115968},
	GeoPoint{42.362285, -71.113515},
	GeoPoint{42.365115, -71.104975},
}

var path2 = []GeoPoint{
	GeoPoint{42.365218, -71.104578},
	GeoPoint{42.362285, -71.113429},
	GeoPoint{42.360644, -71.113225},
	GeoPoint{42.360192, -71.112774},
	GeoPoint{42.364037, -71.108278},
	GeoPoint{42.365223, -71.104852},
}

var path3 = []GeoPoint{
	GeoPoint{42.365127, -71.103168},
	GeoPoint{42.360831, -71.096162},
	GeoPoint{42.359000, -71.100175},
	GeoPoint{42.361164, -71.103930},
	GeoPoint{42.365516, -71.103951},
	GeoPoint{42.364667, -71.102567},
}
	
func TestGridBoundary(t *testing.T) {

	// find the min / max lat long for one path
	boundary1 := GridBoundary(path1)

	// boundary should be:
	// {
	//   maxLat: 42.365592
	//   minLat: 42.361439
	//   maxLon: -71.103875
	//   minLon: -71.116022
	// }

	t1 := GeoRect{
		NorthWest: GeoPoint{42.365592, -71.116022},
		SouthEast: GeoPoint{42.361439, -71.103875},
	}
	if !reflect.DeepEqual(boundary1, t1) {
		t.Errorf("incorrect boundary structure for 1-path example")
	}

	// then for multiple
	boundary2 := GridBoundary(path1, path2, path3)

	// boundary should be
	// {
	//   maxLat: 42.365592
	//   minLat: 42.359000
	//   maxLon: -71.096162
	//   minLon: -71.116022
	// }

	t2 := GeoRect{
		NorthWest: GeoPoint{42.365592, -71.116022},
		SouthEast: GeoPoint{42.359000, -71.096162},
	}

	if !reflect.DeepEqual(boundary2, t2) {
		t.Errorf("incorrect boundary structure for 3-path example")
	}
}

func TestGeoDistance(t *testing.T) {

	// two points relatively close together (half a mile)
	d1 := GeoDistance(GeoPoint{42.353381, -71.107131}, GeoPoint{42.356941, -71.092647})
	if math.Abs(d1 - 1250) > 7 {
		t.Errorf("GeoDistance off for two close points")
	}

	// two points very close together (a few feet)
	d2 := GeoDistance(GeoPoint{42.364378, -71.114549}, GeoPoint{42.364447, -71.114603})
	if math.Abs(d2 - 10) > 7 {
		t.Errorf("GeoDistance off for two very close points")
	}
	
	// the same two points
	d3 := GeoDistance(GeoPoint{42.365673, -71.104100}, GeoPoint{42.365673, -71.104100})
	if d3 != 0 {
		t.Errorf("GeoDistance off for exact two points")
	}
	
	// two points at the same longitude
	d4 := GeoDistance(GeoPoint{42.365673, -71.104100}, GeoPoint{42.845683, -71.104100})
	if math.Abs(d4 - 53370) > 7 {
		t.Errorf("GeoDistance off for two points on same line of longitude")
	}
	
	// two points at the same latitude
	d5 := GeoDistance(GeoPoint{42.365673, -71.304331}, GeoPoint{42.365673, -71.104100})
	if math.Abs(d5 - 16450) > 7 {
		t.Errorf("GeoDistance off for two points on same line of latitude")
	}
	
	// two points relatively far away with varying lat / lon (~10 kilometers)
	d6 := GeoDistance(GeoPoint{42.365121, -71.212806}, GeoPoint{42.366810, -71.068591})
	if math.Abs(d6 - 11850) > 7 {
		t.Errorf("GeoDistance off for two relatively distant points")
	}	
}


func TestPathToGeoGrid(t *testing.T) {

	width := 7.0
	height := 7.0

	g1 := PathToGeoGrid(width, height, path1)

	// some sanity checks on the grid
	if (g1.CellWidth != width) {
		t.Errorf("grid cell width differs from input")
	}
	if (g1.CellHeight != height) {
		t.Errorf("grid cell height differs from input")
	}
	b1 := GeoRect{
		NorthWest: GeoPoint{42.365592, -71.116022},
		SouthEast: GeoPoint{42.361439, -71.103875},
	}
	if !reflect.DeepEqual(g1.Boundary, b1) {
		t.Errorf("grid boundary does not match expected")
	}
	// checking the cells is more involved
	// split it off into its own function
	checkCells(t, g1)
}

func checkCells(t *testing.T, grid *GeoGrid) {

	minLatCoord := math.MaxInt64
	maxLatCoord := -1
	minLonCoord := math.MaxInt64
	maxLonCoord := -1
	for coord, _ := range(grid.Cells) {
		parts := strings.Split(string(coord), "_")

		latCoord, _ := strconv.Atoi(parts[0])
		if latCoord > maxLatCoord {
			maxLatCoord = latCoord
		} else if latCoord < minLatCoord {
			minLatCoord = latCoord
		}

		lonCoord, _ := strconv.Atoi(parts[1])
		if lonCoord > maxLonCoord {
			
			maxLonCoord = lonCoord
		} else if lonCoord < minLonCoord {
			minLonCoord = lonCoord
		}
	}

	minLat := grid.Boundary.SouthEast.LatitudeInDegrees
	maxLat := grid.Boundary.NorthWest.LatitudeInDegrees
	minLon := grid.Boundary.NorthWest.LongitudeInDegrees
	maxLon := grid.Boundary.SouthEast.LongitudeInDegrees

	meridianDistance := GeoDistance(GeoPoint{minLat, minLon}, GeoPoint{maxLat, minLon})
	parallelDistance := GeoDistance(GeoPoint{minLat, minLon}, GeoPoint{minLat, maxLon})

	// remember, distance along a parallel is a change in longitude
	maxLonCell := int(math.Floor(parallelDistance / grid.CellWidth))
	// and distance along a meridian is a change in latitude
	maxLatCell := int(math.Floor(meridianDistance / grid.CellHeight))

	if minLatCoord < 0 {
		t.Errorf("minimum latitude grid coordinate should be non-negative")
	}
	if minLonCoord < 0 {
		t.Errorf("minimum longitude grid coordinate should be non-negative")
	}
	if maxLatCoord > maxLatCell {
		t.Errorf(fmt.Sprintf("maximum latitude grid coordinate should less than %d", maxLatCell))
	}
	if maxLonCoord > maxLonCell {
		t.Errorf(fmt.Sprintf("maximum longitude grid coordinate should less than %d", maxLonCell))
	}	
}
