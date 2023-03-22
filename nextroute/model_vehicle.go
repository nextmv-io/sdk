package nextroute

import "time"

// ModelVehicle is a vehicle in the model. A vehicle is a sequence of stops.
type ModelVehicle interface {
	ModelData

	// First returns the first stop of the vehicle.
	First() ModelStop

	// Index returns the index of the vehicle.
	Index() int

	// Last returns the last stop of the vehicle.
	Last() ModelStop

	// Model returns the model of the vehicle.
	Model() Model

	// Name returns the name of the vehicle.
	Name() string

	// SetName sets the name of the vehicle.
	SetName(string)
	// Start returns the start time of the vehicle.
	Start() time.Time

	// VehicleType returns the vehicle type of the vehicle.
	VehicleType() ModelVehicleType
}

// ModelVehicles is a slice of ModelVehicle.
type ModelVehicles []ModelVehicle
