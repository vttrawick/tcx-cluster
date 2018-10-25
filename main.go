package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/llehouerou/go-tcx"
	"github.com/vttrawick/tcx-cluster/geo"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"runtime"
	"math"
)

func main() {

	// flags
	var tcxdir string
	flag.StringVar(&tcxdir,
		"tcxdir",
		"",
		"directory containing tcx files to cluster")

	flag.Parse()

	files, dirErr := ioutil.ReadDir(tcxdir)
	fmt.Printf("there are %d files in this dir\n", len(files))
	if dirErr != nil {
		log.Fatal(dirErr)
	}

	// paths list of paths to cluster
	paths := make([]*tcx.Tcx, 0, len(files))

//	for i := range files {
	for i := 0; i < 10; i++ {
		info := files[i]
		tcxMatch, _ := regexp.MatchString(`\.tcx$`, info.Name())
		
		if !info.IsDir() && tcxMatch {
			contents, fileErr := ioutil.ReadFile(filepath.Join(tcxdir, info.Name()))
			if fileErr != nil {
				log.Fatal(fileErr)
			}
			activity, parseErr := tcx.Parse(bytes.NewBuffer(contents))
			if parseErr != nil {
				fmt.Printf("error parsing file %v\n", info.Name())
			}
			paths = append(paths, activity)
		}
	}
	ClusterPaths(paths)
}

func ClusterPaths(paths []*tcx.Tcx) [][]uint16 {

	// the cluster list is an array of arrays of indexes into the path list
	clusters := make([][]uint16, 0, 512)

	for i := range(paths) {

		// compare path to the "reference" path for each existing cluster
		match := -1
		runtime.Breakpoint()
		for j := 0; j < len(clusters) && match < 0; j++ {

			clusterExampleIndex := clusters[j][0]
			score := differenceScore(paths[i], paths[clusterExampleIndex])

			if score > 0.95 {
				match = j
				clusters[j] = append(clusters[j], uint16(i))
			}
		}
		// if no match, create a new cluster
		if match < 0 {
			cluster := make([]uint16, 0, 64)
			cluster = append(cluster, uint16(i))

			clusters = append(clusters, cluster)
		}
	}
	return clusters
}

func gridBoundary(paths ...[]tcx.Trackpoint) gridRect {

	r := gridRect{
		minLat: 91
		maxLat: -91
		minLon: 181
		maxLon: -181
	}
	for _, path := range(paths) {
		for _, pt := range(path) {
			if pt.LatitudeInDegrees < minLat {
				r.minLat = pt.LatitudeInDegrees
			}
			if pt.LatitudeInDegrees > maxLat {
				r.maxLat = pt.LatitudeInDegrees
			}
			if pt.LongitudeInDegrees < minLon {
				r.minLon = pt.LongitudeInDegrees
			}
			if pt.LongitudeInDegrees > maxLon {
				r.maxLon = pt.LongitudeInDegrees
			}
		}
	}

	return r
}

// computes the path difference score between two paths
func differenceScore(path1 *tcx.Tcx, path2 *tcx.Tcx) float64 {

	// first flatten the paths to get rid of 'activities' and 'laps' concepts
	var fpath1 := flattenTcxPts(path1)
	var fpath2 := flattenTcxPts(path2)

	// then find the bounding rectangles for each path
	r1 := gridBoundary(fpath1)
	r2 := gridBoundary(fpath2)

	// if there is no possible overlap in the paths, skip the rest
	if (!gridOverlap(r1, r2)) {
		return 0
	}

	// create an outer boundary for all paths
	merged := mergeRect(r1, r2)
	
	gridmap := gridOverlay(merged)

	gb := gridBoundary(fpath1, fpath2)
	fmt.Printf("grid boundaries: %v", gb)

	runtime.Breakpoint()

	return float64(len(fpath1) - len(fpath2))
}

func gridOverlay(bounds gridRect, paths ...[]tcx.Trackpoint) []map[uint16]uint16 {

	d1 := eastWestDistanceInMeters(bounds)
	d2 := northSouthDistanceInMeters(bounds)

	// all calculations done in meters
	cellWidth = 7

	math.Floor(d1 / cellWidth) + 1
	math.Floor(d2 / cellWidth) + 1

	// represent the grid as an array of integers,
	// with 0 in the northwestern-most corner
	// and n in the southeasternmost corner
	// each rect should be up to 5x5m
	
	// each trackpoint in the set of paths will
	// then be indexed to its place in the grid

	// what does this allow for?
	
	
}

func gridOverlap(r1 gridRect, r2 gridRect) bool {
	// longitude always decreases to the west,
	// until you get to the anti-meridian. If your data is from
	// Taveuni or somewhere like that this will just be wrong.
	if r1.minLon >= r2.maxLon || r2.minLon >= r1.maxLon {
		return false
	}

	// latitude will get pretty weird around the poles
	// this formula will technically work,
	// but the rectangle may be totally misrepresented
	if r1.minLat >= r2.maxLat || r2.minLat >= r1.maxLat {
		return false
	}

	return true
}

func mergeGridRect(rlist ...gridRect) {

	merged := gridRect{
		minLat: 91
		maxLat: -91
		minLon: 181
		maxLon: -181
	}

	for r := range(rlist) {
		if r.minLat < minLat {
			merged.minLat = pt.LatitudeInDegrees
		}
		if r.maxLat > maxLat {
			merged.maxLat = pt.LatitudeInDegrees
		}
		if r.minLon < minLon {
			merged.minLon = pt.LongitudeInDegrees
		}
		if r.maxLon > maxLon {
			merged.maxLon = pt.LongitudeInDegrees
		}
	}
	return merged
}

func trackPointDistance(p1 tcx.Trackpoint, p2 tcx.Trackpoint) float64 {

	var dphi, dlam, avgphi float64
	// north / south degree distance
	dphi = (p1.LatitudeInDegrees - p2.LatitudeInDegrees)
	// east / west degree distance
	dlambda = p1.LongitudeInDegrees - p2.LongitudeInDegrees
	// average latitude
	avgphi = math.Abs((p1.LongitudeInDegrees + p2.LongitudeInDegrees) / 2)

	// convert to radians
	avgphi *= math.Pi / 180
	dphi *= math.Pi / 180
	dlambda *= math.Pi / 180
	// then finally adjust the east-west distance by the latitude
	dlam *= math.Cos(avgphi)

	// then just do the pythagorean theorem
	return earthRadiusInMeters * math.Sqrt((dphi * dphi) + (dlambda * dlambda))
}

func eastWestDistanceInMeters(bounds gridRect) float64 {

	// find the latitude used in the circumfrence calculation
	// probably doesn't make much of a difference,
	// but we should use the longer of the two arcs
	// at the latitude closer to the equator
	var latInRads
	// if we're entirely in the southern hemisphere
	if math.Abs(bounds.minLat) > math.Abs(bounds.maxLat) {
		latInRads = math.Abs(bounds.maxLat) * math.Pi / 180
	} else {
		latInRads = math.Abs(bounds.minLat) * math.Pi / 180
	}

	dLon := (bounds.maxLon - bounds.minLon) * math.Pi / 180
	return earthRadiusInMeters * math.Cos(latInRads) * dLon
}

func northSouthDistanceInMeters(bounds gridRect) float64 {

	// lat is in degrees, convert to radians
	dLat := (bounds.maxLat - bounds.minLat) * math.Pi / 180
	return earthRadiusInMeters * dLat
}

// take the tcx and put all the trackpoints in there
func flattenTcxPts(activity *tcx.Tcx) []tcx.Trackpoint {

	flat := make([]tcx.Trackpoint, 0, 512)

	laps := activity.Activities[0].Laps

	for _, lap := range(laps) {
		for _, point := range(lap.Track) {
			flat = append(flat, point)
		}
	}
	return flat
}
