package geo

import (
	"math"
)

const EarthRadiusInMeters float64 = 6371.0088 * 1000
// a change in one degree of latitude along a meridian is about 111.32km
const LatitudeDegreesToMeters float64 = 111.32 * 1000

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
