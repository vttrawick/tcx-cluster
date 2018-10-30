package geo

import(
	"testing"
	"reflect"
	"math"
)

func TestGridBoundary(t *testing.T) {
	// first create some test cases

	path1 := []GeoPoint{
		GeoPoint{42.365592, -71.103875},
		GeoPoint{42.364776, -71.110749},
		GeoPoint{42.364237, -71.116022},
		GeoPoint{42.361439, -71.115968},
		GeoPoint{42.362285, -71.113515},
		GeoPoint{42.365115, -71.104975},
	}

	// find the min / max lat long for one path
	boundary1 := GridBoundary(path1)

	// boundary should be:
	// {
	//   maxLat: 42.365592
	//   minLat: 42.361439
	//   maxLon: -71.103875
	//   minLon: -71.116022
	// }

	t1 := GeoRect{42.361439, 42.365592, -71.116022, -71.103875}
	if !reflect.DeepEqual(boundary1, t1) {
		t.Errorf("incorrect boundary structure for 1-path example")
	}

	path2 := []GeoPoint{
		GeoPoint{42.365218, -71.104578},
		GeoPoint{42.362285, -71.113429},
		GeoPoint{42.360644, -71.113225},
		GeoPoint{42.360192, -71.112774},
		GeoPoint{42.364037, -71.108278},
		GeoPoint{42.365223, -71.104852},
	}

	path3 := []GeoPoint{
		GeoPoint{42.365127, -71.103168},
		GeoPoint{42.360831, -71.096162},
		GeoPoint{42.359000, -71.100175},
		GeoPoint{42.361164, -71.103930},
		GeoPoint{42.365516, -71.103951},
		GeoPoint{42.364667, -71.102567},
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

	t2 := GeoRect{42.359000, 42.365592, -71.116022, -71.096162}

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
