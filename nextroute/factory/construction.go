package factory

import (
	"context"

	"github.com/nextmv-io/sdk/connect"
	sdkNextRoute "github.com/nextmv-io/sdk/nextroute"
	"github.com/nextmv-io/sdk/nextroute/common"
	"github.com/nextmv-io/sdk/nextroute/schema"
)

// ClusterSolutionOptions configure how the [NewGreedySolution] function builds [sdkNextRoute.Solution].
type ClusterSolutionOptions struct {
	Depth int     `json:"depth" usage:"maximum failed tries to add a cluster to a vehicle" default:"10" minimum:"0"`
	Speed float64 `json:"speed" usage:"speed of the vehicle in meters per second" default:"10" minimum:"0"`
}

// FilterAreaOptions configure how the [NewGreedySolution] function builds [sdkNextRoute.Solution]. It limits the area
// one vehicle can cover during construction. This limit is only applied during the construction of the solution.
type FilterAreaOptions struct {
	MaximumSide float64 `json:"maximum_side" usage:"maximum side of the square area in meters" default:"100000" minimum:"0"`
}

// GreedySolutionOptions configure how the [NewGreedySolution] function builds [sdkNextRoute.Solution].
type GreedySolutionOptions struct {
	ClusterSolutionOptions ClusterSolutionOptions `json:"cluster_solution_options" usage:"options for the cluster solution"`
	FilterAreaOptions      FilterAreaOptions      `json:"filter_area_options" usage:"options for the filter area"`
}

// NewStartSolution returns a start solution. It uses input, factoryOptions and
// modelFactory to create a model to create a start solution. The start solution
// is created using the given solveOptions and clusterSolutionOptions. The
// solveOptions is used to limit the duration and the number of parallel runs at
// the same time. The clusterSolutionOptions is used to create the clusters to
// create the start solution, see [NewClusterSolution].
func NewStartSolution(
	ctx context.Context,
	input schema.Input,
	factoryOptions Options,
	modelFactory ModelFactory,
	solveOptions sdkNextRoute.ParallelSolveOptions,
	clusterSolutionOptions ClusterSolutionOptions,
) (sdkNextRoute.Solution, error) {
	connect.Connect(con, &newStartSolution)
	return newStartSolution(
		ctx,
		input,
		factoryOptions,
		modelFactory,
		solveOptions,
		clusterSolutionOptions,
	)
}

// NewGreedySolution returns a greedy solution for the given input.
func NewGreedySolution(
	ctx context.Context,
	input schema.Input,
	options Options,
	greedySolutionOptions GreedySolutionOptions,
	modelFactory ModelFactory,
) (sdkNextRoute.Solution, error) {
	connect.Connect(con, &newGreedySolution)
	return newGreedySolution(
		ctx,
		input,
		options,
		greedySolutionOptions,
		modelFactory,
	)
}

// StopCluster represents a group of stops that can be added to a vehicle.
type StopCluster interface {
	// Stops returns the stops in the stop cluster.
	Stops() []schema.Stop
	// Centroid returns the centroid of the stop cluster.
	Centroid() schema.Location
}

// NewStopCluster returns a new stop cluster for the given stops.
func NewStopCluster(
	stops []schema.Stop) StopCluster {
	connect.Connect(con, &newStopCluster)
	return newStopCluster(stops)
}

// NewPlanUnitStopClusterGenerator returns a list of stop clusters based
// upon unplanned plan units.
func NewPlanUnitStopClusterGenerator() StopClusterGenerator {
	connect.Connect(con, &newPlanUnitStopClusterGenerator)
	return newPlanUnitStopClusterGenerator()
}

// NewSortStopClustersRandom returns StopClusterSorter which sorts the stop
// clusters randomly. Can be used to randomize the order of the stop clusters
// assigned to the vehicles as the first cluster.
func NewSortStopClustersRandom() StopClusterSorter {
	connect.Connect(con, &newSortStopClustersRandom)
	return newSortStopClustersRandom()
}

// NewSortStopClustersOnDistanceFromCentroid sorts the stop clusters based upon
// the distance from the centroid of the stop cluster to the centroid of all
// stops. Can be used to select the order of the stop clusters assigned to the
// vehicles as the first cluster.
func NewSortStopClustersOnDistanceFromCentroid() StopClusterSorter {
	connect.Connect(con, &newSortStopClustersOnDistanceFromCentroid)
	return newSortStopClustersOnDistanceFromCentroid()
}

// NewStopClusterFilterArea returns a StopClusterFilter that filters out stop
// clusters that result in an area larger than the dimensions specified.
// The area is approximated using haversine distance.
func NewStopClusterFilterArea(
	maximumWidth common.Distance,
	maximumHeight common.Distance,
	maximumRadius common.Distance,
) StopClusterFilter {
	connect.Connect(con, &newStopClusterFilterArea)
	return newStopClusterFilterArea(maximumWidth, maximumHeight, maximumRadius)
}

// NewAndStopClusterFilter returns a StopClusterFilter that filters out stop
// clusters that are filtered out by all the given filters.
func NewAndStopClusterFilter(
	filter StopClusterFilter,
	filters ...StopClusterFilter,
) StopClusterFilter {
	connect.Connect(con, &newAndStopClusterFilter)
	return newAndStopClusterFilter(filter, filters...)
}

// NewOrStopClusterFilter returns a StopClusterFilter that filters out stop
// clusters that are filtered out by any of the given filters.
func NewOrStopClusterFilter(
	filter StopClusterFilter,
	filters ...StopClusterFilter,
) StopClusterFilter {
	connect.Connect(con, &newOrStopClusterFilter)
	return newOrStopClusterFilter(filter, filters...)
}

// StopClusterGenerator returns a list of stop clusters for the given input.
type StopClusterGenerator interface {
	// Generate returns a list of stop clusters for the given input.
	// A cluster is a group of stops that can be added to a vehicle. If a stop
	// is added to a cluster all the stops belonging to the same plan units
	// must be added to the same cluster.
	Generate(
		input schema.Input,
		options Options,
		factory ModelFactory,
	) ([]StopCluster, error)
}

// StopClusterSorter returns a sorted list of stop clusters for the given input.
type StopClusterSorter interface {
	// Sort returns a sorted list of stop clusters for the given input.
	Sort(
		input schema.Input,
		clusters []StopCluster,
		factory ModelFactory,
	) ([]StopCluster, error)
}

// StopClusterFilter returns true if the given stop cluster should be filtered
// out.
type StopClusterFilter interface {
	// Filter returns true if the given stop cluster should be filtered out.
	Filter(
		input schema.Input,
		cluster StopCluster,
		factory ModelFactory,
	) (bool, error)
}

// ModelFactory returns a new model for the given input and options.
type ModelFactory interface {
	// NewModel returns a new model for the given input and options.
	NewModel(schema.Input, Options) (sdkNextRoute.Model, error)
}

// NewDefaultModelFactory returns a default model factory.
func NewDefaultModelFactory() ModelFactory {
	connect.Connect(con, &newDefaultModelFactory)
	return newDefaultModelFactory()
}

// NewClusterSolution returns a solution for the given input using the given
// options. The solution is constructed by first creating a solution for each
// vehicle and then adding stop groups to the vehicles in a greedy fashion.
//
//   - Raises an error if the input has initial stops on any of the vehicles.
//   - Uses haversine distance independent of the input's distance/duration
//     matrix. Uses the correct distance matrix in the solution returned.
//   - Uses the speed of the vehicle if defined, otherwise the speed defined in
//     the options.
//   - Ignores stop duration groups in construction but not in the solution
//     returned.
//
// # The initial solution is created as following:
//
//	Creates the clusters using the stopClusterGenerator
//
//	In random order of the vehicles in the input:
//
//	 - Add a first cluster to the empty vehicle defined by the
//	   initialStopClusterSorter
//	 - If the vehicle is not solved, the cluster is removed and the next cluster
//	   will be added
//	 - If no clusters can be added, the vehicle will not be used
//	 - If a cluster has been added we continue adding clusters to the vehicle in
//	   the order defined by additionalStopClusterSorter until no more clusters
//	   can be added
//
//	We repeat until no more vehicles or no more clusters to add to the solution.
func NewClusterSolution(
	ctx context.Context,
	input schema.Input,
	options Options,
	stopClusterGenerator StopClusterGenerator,
	initialStopClusterSorter StopClusterSorter,
	additionalStopClusterSorter StopClusterSorter,
	stopClusterFilter StopClusterFilter,
	stopClusterOptions ClusterSolutionOptions,
	modelFactory ModelFactory,
) (sdkNextRoute.Solution, error) {
	connect.Connect(con, &newClusterSolution)
	return newClusterSolution(
		ctx,
		input,
		options,
		stopClusterGenerator,
		initialStopClusterSorter,
		additionalStopClusterSorter,
		stopClusterFilter,
		stopClusterOptions,
		modelFactory,
	)
}
