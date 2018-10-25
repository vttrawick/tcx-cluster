package geo

// given a lat-long rectangle, create a grid where each box is nxn meters
// if the rectangle spans a long north-south distance there may be more boxes
// at the top than at the bottom

// the grid is a slice where each element is a slice of boxes
// each element in the primary array represents a swath latitude with height boxHeight

const EarthRadiusInMeters float64 = 6371.009 * 1000

// a coordinate on the earth
// represented by lat / lon degrees
type GeoPoint struct {
	lat, lon float64
}

// the minimum area bounded by the four lat / lon lines
type GeoRect struct {
	minLat, maxLat, minLon, maxLon float64
}

func GridBoundary(paths ...[]tcx.Trackpoint) gridRect {

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

func Create(g GridRect) (*GridMap) {
	
}
