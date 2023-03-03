// Package nextroute is a package
package nextroute

import (
	"time"

	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/measure"
	"github.com/nextmv-io/sdk/nextroute/common"
)

var (
	con = connect.NewConnector("sdk", "NextRoute")

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
		ModelExpression,
		time.Duration,
	) DurationExpression
	newEmptyMove          func() Move
	newFromStopExpression func(
		string,
		float64,
	) FromStopExpression
	newFromToExpression func(
		string,
		float64,
	) FromToExpression
	newHaversineExpression func(
		bool,
	) DistanceExpression
	newAttributesConstraint func() (AttributesConstraint, error)
	newInwardnessConstraint func() (InwardnessConstraint, error)
	newLatestEnd            func(
		StopExpression,
	) (LatestEnd, error)
	newLatestStart func(
		StopExpression,
	) (LatestStart, error)
	newMaximumConstraint func(
		StopExpression,
		VehicleTypeExpression,
	) (MaximumConstraint, error)
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
	newNoStopPositionsHint  func() StopPositionsHint
	newOperatorExpression   func(
		ModelExpression,
		ModelExpression,
		BinaryFunction,
	) BinaryExpression
	newSkipVehiclePositionsHint func(
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
		common.Speed,
	) TravelDurationExpression
	newTravelDurationObjective func() TravelDurationObjective
	newUnPlannedObjective      func(
		StopExpression,
	) UnPlannedObjective
	newSolverFactory     func() SolverFactory
	newVehiclesObjective func(
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
	selectRandom func(
		SolutionPlanClusters,
		int,
	) SolutionPlanClusters

	newSolver func(
		Solution,
		SolverOptions,
	) (Solver, error)
)
