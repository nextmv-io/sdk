package nextroute

// SolutionVehicleFilter is a filter to be applied
// to filter out vehicles that should not be evaluated.
type SolutionVehicleFilter interface {
	// FilterVehicle returns true if the vehicle should be filtered out.
	FilterVehicle(vehicle SolutionVehicle) bool
}

// SolutionVehicleFilters is a slice of SolutionVehicleFilter.
type SolutionVehicleFilters []SolutionVehicleFilter
