package geo

// given a lat-long rectangle, create a grid where each box is nxn meters
// if the rectangle spans a long north-south distance there may be more boxes
// at the top than at the bottom

// the grid is a slice where each element is a slice of boxes
// each element in the primary array represents a swath latitude with height boxHeight

const EarthRadiusInMeters float64 = 6371.009 * 1000

// a coordinate on the earth represented by lat / lon degrees
type GeoPoint struct {
	LatitudeInDegrees, LongitudeInDegrees float64
}

// the minimum area bounded by the four lat / lon lines
// all values in degrees
type GeoRect struct {
	MinLat, MaxLat, MinLon, MaxLon float64
}

type GeoGrid struct {
	Boundary GeoRect
	SubHeight, SubWidth int
	// map from grid coordinate to point list
}

func GridBoundary(paths ...[]GeoPoint) GeoRect {

	r := GeoRect{
		MinLat: 91,
		MaxLat: -91,
		MinLon: 181,
		MaxLon: -181,
	}
	for _, path := range(paths) {
		for _, pt := range(path) {
			if pt.LatitudeInDegrees < r.MinLat {
				r.MinLat = pt.LatitudeInDegrees
			}
			if pt.LatitudeInDegrees > r.MaxLat {
				r.MaxLat = pt.LatitudeInDegrees
			}
			if pt.LongitudeInDegrees < r.MinLon {
				r.MinLon = pt.LongitudeInDegrees
			}
			if pt.LongitudeInDegrees > r.MaxLon {
				r.MaxLon = pt.LongitudeInDegrees
			}
		}
	}
	return r
}

