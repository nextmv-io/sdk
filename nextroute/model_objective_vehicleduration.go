package nextroute

import "github.com/nextmv-io/sdk/connect"

// NewVehicleDurationObjective returns a new VehicleDurationObjective that
// uses the vehicle duration as an objective.
func NewVehiclesDurationObjective() VehiclesDurationObjective {
	connect.Connect(con, &newVehiclesDurationObjective)
	return newVehiclesDurationObjective()
}

// VehicleDurationObjective is an objective that uses the vehicle duration as an
// objective.
type VehiclesDurationObjective interface {
	ModelObjective
}
