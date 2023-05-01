package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/common"
)

// NewModelStatistics creates a new ModelStatistics instance.
func NewModelStatistics(model Model) ModelStatistics {
	connect.Connect(con, &newModelStatistics)
	return newModelStatistics(model)
}

// NewVehicleStatistics creates a new VehicleStatistics instance.
func NewVehicleStatistics(vehicle ModelVehicle) VehicleStatistics {
	connect.Connect(con, &newVehicleStatistics)
	return newVehicleStatistics(vehicle)
}

// ModelStatistics provides statistics for a model.
type ModelStatistics interface {
	// FirstLocations returns the number of unique locations that are first
	// locations of a vehicle.
	FirstLocations() int

	// LastLocations returns the number of unique locations that are last
	// locations of a vehicle.
	LastLocations() int
	// Locations returns the number of unique locations excluding first and last
	// locations of a vehicle.
	Locations() int

	// PlanUnits returns the number of plan units.
	PlanUnits() int

	// Report returns a report of the statistics.
	Report() string

	// Stops returns the number of stops.
	Stops() int

	// VehicleTypes returns the number of vehicle types.
	VehicleTypes() int
	// Vehicles returns the number of vehicles.
	Vehicles() int
}

// VehicleStatistics provides statistics for a vehicle.
type VehicleStatistics interface {
	// FirstToLastSeconds returns the travel time from the first location to the
	// last location of a vehicle.
	FirstToLastSeconds() float64
	// FromFirstSeconds returns the travel time in seconds from the first
	// location to all stops as statistics.
	FromFirstSeconds() common.Statistics

	// Report returns a report of the statistics.
	Report() string

	// ToLastSeconds returns the travel time in seconds from all stops to the
	// last location as statistics.
	ToLastSeconds() common.Statistics
}
