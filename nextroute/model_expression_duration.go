package nextroute

import (
	"time"

	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/common"
)

// NewDurationExpression creates a new duration expression.
func NewDurationExpression(
	name string,
	expression ModelExpression,
	unit common.DurationUnit,
) DurationExpression {
	connect.Connect(con, &newDurationExpression)
	return newDurationExpression(name, expression, unit)
}

// NewScaledDurationExpression creates a new scaled duration expression.
func NewScaledDurationExpression(
	expression DurationExpression,
	scale float64,
) DurationExpression {
	connect.Connect(con, &newScaledDurationExpression)
	return newScaledDurationExpression(expression, scale)
}

// NewTravelDurationExpression creates a new travel duration expression.
func NewTravelDurationExpression(
	distanceExpression DistanceExpression,
	speed common.Speed,
) TravelDurationExpression {
	connect.Connect(con, &newTravelDurationExpression)
	return newTravelDurationExpression(distanceExpression, speed)
}

// NewStopDurationExpression creates a new duration expression.
func NewStopDurationExpression(
	name string,
	duration time.Duration,
) StopDurationExpression {
	connect.Connect(con, &newStopDurationExpression)
	return newStopDurationExpression(name, duration)
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

// StopDurationExpression is a ModelExpression that returns a duration per stop
// and allows to set the duration per stop.
type StopDurationExpression interface {
	DurationExpression
	// SetDuration sets the duration for the given stop.
	SetDuration(ModelStop, time.Duration)
}

// VehicleTypeDurationExpression is a ModelExpression that returns a duration
// per vehicle type and allows to set the duration per vehicle type.
type VehicleTypeDurationExpression interface {
	DurationExpression
	VehicleTypeExpression
	// SetDuration sets the duration for the given vehicle type.
	SetDuration(ModelVehicleType, time.Duration)
	DurationForVehicleType(ModelVehicleType) time.Duration
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

// VehicleTypeDistanceExpression is an expression that returns a distance per
// vehicle type and allows to set the duration per vehicle.
type VehicleTypeDistanceExpression interface {
	DistanceExpression
	VehicleTypeExpression

	SetDistance(ModelVehicleType, common.Distance)
	DistanceForVehicleType(ModelVehicleType) common.Distance
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
