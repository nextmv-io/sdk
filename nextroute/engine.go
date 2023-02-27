// Package nextroute is a package
package nextroute

import (
	"math/rand"
	"time"

	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/common"
)

var (
	con = connect.NewConnector("sdk", "NextRoute")

	addConstraintModelOption func(
		ModelConstraint,
	) ModelOption
	constrainStopsPerVehicleTypeModelOption func(
		VehicleTypeExpression,
	) ModelOption
	constrainVehicleCompactnessModelOption func(
		StopExpression,
	) ModelOption
	distanceUnitModelOption func(
		common.DistanceUnit,
	) ModelOption
	durationUnitModelOption func(
		time.Duration,
	) ModelOption
	epochModelOption func(
		time.Time,
	) ModelOption
	minimizeTravelDurationModelOption func(
		float64,
	) ModelOption
	minimizeUnplannedStopsModelOption func(
		float64,
		StopExpression,
	) ModelOption
	minimizeVehicleCost func(
		float64,
		VehicleTypeExpression,
	) ModelOption
	newPlanSingleStopsModelOption func(
		common.Locations,
		...StopOption,
	) ModelOption
	newVehicleTypeModelOption func(
		TravelDurationExpression,
		DurationExpression,
		...VehicleTypeOption,
	) ModelOption
	randomModelOption func(
		*rand.Rand,
	) ModelOption
	seedModelOption func(
		int64,
	) ModelOption
	timeFormatModelOption func(
		string,
	) ModelOption
	earliestStartStopOption func(
		time.Time,
	) StopOption
	nameStopOption func(
		string,
	) StopOption
	newCompactnessConstraint func(
		StopExpression,
	) (CompactnessConstraint, error)
	newConstantDurationExpression func(
		string,
		time.Duration,
	) DurationExpression
	newConstantExpression func(
		string,
		float64,
	) ConstantExpression
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
	newInwardnessConstraint func() (InwardnessConstraint, error)
	newLatestEndConstraint  func(
		StopExpression,
	) (LatestEndConstraint, error)
	newMaximumConstraint func(
		StopExpression,
		VehicleTypeExpression,
	) (MaximumConstraint, error)
	newMaximumStopsConstraint func(
		VehicleTypeExpression,
	) (MaximumStopsConstraint, error)
	newModel func(
		...ModelOption,
	) (Model, error)
	newNoStopPositionsHint func() StopPositionsHint
	newOperatorExpression  func(
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
	newTravelDurationObjective func(
		float64,
	) TravelDurationObjective
	newUnPlannedObjective func(
		float64,
		StopExpression,
	) UnPlannedObjective
	newVehiclesObjective func(
		float64,
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
	selectRandomSolutionPlanClusters func(
		SolutionPlanClusters,
		int,
	) SolutionPlanClusters
	stopDataStopOption func(
		any,
	) StopOption
	vehicleTypeDataVehicleTypeOption func(
		any,
	) VehicleTypeOption
)
