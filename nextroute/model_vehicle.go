package nextroute

import "time"

// ModelVehicle is a vehicle in the model. A vehicle is a sequence of stops.
type ModelVehicle interface {
	ModelData

	// AddStop adds a stop to the vehicle. The stop is added to the end of the
	// vehicle, before the last stop. If fixed is true stop will be fixed and
	// can not be unplanned.
	AddStop(stop ModelStop, fixed bool) error

	// First returns the first stop of the vehicle.
	First() ModelStop

	// ID returns the identifier of the vehicle.
	ID() string
	// Index returns the index of the vehicle.
	Index() int

	// Last returns the last stop of the vehicle.
	Last() ModelStop

	// Model returns the model of the vehicle.
	Model() Model

	// SetID sets the identifier of the vehicle. This identifier is not used by
	// nextroute, and therefore it does not have to be unique for nextroute
	// internally. However, to make this ID useful for debugging and reporting it
	// should be made unique.
	SetID(string)
	// Start returns the start time of the vehicle.
	Start() time.Time
	// Stops returns the stops of the vehicle that are provided as a start
	// assignment. The first and last stop of the vehicle are not included in
	// the returned slice.
	Stops() ModelStops

	// VehicleType returns the vehicle type of the vehicle.
	VehicleType() ModelVehicleType
}

// ModelVehicles is a slice of ModelVehicle.
type ModelVehicles []ModelVehicle
