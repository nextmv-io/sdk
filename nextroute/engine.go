package nextroute

import (
	"math/rand"
	"time"

	"github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/measure"
)

var (
	con = connect.NewConnector("sdk", "NextRoute")

	newRandomSolution             func(Model) (Solution, error)
	newSweepSolution              func(Model) (Solution, error)
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
		DistanceUnit,
	) DistanceExpression
	newDurationExpression func(
		ModelExpression,
		time.Duration,
	) DurationExpression
	newNotExecutableMove  func() Move
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
	newInwardnessConstraint func() (InwardnessConstraint, error)
	newLatestEnd            func(
		StopTimeExpression,
	) (LatestEnd, error)
	newLatestStart func(
		StopTimeExpression,
	) (LatestStart, error)
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
	newMaximumStopsConstraint func(
		VehicleTypeExpression,
	) (MaximumStopsConstraint, error)
	newMeasureByIndexExpression func(
		measure.ByIndex,
	) ModelExpression
	newMeasureByPointExpression func(
		measure.ByPoint,
	) ModelExpression
	newModel                func() (Model, error)
	newModelConstraintIndex func() int
	newModelExpressionIndex func() int
	newModelObjectiveIndex  func() int
	noPositionsHint         func() StopPositionsHint
	newOperatorExpression   func(
		ModelExpression,
		ModelExpression,
		BinaryFunction,
	) BinaryExpression
	skipVehiclePositionsHint func(
		bool,
	) StopPositionsHint
	newStopExpression func(
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
		Speed,
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
	newVehicleTypeExpression func(
		string,
		float64,
	) VehicleTypeExpression
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
	newBasicFormatter                func() Formatter
	newVerboseFormatter              func([]alns.ProgressionEntry) Formatter
	newSolutionPlanClusterCollection func(
		*rand.Rand,
		SolutionPlanClusters,
	) SolutionPlanClusterCollection
	newPerformanceObserver func(
		Model,
	) PerformanceObserver
	newModelStatistics func(
		Model,
	) ModelStatistics
	newVehicleStatistics func(
		ModelVehicle,
	) VehicleStatistics
	newLocation        func(float64, float64) (Location, error)
	newInvalidLocation func() Location
	unique             func(Locations) Locations
	centroid           func(Locations) (Location, error)
	newSpeed           func(float64, SpeedUnit) Speed
	newSpeedUnit       func(DistanceUnit, time.Duration) SpeedUnit
	newDistance        func(float64, DistanceUnit) Distance
	haversine          func(Location, Location) (Distance, error)
	durationValue      func(Distance, Speed, time.Duration) float64
)
