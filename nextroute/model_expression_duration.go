package nextroute

import (
	"time"

	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/common"
)

// NewDurationExpression creates a new duration expression.
func NewDurationExpression(
	expression ModelExpression,
	multiplier time.Duration,
) DurationExpression {
	connect.Connect(con, &newDurationExpression)
	return newDurationExpression(expression, multiplier)
}

// NewTravelDurationExpression creates a new travel duration expression.
func NewTravelDurationExpression(
	distanceExpression DistanceExpression,
	speed common.Speed,
) TravelDurationExpression {
	connect.Connect(con, &newTravelDurationExpression)
	return newTravelDurationExpression(distanceExpression, speed)
}

// NewVehicleTypeDurationExpression creates a new duration expression.
func NewVehicleTypeDurationExpression(
	name string,
	duration time.Duration,
) VehicleTypeDurationExpression {
	connect.Connect(con, &newVehicleTypeDurationExpression)
	return newVehicleTypeDurationExpression(name, duration)
}

// NewConstantDurationExpression creates a new constant duration expression.
func NewConstantDurationExpression(
	name string,
	duration time.Duration,
) DurationExpression {
	connect.Connect(con, &newConstantDurationExpression)
	return newConstantDurationExpression(name, duration)
}

type VehicleTypeDurationExpression interface {
	ModelExpression
	// Duration returns the duration for the given vehicle type
	Duration(ModelVehicleType) time.Duration
	// SetDuration sets the duration for the given vehicle type.
	SetDuration(ModelVehicleType, time.Duration)
}

// DurationExpression is an expression that returns a duration.
type DurationExpression interface {
	ModelExpression
	// Duration returns the duration for the given vehicle type, start and
	// end stop.
	Duration(ModelVehicleType, ModelStop, ModelStop) time.Duration
}

// DistanceExpression is an expression that returns a distance.
type DistanceExpression interface {
	ModelExpression
	// Distance returns the distance for the given vehicle type, start and
	// end stop.
	Distance(ModelVehicleType, ModelStop, ModelStop) common.Distance
}

// TravelDurationExpression is an expression that returns a duration based on
// a distance and a speed.
type TravelDurationExpression interface {
	DurationExpression
	// DistanceExpression returns the distance expression.
	DistanceExpression() DistanceExpression
	// Speed returns the speed.
	Speed() common.Speed
}
