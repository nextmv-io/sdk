package nextroute

import (
	"fmt"

	"github.com/nextmv-io/sdk/connect"
)

// NewDistance returns a new distance from the given meters.
func NewDistance(meters float64, unit DistanceUnit) Distance {
	connect.Connect(con, &newDistance)
	return newDistance(meters, unit)
}

// Haversine calculates the distance between two locations using the
// Haversine formula. Haversine is a good approximation for short
// distances (up to a few hundred kilometers).
func Haversine(from, to Location) (Distance, error) {
	connect.Connect(con, &haversine)
	return haversine(from, to)
}

// Distance is the distance between two points.
type Distance interface {
	// Value returns the distance in the specified unit.
	Value(unit DistanceUnit) float64
}

// DistanceUnit is the unit of distance.
type DistanceUnit int

const (
	// Kilometers represents kilometers.
	Kilometers DistanceUnit = iota
	// Meters represents meters.
	Meters
	// Miles represents miles.
	Miles
)

// String returns the string representation of the distance unit.
func (d DistanceUnit) String() string {
	switch d {
	case Kilometers:
		return "kilometers"
	case Meters:
		return "meters"
	case Miles:
		return "miles"
	}
	return fmt.Sprintf("unknown distance unit %v", int(d))
}
