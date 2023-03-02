package nextroute

import "time"

// ModelVehicle is a vehicle in the model. A vehicle is a sequence of stops.
type ModelVehicle interface {
	// VehicleType returns the vehicle type of the vehicle.
	VehicleType() ModelVehicleType

	// First returns the first stop of the vehicle.
	First() ModelStop

	// Index returns the index of the vehicle.
	Index() int

	// Last returns the last stop of the vehicle.
	Last() ModelStop

	// Name returns the name of the vehicle.
	Name() string

	// SetName sets the name of the vehicle.
	SetName(string)

	// Start returns the start time of the vehicle.
	Start() time.Time
}

// ModelVehicles is a slice of ModelVehicle.
type ModelVehicles []ModelVehicle
