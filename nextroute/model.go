package nextroute

import (
	"math/rand"
	"time"

	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/common"
)

// NewModel creates a new model. The model is used to define a routing problem.
func NewModel(
	options ...ModelOption,
) (Model, error) {
	connect.Connect(con, &newModel)
	return newModel(options...)
}

// Model defines routing problem.
type Model interface {
	SolutionObserved

	// AddConstraint adds a constraint to the model. The constraint is
	// checked at the specified violation.
	AddConstraint(constraint ModelConstraint) error

	// Constraints returns all constraints of the model.
	Constraints() ModelConstraints

	// ConstraintsCheckedAt returns all constraints of the model that
	// are checked at the specified time of having calculated the new
	// information for the changed solution.
	ConstraintsCheckedAt(violation CheckedAt) ModelConstraints

	// DistanceUnit returns the unit of distance used in the model. The
	// unit is used to convert distances to values and vice versa. This is
	// also used for reporting.
	DistanceUnit() common.DistanceUnit
	// DurationUnit returns the unit of duration used in the model. The
	// unit is used to convert durations to values and vice versa. This is
	// also used for reporting.
	DurationUnit() time.Duration
	// DurationToValue converts the specified duration to a value as it used
	// internally in the model.
	DurationToValue(duration time.Duration) float64

	// Epoch returns the epoch of the model. The epoch is used to convert
	// time.Time to float64 and vice versa. All float64 values are relative
	// to the epoch.
	Epoch() time.Time
	// Expressions returns all expressions of the model for which a solution
	// has to calculate values. The expressions are sorted by their index. The
	// constraints register their expressions with the model.
	Expressions() ModelExpressions

	// IsImmutable returns true if the model is immutable. The model is
	// immutable after a solution has been created using the model.
	IsImmutable() bool

	// NewPlanSingleStop creates a new plan single stop. A plan single stop
	// is a plan cluster of a single stop. A plan cluster is a collection of
	// stops which are always planned and unplanned as a single entity.
	NewPlanSingleStop(
		stop ModelStop,
	) ModelPlanSingleStop
	NewStop(
		location common.Location,
		options ...StopOption,
	) (ModelStop, error)

	// NewVehicleType creates a new vehicle type. The vehicle type is used
	// to create vehicles.
	NewVehicleType(
		travelDuration TravelDurationExpression,
		processDuration DurationExpression,
		options ...VehicleTypeOption,
	) (ModelVehicleType, error)
	// NumberOfStops returns the number of stops in the model.
	NumberOfStops() int

	// Objective returns the objective of the model.
	Objective() ModelObjectiveSum

	// PlanClusters returns all plan clusters of the model. A plan cluster
	// is a collection of stops which are always planned and unplanned as a
	// single entity.
	PlanClusters() ModelPlanClusters

	// Random returns a random number generator.
	Random() *rand.Rand

	// Stops returns all stops of the model.
	Stops() ModelStops
	// Stop returns the stop with the specified index.
	Stop(index int) ModelStop

	// TimeFormat returns the time format used for reporting.
	TimeFormat() string

	// VehicleTypes returns all vehicle types of the model.
	VehicleTypes() ModelVehicleTypes
	// VehicleType returns the vehicle type with the specified index.
	VehicleType(index int) ModelVehicleType
}

// ModelOption is a function that configures a model.
type ModelOption func(m Model) error

func DistanceUnit(
	distanceUnit common.DistanceUnit,
) ModelOption {
	connect.Connect(con, &distanceUnitModelOption)
	return distanceUnitModelOption(distanceUnit)
}

// DurationUnit is a model option to set the duration unit.
func DurationUnit(
	durationUnit time.Duration,
) ModelOption {
	connect.Connect(con, &durationUnitModelOption)
	return durationUnitModelOption(durationUnit)
}

// TimeFormat is a model option to set the time format.
func TimeFormat(
	timeFormat string,
) ModelOption {
	connect.Connect(con, &timeFormatModelOption)
	return timeFormatModelOption(timeFormat)
}

// Epoch is a model option to set the epoch.
func Epoch(
	epoch time.Time,
) ModelOption {
	connect.Connect(con, &epochModelOption)
	return epochModelOption(epoch)
}

// NewVehicleType is a model option to create a new vehicle type.
func NewVehicleType(
	travelDuration TravelDurationExpression,
	processDuration DurationExpression,
	options ...VehicleTypeOption,
) ModelOption {
	connect.Connect(con, &newVehicleTypeModelOption)
	return newVehicleTypeModelOption(
		travelDuration,
		processDuration,
		options...,
	)
}

// NewPlanSingleStops is a model option to create new plan single stops from
// locations.
func NewPlanSingleStops(
	locations common.Locations,
	options ...StopOption,
) ModelOption {
	connect.Connect(con, &newPlanSingleStopsModelOption)
	return newPlanSingleStopsModelOption(
		locations,
		options...,
	)
}

// MinimizeTravelDuration is a model option to minimize the travel duration.
func MinimizeTravelDuration(factor float64) ModelOption {
	connect.Connect(con, &minimizeTravelDurationModelOption)
	return minimizeTravelDurationModelOption(factor)
}

// MinimizeUnplannedStops is a model option to minimize the unplanned stops.
func MinimizeUnplannedStops(
	factor float64,
	costs StopExpression,
) ModelOption {
	connect.Connect(con, &minimizeUnplannedStopsModelOption)
	return minimizeUnplannedStopsModelOption(
		factor,
		costs,
	)
}

// MinimizeVehicleCost is a model option to minimize the vehicle cost.
func MinimizeVehicleCost(
	factor float64,
	costs VehicleTypeExpression,
) ModelOption {
	connect.Connect(con, &minimizeVehicleCost)
	return minimizeVehicleCost(
		factor,
		costs,
	)
}

// ConstrainStopsPerVehicleType is a model option to constrain the stops per
// vehicle type.
func ConstrainStopsPerVehicleType(
	maximumStops VehicleTypeExpression,
) ModelOption {
	connect.Connect(con, &constrainStopsPerVehicleTypeModelOption)
	return constrainStopsPerVehicleTypeModelOption(
		maximumStops,
	)
}

// ConstrainVehicleCompactness is a model option to constrain the vehicle
// compactness.
func ConstrainVehicleCompactness(
	limit StopExpression,
) ModelOption {
	connect.Connect(con, &constrainVehicleCompactnessModelOption)
	return constrainVehicleCompactnessModelOption(
		limit,
	)
}

// Random is a model option to set the random number generator.
func Random(random *rand.Rand) ModelOption {
	connect.Connect(con, &randomModelOption)
	return randomModelOption(random)
}

// Seed is a model option to set the random number generator seed.
func Seed(seed int64) ModelOption {
	connect.Connect(con, &seedModelOption)
	return seedModelOption(
		seed,
	)
}

// AddConstraint is a model option to add a constraint.
func AddConstraint(constraint ModelConstraint) ModelOption {
	connect.Connect(con, &addConstraintModelOption)
	return addConstraintModelOption(
		constraint,
	)
}
