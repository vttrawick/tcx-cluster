package geo

import (
	"testing"
	"math"
)

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
