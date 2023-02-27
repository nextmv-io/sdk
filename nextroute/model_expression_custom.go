package nextroute

import "github.com/nextmv-io/sdk/connect"

// NewConstantExpression returns an expression that always returns the same
// value.
func NewConstantExpression(
	name string,
	value float64,
) ConstantExpression {
	connect.Connect(con, &newConstantExpression)
	return newConstantExpression(name, value)
}

// NewFromStopExpression returns an expression whose value is based on the
// from stop in Value method. Expression values are calculated by calling the
// expression's Value method which takes a vehicle type, a from stop and a to
// stop as arguments.
func NewFromStopExpression(
	name string,
	defaultValue float64,
) FromStopExpression {
	connect.Connect(con, &newFromStopExpression)
	return newFromStopExpression(name, defaultValue)
}

// NewStopExpression returns an expression whose value is based on the
// to stop in Value method. Expression values are calculated by calling the
// expression's Value method which takes a vehicle type, a from stop and a to
// stop as arguments.
func NewStopExpression(
	name string,
	defaultValue float64,
) StopExpression {
	connect.Connect(con, &newStopExpression)
	return newStopExpression(name, defaultValue)
}

// NewVehicleTypeExpression returns a VehicleTypeExpression whose value is
// based on the  vehicle type in Value method. Expression values are calculated
// by calling the expression's Value method which takes a vehicle type, a from
// stop and a to stop as arguments.
func NewVehicleTypeExpression(
	name string,
	defaultValue float64,
) VehicleTypeExpression {
	connect.Connect(con, &newVehicleTypeExpression)
	return newVehicleTypeExpression(name, defaultValue)
}

// NewFromToExpression returns an expression whose value is based on the
// from and to stops in Value method. Expression values are calculated by
// calling the expression's Value method which takes a vehicle type, a
// from stop and a to stop as arguments.
func NewFromToExpression(
	name string,
	defaultValue float64,
) FromToExpression {
	connect.Connect(con, &newFromToExpression)
	return newFromToExpression(name, defaultValue)
}

// NewVehicleTypeFromToExpression returns an expression whose value is based on
// the vehicle type, from and to stops in Value method. Expression values are
// calculated by calling the expression's Value method which takes a vehicle
// type, a from stop and a to stop as arguments.
func NewVehicleTypeFromToExpression(
	name string,
	defaultValue float64,
) VehicleFromToExpression {
	connect.Connect(con, &newVehicleTypeFromToExpression)
	return newVehicleTypeFromToExpression(name, defaultValue)
}

// ConstantExpression is an expression that always returns the same value.
type ConstantExpression interface {
	ModelExpression

	// SetValue sets the value of the expression.
	SetValue(value float64)
}

// DefaultExpression is an expression that has a default value if no other
// values are defined.
type DefaultExpression interface {
	ModelExpression
	// DefaultValue returns the default value of the expression.
	DefaultValue() float64
}

// FromStopExpression is an expression that has a value for each from stop.
type FromStopExpression interface {
	DefaultExpression

	// SetValue sets the value of the expression for the given from stop.
	SetValue(
		stop ModelStop,
		value float64,
	)
}

// StopExpression is an expression that has a value for each to stop.
type StopExpression interface {
	DefaultExpression

	// SetValue sets the value of the expression for the given to stop.
	SetValue(
		stop ModelStop,
		value float64,
	)
}

// VehicleTypeExpression is an expression that has a value for each vehicle
// type.
type VehicleTypeExpression interface {
	DefaultExpression

	// SetValue sets the value of the expression for the given vehicle type.
	SetValue(
		vehicle ModelVehicleType,
		value float64,
	)
}

// FromToExpression is an expression that has a value for each combination
// of from and to stop.
type FromToExpression interface {
	DefaultExpression

	// SetValue sets the value of the expression for the given
	// from and to stops.
	SetValue(
		from ModelStop,
		to ModelStop,
		value float64,
	)
}

// VehicleFromToExpression is an expression that has a value for each
// combination of vehicle type, from and to stop.
type VehicleFromToExpression interface {
	DefaultExpression

	// SetValue sets the value of the expression for the given vehicle type,
	// from and to stops.
	SetValue(
		vehicle ModelVehicleType,
		from ModelStop,
		to ModelStop,
		value float64,
	)
}
