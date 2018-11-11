package geo

import(
	"testing"
	"math"
	"reflect"
)

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

func TestPathSimilarity(t *testing.T) {

	// test that the same path has a similarity of 1
	// and two paths with different boundaries have a similarity of 0
	// and two paths that are somewhat similar have a score between those two numbers
	score := PathSimilarity(10.0, 10.0, path1, path1)
	if score != 1 {
		t.Errorf("the exact same path should return a score of 1")
	}

	score = PathSimilarity(10.0, 10.0, path1, path3)
	if score != 0 {
		t.Errorf("two paths with non-overlapping boundaries should return a score of 0")
	}

	score = PathSimilarity(10.0, 10.0, path1, path2)
	if !(score > 0 && score < 1) {
		t.Errorf("two somewhat similar paths should return a score between 0 and 1")
	}
}
