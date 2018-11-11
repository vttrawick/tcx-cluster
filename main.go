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
)

func main() {

	// command line flags
	var tcxdir string
	flag.StringVar(&tcxdir,
		"tcxdir",
		"",
		"directory containing tcx files to cluster")

	flag.Parse()

	tcxlist := LoadTCXDir(tcxdir)
	paths := tcx2path(tcxlist)
	clusters := geo.ClusterPaths(7.0, paths...)
	fmt.Printf("%d paths have been filtered into %d clusters\n", len(paths), len(clusters))
}

func LoadTCXDir(tcxdir string) []*tcx.Tcx {

	files, dirErr := ioutil.ReadDir(tcxdir)
	fmt.Printf("there are %d files in this dir\n", len(files))
	if dirErr != nil {
		log.Fatal(dirErr)
	}

	tcxlist := make([]*tcx.Tcx, 0, len(files))

	for i := range files {
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
			tcxlist = append(tcxlist, activity)
		}
	}
	return tcxlist
}

// format the cumbersome tcx structs into the simpler TraveledPath struct
func tcx2path(tcxlist []*tcx.Tcx) []geo.TraveledPath {

	pathlist := make([]geo.TraveledPath, 0)

	for _, tcxdata := range(tcxlist) {

		activity := tcxdata.Activities[0]

		pts, distance := extractPtsAndDistance(activity)
		// skip over treadmill activities and things without geo data
		if len(pts) > 0 {
			pathlist = append(pathlist, geo.TraveledPath{
				Date: tcxdata.Activities[0].ID,
				DistanceInMeters: distance,
				Points: pts,
			})
		}
	}
	return pathlist;
}

// take the tcx and put all the trackpoints in there
func extractPtsAndDistance(activity tcx.Activity) ([]geo.GeoPoint, float64) {

	pts := make([]geo.GeoPoint, 0, 512)
	distance := 0.0

	for _, lap := range(activity.Laps) {
		distance += lap.DistanceInMeters
		for _, point := range(lap.Track) {
			pts = append(pts, geo.GeoPoint{
				LatitudeInDegrees: point.LatitudeInDegrees,
				LongitudeInDegrees: point.LongitudeInDegrees,
			})
		}
	}
	return pts, distance
}
