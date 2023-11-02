package factory

import (
	"context"

	"github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute"
	"github.com/nextmv-io/sdk/nextroute/common"
	"github.com/nextmv-io/sdk/nextroute/schema"
	runSchema "github.com/nextmv-io/sdk/run/schema"
)

var (
	con = connect.NewConnector("sdk", "NextRouteFactory")

	newModel func(
		schema.Input,
		Options,
	) (nextroute.Model, error)

	format func(
		context.Context,
		any,
		alns.Progressioner,
		...nextroute.Solution,
	) runSchema.Output

	newStartSolution func(
		context.Context,
		schema.Input,
		Options,
		ModelFactory,
		nextroute.ParallelSolveOptions,
		ClusterSolutionOptions,
	) (nextroute.Solution, error)

	newGreedySolution func(
		context.Context,
		schema.Input,
		Options,
		GreedySolutionOptions,
		ModelFactory,
	) (nextroute.Solution, error)

	newStopCluster func(
		[]schema.Stop,
	) StopCluster

	newPlanUnitStopClusterGenerator func() StopClusterGenerator

	newSortStopClustersRandom func() StopClusterSorter

	newSortStopClustersOnDistanceFromCentroid func() StopClusterSorter

	newStopClusterFilterArea func(
		common.Distance,
		common.Distance,
		common.Distance,
	) StopClusterFilter

	newAndStopClusterFilter func(
		StopClusterFilter,
		...StopClusterFilter,
	) StopClusterFilter

	newOrStopClusterFilter func(
		StopClusterFilter,
		...StopClusterFilter,
	) StopClusterFilter

	newClusterSolution func(
		context.Context,
		schema.Input,
		Options,
		StopClusterGenerator,
		StopClusterSorter,
		StopClusterSorter,
		StopClusterFilter,
		ClusterSolutionOptions,
		ModelFactory,
	) (nextroute.Solution, error)

	newDefaultModelFactory func() ModelFactory
)
