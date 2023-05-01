package nextroute

import "github.com/nextmv-io/sdk/connect"

// AttributesConstraint is a constraint that limits the vehicles a plan unit
// can be added to. The Attribute constraint configures compatibility
// attributes for stops and vehicles separately. This is done by specifying
// a list of attributes for stops and vehicles, respectively. Stops that
// have configured attributes are only compatible with vehicles that match
// at least one of them. Stops that do not have any specified attributes are
// compatible with any vehicle. Vehicles that do not have any specified
// attributes are only compatible with stops without attributes.
type AttributesConstraint interface {
	ModelConstraint

	// SetStopAttributes sets the attributes for the given stop. The attributes
	// are specified as a list of strings. The attributes are not interpreted
	// in any way. They are only used to determine compatibility between stops
	// and vehicle types.
	SetStopAttributes(
		stop ModelStop,
		stopAttributes []string,
	)
	// SetVehicleTypeAttributes sets the attributes for the given vehicle type.
	// The attributes are specified as a list of strings. The attributes are not
	// interpreted in any way. They are only used to determine compatibility
	// between stops and vehicle types.
	SetVehicleTypeAttributes(
		vehicle ModelVehicleType,
		vehicleAttributes []string,
	)
	// StopAttributes returns the attributes for the given stop. The attributes
	// are specified as a list of strings.
	StopAttributes(stop ModelStop) []string

	// VehicleTypeAttributes returns the attributes for the given vehicle type.
	// The attributes are specified as a list of strings.
	VehicleTypeAttributes(vehicle ModelVehicleType) []string
}

// NewAttributesConstraint creates a new attributes constraint. The constraint
// needs to be added to the model to be taken into account.
func NewAttributesConstraint() (AttributesConstraint, error) {
	connect.Connect(con, &newAttributesConstraint)
	return newAttributesConstraint()
}
