package nextroute

import (
	"math/rand"
	"time"

	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/common"
)

// NewModel creates a new model. The model is used to define a routing problem.
func NewModel() (Model, error) {
	connect.Connect(con, &newModel)
	return newModel()
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
	NewPlanSingleStop(stop ModelStop) ModelPlanSingleStop

	// NewStop creates a new stop. The stop is used to create plan clusters.
	NewStop(location common.Location) (ModelStop, error)

	// NewVehicle creates a new vehicle. The vehicle is used to create
	// solutions.
	NewVehicle(
		vehicleType ModelVehicleType,
		start time.Time,
		first ModelStop,
		last ModelStop,
	) (ModelVehicle, error)
	// NewVehicleType creates a new vehicle type. The vehicle type is used
	// to create vehicles.
	NewVehicleType(
		travelDuration TravelDurationExpression,
		processDuration DurationExpression,
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

	// SetDistanceUnit sets the distance unit of the model.
	SetDistanceUnit(distanceUnit common.DistanceUnit)

	// SetDurationUnit sets the duration unit of the model.
	SetDurationUnit(durationUnit time.Duration)

	// SetEpoch sets the epoch of the model. The epoch is used to convert
	// time.Time to float64 and vice versa. All float64 values are relative
	// to the epoch.
	SetEpoch(epoch time.Time)

	// SetRandom sets the random number generator of the model.
	SetRandom(random *rand.Rand)

	// SetSeed sets the seed of the random number generator of the model.
	SetSeed(seed int64)

	// Stops returns all stops of the model.
	Stops() ModelStops

	// Stop returns the stop with the specified index.
	Stop(index int) ModelStop

	// TimeFormat returns the time format used for reporting.
	TimeFormat() string

	// Vehicles returns all vehicles of the model.
	Vehicles() ModelVehicles
	// VehicleTypes returns all vehicle types of the model.
	VehicleTypes() ModelVehicleTypes

	// Vehicle returns the vehicle with the specified index.
	Vehicle(index int) ModelVehicle
}
