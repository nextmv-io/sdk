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
	modelOptions ModelOptions,
	options ...nextroute.Option,
) (nextroute.Model, error) {
	connect.Connect(con, &newModel)
	return newModel(input, modelOptions, options...)
}

// ModelOptions represents options for a model.
type ModelOptions struct {
	EnableClusterConstraint           bool `json:"enable_cluster_constraint" usage:"Enable the cluster constraint creating compact clusters of stops"`
	EnableClusterObjective            bool `json:"enable_cluster_objective" usage:"Enable the cluster objective creating compact clusters of stops"`
	IgnoreAttributesConstraint        bool `json:"ignore_attributes_constraint" usage:"Ignore the compatibility attributes constraint"`
	IgnoreCapacityConstraint          bool `json:"ignore_capacity_constraint" usage:"Ignore the capacity constraint"`
	IgnoreDistanceLimitConstraint     bool `json:"ignore_distance_limit_constraint" usage:"Ignore the distance limit constraint"`
	IgnoreInitializationCostObjective bool `json:"ignore_initialization_cost_objective" usage:"Ignore the initialization cost objective"`
	IgnoreMaximumDurationConstraint   bool `json:"ignore_maximum_duration_constraint" usage:"Ignore the maximum duration constraint"`
	IgnorePrecedenceConstraint        bool `json:"ignore_precedence_constraint" usage:"Ignore the precedence (pickups & deliveries) constraint"`
	IgnoreServiceDurations            bool `json:"ignore_service_durations" usage:"Ignore the service durations of stops"`
	IgnoreShiftConstraint             bool `json:"ignore_shift_constraint" usage:"Ignore the shift constraint"`
	IgnoreTravelDurationObjective     bool `json:"ignore_travel_duration_objective" usage:"Ignore the travel duration objective"`
	IgnoreUnassignedStopsObjective    bool `json:"ignore_unassigned_stops_objective" usage:"Ignore the unplanned objective"`
	IgnoreWindowsConstraint           bool `json:"ignore_windows_constraint" usage:"Ignore the windows constraint"`
	IgnoreEarlinessObjective          bool `json:"ignore_earliness_objective" usage:"Ignore the earliness objective"`
}

var (
	con = connect.NewConnector("sdk", "NextRouteFactory")

	newModel func(
		schema.Input,
		ModelOptions,
		...nextroute.Option,
	) (nextroute.Model, error)
)
