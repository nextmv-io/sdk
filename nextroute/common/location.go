package common

import (
	"fmt"
)

// NewLocation creates a new Location. Will panic if longitude or latitude are
// out of range. Longitude must be between -180 and 180. Latitude must be
// between -90 and 90.
func NewLocation(
	longitude float64,
	latitude float64,
) Location {
	if longitude < -180 || longitude > 180 {
		panic("longitude must be between -180 and 180")
	}
	if latitude < -90 || latitude > 90 {
		panic("latitude must be between -90 and 90")
	}
	return Location{
		longitude: longitude,
		latitude:  latitude,
		valid:     true,
	}
}

// NewInvalidLocation creates a new invalid Location. Longitude and latitude
// are not important.
func NewInvalidLocation() Location {
	return Location{
		valid: false,
	}
}

// Locations is a slice of Location.
type Locations []Location

// Unique returns a new slice of Locations with unique locations.
func (l Locations) Unique() Locations {
	unique := make(map[string]Location)
	for _, location := range l {
		// TODO: in Go 1.20 we don't need to use fmt.Sprintf here.
		// This can simply become unique[location] = struct{}{}
		unique[fmt.Sprintf("%v", location)] = location
	}
	result := make(Locations, 0, len(unique))
	for _, location := range unique {
		result = append(result, location)
	}
	return result
}

// Centroid returns the centroid of the locations. If locations is empty, the
// centroid will be (0, 0).
func (l Locations) Centroid() Location {
	if len(l) == 0 {
		return NewLocation(0, 0)
	}
	lat := 0.0
	lon := 0.0
	for _, location := range l {
		lat += location.latitude
		lon += location.longitude
	}
	return NewLocation(lon/float64(len(l)), lat/float64(len(l)))
}

// Location represents a location on the earth.
type Location struct {
	longitude float64
	latitude  float64
	valid     bool
}

func (l Location) String() string {
	return fmt.Sprintf(
		"{lat: %v,lon: %v}",
		l.latitude,
		l.longitude,
	)
}

func (l Location) Longitude() float64 {
	return l.longitude
}

func (l Location) Latitude() float64 {
	return l.latitude
}

func (l Location) Equals(other Location) bool {
	return l.longitude == other.longitude && l.latitude == other.latitude
}

func (l Location) IsValid() bool {
	return l.valid
}
