// Package nextroute is a package
package nextroute

import (
	"context"
	"math/rand"
	"time"

	"github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/measure"
	"github.com/nextmv-io/sdk/nextroute/common"
	"github.com/nextmv-io/sdk/run/schema"
)

var (
	con = connect.NewConnector("sdk", "NextRoute")

	newRandomSolution func(
		context.Context,
		Model,
	) (Solution, error)
	newSweepSolution func(
		context.Context,
		Model,
	) (Solution, error)
	newClusterSolution func(
		context.Context,
		Model,
		int,
	) (Solution, error)
	newConstantDurationExpression func(
		string,
		time.Duration,
	) DurationExpression
	newConstantExpression func(
		string,
		float64,
	) ConstantExpression
	newDistanceExpression func(
		string,
		ModelExpression,
		common.DistanceUnit,
	) DistanceExpression
	newDurationExpression func(
		string,
		ModelExpression,
		common.DurationUnit,
	) DurationExpression
	newScaledDurationExpression func(
		DurationExpression,
		float64,
	) DurationExpression
	newNotExecutableMove  func() SolutionMove
	newFromStopExpression func(
		string,
		float64,
	) FromStopExpression
	newFromToExpression func(
		string,
		float64,
	) FromToExpression
	newHaversineExpression  func() DistanceExpression
	newAttributesConstraint func() (AttributesConstraint, error)
	newCluster              func() (Cluster, error)
	newLatestEnd            func(
		StopTimeExpression,
	) (LatestEnd, error)
	newLatestStart func(
		StopTimeExpression,
	) (LatestStart, error)
	newLatestArrival func(
		StopTimeExpression,
	) (LatestArrival, error)
	newMaximumConstraint func(
		ModelExpression,
		VehicleTypeExpression,
	) (MaximumConstraint, error)
	newComposedPerVehicleTypeExpression func(
		ModelExpression,
	) ComposedPerVehicleTypeExpression
	newTimeExpression func(
		ModelExpression,
		time.Time,
	) TimeExpression
	newStopTimeExpression func(
		string,
		time.Time,
	) StopTimeExpression
	newStopDurationExpression func(
		string,
		time.Duration,
	) StopDurationExpression
	newVehicleTypeDurationExpression func(
		string,
		time.Duration,
	) VehicleTypeDurationExpression

	newMaximumDurationConstraint func(
		VehicleTypeDurationExpression,
	) (MaximumDurationConstraint, error)

	newMaximumTravelDurationConstraint func(
		VehicleTypeDurationExpression,
	) (MaximumTravelDurationConstraint, error)

	newMaximumWaitStopConstraint    func() (MaximumWaitStopConstraint, error)
	newMaximumWaitVehicleConstraint func() (MaximumWaitVehicleConstraint, error)

	newMaximumStopsConstraint func(
		VehicleTypeExpression,
	) (MaximumStopsConstraint, error)

	newNoMixConstraint func(
		map[ModelStop]MixItem,
	) (NoMixConstraint, error)

	newMeasureByIndexExpression func(
		measure.ByIndex,
	) ModelExpression
	newMeasureByPointExpression func(
		measure.ByPoint,
	) ModelExpression
	newModel                 func() (Model, error)
	newModelExpressionIndex  func() int
	newSolutionStopGenerator func(
		SolutionMoveStops,
		bool,
		bool,
	) SolutionStopGenerator
	noPositionsHint       func() StopPositionsHint
	newOperatorExpression func(
		ModelExpression,
		ModelExpression,
		BinaryFunction,
	) BinaryExpression
	skipVehiclePositionsHint func() StopPositionsHint
	newStopExpression        func(
		string,
		float64,
	) StopExpression
	newSolution func(
		Model,
	) (Solution, error)
	newSumExpression func(
		ModelExpressions,
	) SumExpression
	newTermExpression func(
		float64,
		ModelExpression,
	) TermExpression
	newTravelDurationExpression func(
		DistanceExpression,
		common.Speed,
	) TravelDurationExpression
	newTravelDurationObjective func() TravelDurationObjective
	newUnPlannedObjective      func(
		StopExpression,
	) UnPlannedObjective
	newExpressionObjective func(ModelExpression) ExpressionObjective
	newSolverFactory       func() SolverFactory
	newVehiclesObjective   func(
		VehicleTypeExpression,
	) VehiclesObjective
	newVehiclesDurationObjective  func() VehiclesDurationObjective
	newVehicleTypeValueExpression func(
		string,
		float64,
	) VehicleTypeValueExpression
	newVehicleTypeDistanceExpression func(
		string,
		common.Distance,
	) VehicleTypeDistanceExpression
	newVehicleTypeFromToExpression func(
		string,
		float64,
	) VehicleFromToExpression

	newSolver func(
		Model,
		SolverOptions,
	) (Solver, error)

	newParallelSolver func(
		Model,
	) (ParallelSolver, error)

	newSolutionPlanUnitCollection func(
		*rand.Rand,
		SolutionPlanUnits,
	) SolutionPlanUnitCollection

	newPerformanceObserver func(
		Model,
	) PerformanceObserver

	newModelStatistics func(
		Model,
	) ModelStatistics

	newVehicleStatistics func(
		ModelVehicle,
	) VehicleStatistics
	newEarlinessObjective func(
		StopTimeExpression,
		StopExpression,
		TemporalReference,
	) (EarlinessObjective, error)
	newDirectedAcyclicGraph func() DirectedAcyclicGraph

	newTimeDependentDurationExpression func(
		Model,
		DurationExpression,
	) (TimeDependentDurationExpression, error)

	newTimeIndependentDurationExpression func(
		DurationExpression,
	) TimeDependentDurationExpression

	format func(
		context.Context,
		any,
		alns.Progressioner,
		func(Solution) any,
		...Solution,
	) schema.Output

	newModelStopsDistanceQueries func(
		ModelStops,
	) (ModelStopsDistanceQueries, error)
)
