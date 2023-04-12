package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// NewLocation creates a new Location. An error is returned if the longitude is
// not between (-180, 180) or the latitude is not between (-90, 90).
func NewLocation(longitude, latitude float64) (Location, error) {
	connect.Connect(con, &newLocation)
	return newLocation(longitude, latitude)
}

// NewInvalidLocation creates a new invalid Location. Longitude and latitude
// are not important.
func NewInvalidLocation() Location {
	connect.Connect(con, &newInvalidLocation)
	return newInvalidLocation()
}

// Location represents a physical location on the earth.
type Location interface {
	// Longitude returns the longitude of the location.
	Longitude() float64
	// Latitude returns the latitude of the location.
	Latitude() float64
	// Equals returns true if the location is equal to the location given as an
	// argument.
	Equals(Location) bool
	// IsValid returns true if the location is valid. A location is valid if
	// the bounds of the longitude and latitude are correct.
	IsValid() bool
}

// Locations is a slice of Location.
type Locations []Location

// Unique returns a new slice of Locations with unique locations.
func Unique(l Locations) Locations {
	connect.Connect(con, &unique)
	return unique(l)
}

// Centroid returns the centroid of the locations. If locations is empty, the
// centroid will be an invalid location.
func Centroid(l Locations) (Location, error) {
	connect.Connect(con, &centroid)
	return centroid(l)
}
