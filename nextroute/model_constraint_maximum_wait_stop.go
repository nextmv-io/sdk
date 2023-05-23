package nextroute

import "github.com/nextmv-io/sdk/connect"

// MaximumWaitStopConstraint is a constraint that limits the time a vehicle can
// wait between two stops. Wait is defined as the time between arriving at a
// stop and starting to do whatever you need to do,
// [SolutionStop.StartValue()] - [SolutionStop.ArrivalValue()].
type MaximumWaitStopConstraint interface {
	ModelConstraint

	// Maximum returns the maximum expression which defines the maximum time a
	// vehicle can wait at a stop. Returns nil if not set.
	Maximum() StopDurationExpression
}

// NewMaximumWaitStopConstraint returns a new MaximumWaitStopConstraint. The
// maximum wait constraint limits the time a vehicle can wait between two stops.
// Wait is defined as the time between arriving at a stop and starting to do
// whatever you need to do, [SolutionStop.StartValue()] -
// [SolutionStop.ArrivalValue()].
func NewMaximumWaitStopConstraint() (MaximumWaitStopConstraint, error) {
	connect.Connect(con, &newMaximumWaitStopConstraint)
	return newMaximumWaitStopConstraint()
}
