package nextroute

import "github.com/nextmv-io/sdk/connect"

// NewVehiclesDurationObjective returns a new VehiclesDurationObjective that
// uses the vehicle duration as an objective.
func NewVehiclesDurationObjective() VehiclesDurationObjective {
	connect.Connect(con, &newVehiclesDurationObjective)
	return newVehiclesDurationObjective()
}

// VehiclesDurationObjective is an objective that uses the vehicle duration as an
// objective.
type VehiclesDurationObjective interface {
	ModelObjective
}
