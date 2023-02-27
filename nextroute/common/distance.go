package common

// DistanceUnit is the unit of distance.
type DistanceUnit int

// NewDistance returns a new distance.
func NewDistance(
	distance float64,
	unit DistanceUnit,
) Distance {
	meters := distance
	switch unit {
	case Kilometers:
		meters *= factorKilometersToMeters
	case Miles:
		meters *= factorMilesToMeters
	}

	return &distanceImpl{
		meters: meters,
		unit:   unit,
	}
}

const (
	// Kilometers is 1000 meters.
	Kilometers DistanceUnit = iota
	// Meters is the distance travelled by light in a vacuum in
	// 1/299,792,458 seconds.
	Meters
	// Miles is 1609.34 meters.
	Miles
)

const (
	factorMetersToKilometers = 0.001
	factorMetersToMiles      = 0.000621371
)

const (
	factorKilometersToMeters = 1000
	factorMilesToMeters      = 1609.34
)

// Distance is the interface for a distance.
type Distance interface {
	// Value returns the distance in the specified unit.
	Value(unit DistanceUnit) float64
}

type distanceImpl struct {
	meters float64
	unit   DistanceUnit
}

func (d *distanceImpl) Value(unit DistanceUnit) float64 {
	returnValue := d.meters
	switch unit {
	case Kilometers:
		returnValue *= factorMetersToKilometers
	case Miles:
		returnValue *= factorMetersToMiles
	}
	return returnValue
}
