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
	GeoPoint{42.360352, -71.102460}, 
	GeoPoint{42.361864, -71.100765},
	GeoPoint{42.364667, -71.102567},
}
	
func TestPathBoundary(t *testing.T) {

	// find the min / max lat long for one path
	boundary1 := PathBoundary(path1)

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
	boundary2 := PathBoundary(path1, path2, path3)

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

func TestPathLengthInMeters(t *testing.T) {

	l1 := PathLengthInMeters(path1)
	if math.Abs(l1 - 2310) > 10 {
		t.Errorf("path1 length calculated incorrectly")
	}

	l2 := PathLengthInMeters(path2)
	if math.Abs(l2 - 1920) > 10 {
		t.Errorf("path2 length calculated incorrectly")
	}
}

func TestMakeGrid(t *testing.T) {

	width := 7.0
	height := 7.0

	boundary := GeoRect{
		NorthWest: GeoPoint{42.365592, -71.116022},
		SouthEast: GeoPoint{42.361439, -71.103875},
	}
	
	g1 := MakeGrid(width, height, boundary)

	if (g1.CellWidth != width) {
		t.Errorf("grid cell width differs from input")
	}
	if (g1.CellHeight != height) {
		t.Errorf("grid cell height differs from input")
	}

	if !reflect.DeepEqual(g1.Boundary, boundary) {
		t.Errorf("grid boundary does not match expected")
	}	
}

func TestMapPoint(t *testing.T) {
	grid := MakeGrid(7.0, 7.0, PathBoundary(path1))

	c1 := grid.MapPoint(path1[0])
	expected := cellCoord("0_142")
	if c1 != expected {
		t.Errorf(fmt.Sprintf("point %v mapped incorrectly: was %v but expected %v",
			path1[0], c1, expected))
	}

	c2 := grid.MapPoint(path1[2])
	expected = cellCoord("21_0")
	if c2 != expected {
		t.Errorf(fmt.Sprintf("point %v mapped incorrectly: was %v but expected %v",
			path1[0], c2, expected))
	}
	
	c3 := grid.MapPoint(path1[4])
	expected = cellCoord("52_29")
	if c3 != expected {
		t.Errorf(fmt.Sprintf("point %v mapped incorrectly: was %v but expected %v",
			path1[0], c3, expected))
	}
}

func TestMapPath(t *testing.T) {

	grid := MakeGrid(7.0, 7.0, PathBoundary(path1))
	coords := grid.MapPath(path1)

	if len(coords) != len(path1) {
		t.Errorf("length mismatch between path and grid map")
	}
		
	// ensure all points are within the grid
	minLatCoord := math.MaxInt64
	maxLatCoord := -1
	minLonCoord := math.MaxInt64
	maxLonCoord := -1
	for _, coord := range(coords) {
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

func TestMergeGeoRect(t *testing.T) {

	r1 := PathBoundary(path1)

	m1 := MergeGeoRect(r1, r1)
	if !reflect.DeepEqual(m1, r1) {
		t.Errorf("merging two of the same GeoRect should yield the original")
	}

	r2 := PathBoundary(path2)
	m2 := MergeGeoRect(r1, r2)
	t2 := GeoRect{
		NorthWest: GeoPoint{42.365592, -71.116022},
		SouthEast: GeoPoint{42.360192, -71.103875},
	}
	if !reflect.DeepEqual(m2, t2) {
		t.Errorf(fmt.Sprintf("merging of two overlapping rects failed, got %v but expected %v",
			t2, m2))
	}

	r3 := PathBoundary(path3)
	m3 := MergeGeoRect(r1, r3)
	t3 := GeoRect{
		NorthWest: GeoPoint{42.365592, -71.116022},
		SouthEast: GeoPoint{42.359000, -71.096162},
	}
	if !reflect.DeepEqual(m3, t3) {
		t.Errorf(fmt.Sprintf("merging of two non-overlapping rects failed, got %v but expected %v",
			t3, m3))
	}
}

func TestOverlaps(t *testing.T) {

	r1 := PathBoundary(path1)
	r2 := PathBoundary(path2)
	
	if r1.Overlaps(r2) != true {
		t.Errorf("false negative error on overlap detection")
	}

	r3 := PathBoundary(path3)
	if r1.Overlaps(r3) != false {
		t.Errorf("false positive error on overlap detection")
	}
}

func TestPathSimilarity(t *testing.T) {
	score := PathSimilarity(50.0, 50.0, path1, path2)
	if score < 0 {
		t.Errorf("similarity score should be between 0 and 1")
	}
}
