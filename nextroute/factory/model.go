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
	options ...nextroute.Option,
) (nextroute.Model, error) {
	connect.Connect(con, &newModel)
	return newModel(input, modelOptions, options...)
}

var (
	con = connect.NewConnector("sdk", "NextRouteFactory")

	newModel func(
		schema.Input,
		Options,
		...nextroute.Option,
	) (nextroute.Model, error)
)

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
			Precedence         bool `json:"precedence" usage:"ignore the precedence (pickups & deliveries) constraint"`
			VehicleStartTime   bool `json:"vehicle_start_time" usage:"ignore the vehicle start time constraint"`
			VehicleEndTime     bool `json:"vehicle_end_time" usage:"ignore the vehicle end time constraint"`
			ArrivalTimeWindows bool `json:"arrival_time_windows" usage:"ignore the arrival time windows constraint"`
		} `json:"disable"`
		Enable struct {
			Cluster bool `json:"cluster" usage:"enable the cluster constraint"`
		} `json:"enable"`
	} `json:"constraints"`
	Objectives struct {
		EarlyArrivalPenalty      float64 `json:"early_arrival_penalty" usage:"factor to weigh the early arrival objective" default:"1.0"`
		LateArrivalPenalty       float64 `json:"late_arrival_penalty" usage:"factor to weigh the late arrival objective" default:"1.0"`
		VehicleActivationPenalty float64 `json:"vehicle_activation_penalty" usage:"factor to weigh the vehicle activation objective" default:"1.0"`
		TravelDuration           float64 `json:"travel_duration" usage:"factor to weigh the travel duration objective" default:"1.0"`
		UnplannedPenalty         float64 `json:"unplanned_penalty" usage:"factor to weigh the unplanned objective" default:"1.0"`
		Cluster                  float64 `json:"cluster" usage:"factor to weigh the cluster objective" default:"0.0"`
	} `json:"objectives"`
	Properties struct {
		Disable struct {
			Durations               bool `json:"durations" usage:"ignore the durations of stops"`
			StopDurationMultipliers bool `json:"stop_duration_multipliers" usage:"ignore the stop duration multipliers defined on vehicles"`
			DurationGroups          bool `json:"duration_groups" usage:"ignore the durations groups of stops"`
			InitialSolution         bool `json:"initial_solution" usage:"ignore the initial solution"`
		} `json:"disable"`
	} `json:"properties"`
}
