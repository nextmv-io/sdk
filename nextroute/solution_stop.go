package nextroute

import (
	"time"
)

// A SolutionStop is a stop that is planned to be visited by a vehicle. It is
// part of a SolutionPlanCluster and is based on a ModelStop.
type SolutionStop interface {
	// Arrival returns the arrival time of the stop. If the stop is unplanned,
	// the arrival time has no semantic meaning.
	Arrival() time.Time
	// ArrivalValue returns the arrival time of the stop as a float64. If the
	// stop is unplanned, the arrival time has no semantic meaning.
	ArrivalValue() float64

	// ConstraintData returns the value of the constraint for the stop. The
	// constraint value of a stop is set by the ConstraintDataUpdater.Update
	// method of the constraint. If the constraint is not set on the stop,
	// nil is returned. If the stop is unplanned, the constraint value has no
	// semantic meaning.
	ConstraintData(constraint ModelConstraint) any
	// CumulativeTravelDurationValue returns the cumulative travel duration of
	// the stop as a float64. The cumulative travel duration is the sum of the
	// travel durations of all stops that are visited before the stop. If the
	// stop is unplanned, the cumulative travel duration has no semantic
	// meaning. The returned value is the number of Model.DurationUnit units.
	CumulativeTravelDurationValue() float64
	// CumulativeTravelDuration returns the cumulative value of the expression
	// for the stop as a time.Duration. The cumulative travel duration is the
	// sum of the travel durations of all stops that are visited before the
	// stop and the stop itself. If the stop is unplanned, the cumulative
	// travel duration has no semantic meaning.
	CumulativeTravelDuration() time.Duration
	// CumulativeValue returns the cumulative value of the expression for the
	// stop as a float64. The cumulative value is the sum of the values of the
	// expression for all stops that are visited before the stop and the stop
	// itself. If the stop is unplanned, the cumulative value has no semantic
	// meaning.
	CumulativeValue(expression ModelExpression) float64

	// End returns the end time of the stop. If the stop is unplanned, the end
	// time has no semantic meaning.
	End() time.Time
	// EndValue returns the end time of the stop as a float64. If the stop is
	// unplanned, the end time has no semantic meaning. The returned value is
	// the number of Model.DurationUnit units since Model.Epoch.
	EndValue() float64

	// Index returns the index of the stop in the Solution.
	Index() int
	// IsFixed returns true if the stop is fixed. A fixed stop is a stop that
	// that can not transition form being planned to unplanned or vice versa.
	IsFixed() bool
	// IsFirst returns true if the stop is the first stop of a vehicle.
	IsFirst() bool
	// IsLast returns true if the stop is the last stop of a vehicle.
	IsLast() bool
	// IsPlanned returns true if the stop is planned. A planned stop is a stop
	// that is visited by a vehicle. An unplanned stop is a stop that is not
	// visited by a vehicle.
	IsPlanned() bool

	// ModelStop returns the ModelStop that is the basis of the SolutionStop.
	ModelStop() ModelStop
	// ModelStopIndex is the index of the ModelStop in the Model.
	ModelStopIndex() int

	// Next returns the next stop the vehicle will visit after the stop. If
	// the stop is the last stop of a vehicle, the solution stop itself is
	// returned. If the stop is unplanned, the next stop has no semantic
	// meaning and the stop itself is returned.
	Next() SolutionStop
	// NextIndex returns the index of the next solution stop the vehicle will
	// visit after the stop. If the stop is the last stop of a vehicle,
	// the index of the stop itself is returned. If the stop is unplanned,
	// the next stop has no semantic meaning and the index of the stop itself
	// is returned.
	NextIndex() int

	// ObjectiveData returns the value of the objective for the stop. The
	// objective value of a stop is set by the ObjectiveDataUpdater.Update
	// method of the objective. If the objective is not set on the stop,
	// nil is returned. If the stop is unplanned, the objective value has no
	// semantic meaning.
	ObjectiveData(objective ModelObjective) any
	// PlanCluster returns the SolutionPlanCluster that the stop is part of.
	PlanCluster() SolutionPlanCluster
	// Previous returns the previous stop the vehicle visited before the stop.
	// If the stop is the first stop of a vehicle, the solution stop itself is
	// returned. If the stop is unplanned, the previous stop has no semantic
	// meaning and the stop itself is returned.
	Previous() SolutionStop
	// PreviousIndex returns the index of the previous solution stop the
	// vehicle visited before the stop. If the stop is the first stop of a
	// vehicle, the index of the stop itself is returned. If the stop is
	// unplanned, the previous stop has no semantic meaning and the index of
	// the stop itself is returned.
	PreviousIndex() int

	// Slack returns the slack of the stop as a time.Duration. Slack is defined
	// as the duration you can start the invoking stop later without
	// postponing the last stop of the vehicle. If the stop is unplanned,
	// the slack has no semantic meaning. Slack is a consequence of the
	// earliest start of stops, if no earliest start is set, the slack is
	// always zero.
	Slack() time.Duration
	// SlackValue returns the slack of the stop as a float64.
	SlackValue() float64

	// Vehicle returns the SolutionVehicle that visits the stop. If the stop
	// is unplanned, the vehicle has no semantic meaning and a panic will be
	// raised.
	Vehicle() SolutionVehicle
	// VehicleIndex returns the index of the SolutionVehicle that visits the
	// stop. If the stop is unplanned, a panic will be raised.
	VehicleIndex() int

	// Solution returns the Solution that the stop is part of.
	Solution() Solution
	// Start returns the start time of the stop. If the stop is unplanned, the
	// start time has no semantic meaning.
	Start() time.Time
	// StartValue returns the start time of the stop as a float64. If the stop
	// is unplanned, the start time has no semantic meaning. The returned
	// value is the number of Model.DurationUnit units since Model.Epoch.
	StartValue() float64
	// Position returns the position of the stop in the vehicle starting with
	// 0 for the first stop. If the stop is unplanned, a panic will be raised.
	Position() int

	// TravelDuration returns the travel duration of the stop as a
	// time.Duration. If the stop is unplanned, the travel duration has no
	// semantic meaning. The travel duration is the time it takes to get to
	// the invoking stop.
	TravelDuration() time.Duration
	// TravelDurationValue returns the travel duration of the stop as a
	// float64. If the stop is unplanned, the travel duration has no semantic
	// meaning. The travel duration is the time it takes to get to the
	// invoking stop. The returned value is the number of
	// Model.DurationUnit units.
	TravelDurationValue() float64

	// Value returns the value of the expression for the stop as a float64.
	// If the stop is unplanned, the value has no semantic meaning.
	Value(expression ModelExpression) float64
}

// SolutionStops is a slice of SolutionStop.
type SolutionStops []SolutionStop
