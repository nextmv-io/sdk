package nextroute

import (
	"time"

	"github.com/nextmv-io/sdk/connect"
)

// NewSpeed creates a new speed.
func NewSpeed(distance float64, unit SpeedUnit) Speed {
	connect.Connect(con, &newSpeed)
	return newSpeed(distance, unit)
}

// Speed is the interface for a speed.
type Speed interface {
	// Value returns the speed in the specified unit.
	Value(unit SpeedUnit) float64
}

// NewSpeedUnit returns a new speed unit.
func NewSpeedUnit(distance DistanceUnit, duration time.Duration) SpeedUnit {
	connect.Connect(con, &newSpeedUnit)
	return newSpeedUnit(distance, duration)
}

// SpeedUnit represents a unit of speed.
type SpeedUnit interface {
	// DistanceUnit returns the distance unit of the speed unit.
	DistanceUnit() DistanceUnit
	// Duration returns the duration of the speed unit.
	Duration() time.Duration
}

// KilometersPerHour is a speed unit of kilometers per hour.
var KilometersPerHour = NewSpeedUnit(Kilometers, time.Hour)

// MilesPerHour is a speed unit of miles per hour.
var MilesPerHour = NewSpeedUnit(Miles, time.Hour)

// MetersPerSecond is a speed unit of meters per second.
var MetersPerSecond = NewSpeedUnit(Meters, time.Second)
