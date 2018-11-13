package geo

import (
	"time"
	"fmt"
	"math"
)

type TraveledPath struct {
	Points []GeoPoint
	Date time.Time
	DistanceInMeters float64
}

type PathCluster struct {
	ReferencePath TraveledPath
	ContainedPaths []TraveledPath
	DistanceInMeters float64
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

func PathLengthInMeters(path []GeoPoint) float64 {

	pathLength := 0.0
	for i := 1; i < len(path); i++ {
		pathLength += GeoDistance(path[i], path[i-1])
	}
	return pathLength
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
	matchCount := 0
	for coord := range(inCoords1) {
		if !inCoords2[coord] {
			diffCount++
		} else {
			matchCount++
		}
	}
	for coord := range(inCoords2) {
		if !inCoords1[coord] {
			diffCount++
		} else {
			matchCount++
		}
	}
	return 1 - float64(diffCount) / float64(diffCount + matchCount)
}

func ClusterPaths(res float64, paths ...TraveledPath) []PathCluster {

	slimit := 0.35
	dlimit := 0.05
	clusters := make([]PathCluster, 0)

	for _, path := range(paths) {

		match := -1

		for j := 0; j < len(clusters) && match < 0; j++ {

			similarity := PathSimilarity(res, res, clusters[j].ReferencePath.Points, path.Points)
			
			refPathDistance := clusters[j].ReferencePath.DistanceInMeters
			avgDistance := (path.DistanceInMeters + refPathDistance) / 2
			distanceDelta := math.Abs(path.DistanceInMeters - refPathDistance) / avgDistance
			if similarity > slimit && distanceDelta < dlimit {
				match = j
				clusters[j].ContainedPaths = append(clusters[j].ContainedPaths, path)
			}
		}

		if match < 0 {
			cluster := PathCluster{
				ReferencePath: path,
				ContainedPaths: []TraveledPath{ path },
				DistanceInMeters: path.DistanceInMeters,
			}
			clusters = append(clusters, cluster)
		}
	}
	return clusters
}

func (c PathCluster) Print() {
	fmt.Printf("\tCluster Size: %d\n\tDistance: %f miles\n\tDate of Exemplar Path: %v",
		len(c.ContainedPaths),
		(c.DistanceInMeters * FeetPerMeter / FeetPerMile),
		c.ReferencePath.Date)
}
