package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// NewComposedPerVehicleTypeExpression returns a new
// ComposedPerVehicleTypeExpression which is a composed expression that uses
// an expression for each vehicle type.
func NewComposedPerVehicleTypeExpression(
	defaultExpression ModelExpression,
) ComposedPerVehicleTypeExpression {
	connect.Connect(con, &newComposedPerVehicleTypeExpression)
	return newComposedPerVehicleTypeExpression(defaultExpression)
}

// ComposedPerVehicleTypeExpression is an expression that uses an expression for
// each vehicle type.
type ComposedPerVehicleTypeExpression interface {
	ModelExpression

	// DefaultExpression returns the default expression that is used if no
	// expression is defined for a specific vehicle type.
	DefaultExpression() ModelExpression

	// Get returns the expression that is defined for the given vehicle type. If
	// no expression is defined for the given vehicle type, the default
	// expression is returned.
	Get(vehicleType ModelVehicleType) ModelExpression
	// Set sets the expression for the given vehicle type.
	Set(vehicleType ModelVehicleType, expression ModelExpression)
}
