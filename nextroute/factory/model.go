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
			Attributes      bool `json:"attributes" usage:"ignore the compatibility attributes constraint"`
			Capacity        bool `json:"capacity" usage:"ignore the capacity constraint"`
			DistanceLimit   bool `json:"distance_limit" usage:"ignore the distance limit constraint"`
			MaximumDuration bool `json:"maximum_duration" usage:"ignore the maximum duration constraint"`
			Precedence      bool `json:"precedence" usage:"ignore the precedence (pickups & deliveries) constraint"`
			Shift           bool `json:"shift" usage:"ignore the shift constraint"`
			Windows         bool `json:"windows" usage:"ignore the windows constraint"`
		}
		Enable struct {
			Cluster bool `json:"cluster" usage:"enable the cluster constraint"`
		}
	}
	Objectives struct {
		InitializationCost float64 `json:"initialization_cost" usage:"factor to weigh the initialization cost objective" default:"1.0"`
		TravelDuration     float64 `json:"travel_duration" usage:"factor to weigh the travel duration objective" default:"1.0"`
		UnassignedStops    float64 `json:"unassigned_stops" usage:"factor to weigh the unplanned objective" default:"1.0"`
	}
	Properties struct {
		Disable struct {
			ServiceDurations bool `json:"service_durations" usage:"ignore the service durations of stops"`
		}
	}
}
