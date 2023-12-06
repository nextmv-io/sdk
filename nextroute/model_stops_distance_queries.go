package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/common"
)

// ModelStopsDistanceQueries is an interface to query distances between stops.
type ModelStopsDistanceQueries interface {
	// ModelStops returns the original set of stops that the distance queries
	// were created from.
	ModelStops() ModelStops
	// NearestStops returns the n nearest stops to the given stop, the stop must
	// be present in the original set of stops.
	NearestStops(stop ModelStop, n int) (ModelStops, error)
	// WithinDistanceStops returns the stops within the given distance of the
	// given stop, the stop must be present in the original set of stops.
	WithinDistanceStops(
		stop ModelStop,
		distance common.Distance,
	) (ModelStops, error)
}

// NewModelStopsDistanceQueries returns a new ModelStopsDistanceQueries.
// The ModelStopsDistanceQueries can be used to query distances between stops.
// The stops must be a set of stops with valid locations.
func NewModelStopsDistanceQueries(stops ModelStops) (ModelStopsDistanceQueries, error) {
	connect.Connect(con, &newModelStopsDistanceQueries)
	return newModelStopsDistanceQueries(stops)
}
