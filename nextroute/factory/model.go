// Package factory is a package containing factory functions for creating
// next-route models.
package factory

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute"
	"github.com/nextmv-io/sdk/nextroute/schema"
)

// NewModel builds a ready-to-go vehicle routing problem. The difference with
// [nextroute.NewModel] is that NewModel processes the input and options to add
// all features to the model, such as constraints and objectives. On the other
// hand, [nextroute.NewModel] creates an empty vehicle routing model which
// must be built from the ground up.
func NewModel(
	input schema.Input,
	modelOptions Options,
) (nextroute.Model, error) {
	connect.Connect(con, &newModel)
	return newModel(input, modelOptions)
}

// Options configure how the [NewModel] function builds [nextroute.Model].
type Options struct {
	Constraints struct {
		Disable struct {
			Attributes         bool `json:"attributes" usage:"ignore the compatibility attributes constraint"`
			Capacity           bool `json:"capacity" usage:"ignore the capacity constraint"`
			DistanceLimit      bool `json:"distance_limit" usage:"ignore the distance limit constraint"`
			Groups             bool `json:"groups" usage:"ignore the groups constraint"`
			MaximumDuration    bool `json:"maximum_duration" usage:"ignore the maximum duration constraint"`
			MaximumStops       bool `json:"maximum_stops" usage:"ignore the maximum stops constraint"`
			MaximumWaitStop    bool `json:"maximum_wait_stop" usage:"ignore the maximum stop wait constraint"`
			MaximumWaitVehicle bool `json:"maximum_wait_vehicle" usage:"ignore the maximum vehicle wait constraint"`
			MixingItems        bool `json:"mixing_items" usage:"ignore the do not mix items constraint"`
			Precedence         bool `json:"precedence" usage:"ignore the precedence (pickups & deliveries) constraint"`
			VehicleStartTime   bool `json:"vehicle_start_time" usage:"ignore the vehicle start time constraint"`
			VehicleEndTime     bool `json:"vehicle_end_time" usage:"ignore the vehicle end time constraint"`
			StartTimeWindows   bool `json:"start_time_windows" usage:"ignore the start time windows constraint"`
		} `json:"disable"`
		Enable struct {
			Cluster bool `json:"cluster" usage:"enable the cluster constraint"`
		} `json:"enable"`
	} `json:"constraints"`
	Objectives struct {
		MinStops                 float64 `json:"min_stops" usage:"factor to weigh the min stops objective" default:"1.0"`
		EarlyArrivalPenalty      float64 `json:"early_arrival_penalty" usage:"factor to weigh the early arrival objective" default:"1.0"`
		LateArrivalPenalty       float64 `json:"late_arrival_penalty" usage:"factor to weigh the late arrival objective" default:"1.0"`
		VehicleActivationPenalty float64 `json:"vehicle_activation_penalty" usage:"factor to weigh the vehicle activation objective" default:"1.0"`
		TravelDuration           float64 `json:"travel_duration" usage:"factor to weigh the travel duration objective" default:"0.0"`
		VehiclesDuration         float64 `json:"vehicles_duration" usage:"factor to weigh the vehicles duration objective" default:"1.0"`
		UnplannedPenalty         float64 `json:"unplanned_penalty" usage:"factor to weigh the unplanned objective" default:"1.0"`
		Cluster                  float64 `json:"cluster" usage:"factor to weigh the cluster objective" default:"0.0"`
		BalanceStopsPenalty      float64 `json:"balance_stops_penalty" usage:"factor to weigh the stops balancing objective" default:"0.0"`
		BalanceDurationPenalty   float64 `json:"balance_duration_penalty" usage:"factor to weigh the duration balancing objective" default:"0.0"`
	} `json:"objectives"`
	Properties struct {
		Disable struct {
			Durations               bool `json:"durations" usage:"ignore the durations of stops"`
			StopDurationMultipliers bool `json:"stop_duration_multipliers" usage:"ignore the stop duration multipliers defined on vehicles"`
			DurationGroups          bool `json:"duration_groups" usage:"ignore the durations groups of stops"`
			InitialSolution         bool `json:"initial_solution" usage:"ignore the initial solution"`
		} `json:"disable"`
	} `json:"properties"`
	Settings struct {
		BalanceStops struct {
			Mode string `json:"mode" usage:"mode to use for balancing stops (one of: max, maxmin)" default:"max"`
		} `json:"balance_stops"`
		BalanceDuration struct {
			Mode    string  `json:"mode" usage:"mode to use for balancing duration (one of: max, maxmin)" default:"max"`
			RoundTo float64 `json:"round_to" usage:"rounds duration to nearest multiple to avoid nervosity" default:"1800.0"`
		} `json:"balance_duration"`
	} `json:"settings"`
	Validate struct {
		Disable struct {
			StartTime bool `json:"start_time" usage:"disable the start time validation" default:"false"`
			Resources bool `json:"resources" usage:"disable the resources validation" default:"false"`
		} `json:"disable"`
		Enable struct {
			Matrix                   bool `json:"matrix" usage:"enable matrix validation" default:"false"`
			MatrixAsymmetryTolerance int  `json:"matrix_asymmetry_tolerance" usage:"percentage of acceptable matrix asymmetry, requires matrix validation enabled" default:"20"`
		} `json:"enable"`
	} `json:"validate"`
}
