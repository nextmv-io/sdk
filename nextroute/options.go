package nextroute

import (
	"github.com/nextmv-io/sdk/nextroute/schema"
)

// An Option configures a model.
type Option func(Model, schema.Input) error

// ModelOptions configures how a [Model] is built.
type ModelOptions struct {
	Enable struct {
		Constraint struct {
			Cluster bool `json:"cluster" usage:"enable the cluster constraint"`
		} `json:"constraint"`
	} `json:"enable"`
	Ignore struct {
		Constraint struct {
			Attributes      bool `json:"attributes" usage:"ignore the compatibility attributes constraint"`
			Capacity        bool `json:"capacity" usage:"ignore the capacity constraint"`
			DistanceLimit   bool `json:"distance_limit" usage:"ignore the distance limit constraint"`
			MaximumDuration bool `json:"maximum_duration" usage:"ignore the maximum duration constraint"`
			Precedence      bool `json:"precedence" usage:"ignore the precedence (pickups & deliveries) constraint"`
			Shift           bool `json:"shift" usage:"ignore the shift constraint"`
			Windows         bool `json:"windows" usage:"ignore the windows constraint"`
		} `json:"constraint"`
		Objective struct {
			InitializationCost bool `json:"initialization_cost" usage:"ignore the initialization cost objective"`
			TravelDuration     bool `json:"travel_duration" usage:"ignore the travel duration objective"`
			UnassignedStops    bool `json:"unassigned_stops" usage:"ignore the unplanned objective"`
		} `json:"objective"`
		Characteristics struct {
			ServiceDurations bool `json:"service_durations" usage:"ignore the service durations of stops"`
		} `json:"characteristics"`
	} `json:"ignore"`
}

// RunnerOptions are used by the runner and configure the model and the solver.
type RunnerOptions struct {
	Model  ModelOptions         `json:"model"`
	Solver ParallelSolveOptions `json:"solver"`
}
