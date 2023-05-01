package nextroute

import "github.com/nextmv-io/sdk/connect"

// NewVehiclesObjective returns a new VehiclesObjective that uses the number of
// vehicles as an objective. Each vehicle that is not empty is scored by the
// given expression. A vehicle is empty if it has no stops assigned to it
// (except for the first and last visit).
func NewVehiclesObjective(expression VehicleTypeExpression) VehiclesObjective {
	connect.Connect(con, &newVehiclesObjective)
	return newVehiclesObjective(expression)
}

// VehiclesObjective is an objective that uses the number of vehicles as an
// objective. Each vehicle that is not empty is scored by the given expression.
// A vehicle is empty if it has no stops assigned to it (except for the first
// and last visit).
type VehiclesObjective interface {
	ModelObjective
}
