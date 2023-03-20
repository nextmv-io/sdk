package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
)

// Copier is the interface that all objects that can be copied must implement.
type Copier interface {
	// Copy returns a copy of the object.
	Copy() Copier
}

type ConstraintReporter interface {
	ReportConstraint(SolutionStop) map[string]any
}

// ConstraintDataUpdater is the interface than can be used by a constraint if
// it wants to store data with each stop in a solution.
type ConstraintDataUpdater interface {
	// UpdateConstraintData is called when a stop is added to a solution.
	// The solutionStop has all it's expression values set and this function
	// can use them to update the constraint data for the stop. The data
	// returned can be used by the estimate function and can be retrieved by the
	// SolutionStop.ConstraintValue function.
	UpdateConstraintData(s SolutionStop) Copier
}

// RegisteredModelExpressions is the interface that exposes the expressions
// that should be registered with the model.
type RegisteredModelExpressions interface {
	// ModelExpressions registers the expressions that are registered with the
	// model.
	ModelExpressions() ModelExpressions
}

// ComputationalEffort is the interface that can be implemented by a constraint
// to indicate the computational effort of checking the constraint.
type ComputationalEffort interface {
	ComputationalEffortForCheck() Cost
}

// SolutionStopViolationCheck is the interface that will be invoked on every
// update of a planned solution stop. The method DoesStopHaveViolations will be
// called to check if the stop violates the constraint. If the method returns
// true the solution will not be accepted. If un-planning can result in a
// violation of the constraint one of the violation checks must be implemented.
type SolutionStopViolationCheck interface {
	ComputationalEffort
	// DoesStopHaveViolations returns true if the stop violates the constraint.
	// The stop is not allowed to be nil. The stop must be part of the solution.
	// This method is only called if CheckedAt returns AtEachStop.
	DoesStopHaveViolations(stop SolutionStop) bool
}

// SolutionVehicleViolationCheck is the interface that will be invoked on every
// update of a last planned solution stop of a vehicle. The method
// DoesVehicleHaveViolations will be called to check if the vehicle violates
// the constraint. If the method returns true the solution will not be accepted.
// If un-planning can result in a violation of the constraint one of the
// violation checks must be implemented.
type SolutionVehicleViolationCheck interface {
	ComputationalEffort
	// DoesVehicleHaveViolations returns true if the vehicle violates the
	// constraint. The vehicle is not allowed to be nil. The vehicle must be
	// part of the solution. This method is only called if CheckedAt returns
	// AtEachVehicle.
	DoesVehicleHaveViolations(vehicle SolutionVehicle) bool
}

// SolutionViolationCheck is the interface that will be invoked once all
// updates on a solution have been executed. The method
// DoesSolutionHaveViolations will be called to check if the vehicle violates
// the constraint. If the method returns true the solution will not be accepted.
// If un-planning can result in a violation of the constraint one of the
// violation checks must be implemented.
type SolutionViolationCheck interface {
	ComputationalEffort
	// DoesSolutionHaveViolations returns true if the solution violates the
	// constraint. The solution is not allowed to be nil. This method is only
	// called if CheckedAt returns AtEachSolution.
	DoesSolutionHaveViolations(solution Solution) bool
}

// ModelConstraint is the interface that all constraints must implement.
// Constraints are used to estimate if a move is allowed and can be used to
// check if a solution is valid after a move is executed or plan clusters have
// been unplanned.
type ModelConstraint interface {
	RegisteredModelExpressions
	// EstimateIsViolated estimates if the solution is changed by the given
	// new positions described in stopPositions if it will be violated or not.
	// The stopPositions is not  allowed to be nil. Should be a pure function,
	// i.e. not change any state of the constraint. The stopPositionsHint can
	// The stopPositionsHint can be used to speed up the estimation of the
	// constraint violation.
	EstimateIsViolated(
		stopPositions StopPositions,
	) (isViolated bool, stopPositionsHint StopPositionsHint)

	// Index returns the index of the constraint. The index should be
	// unique for each constraint.
	Index() int

	// Name returns the name of the constraint.
	Name() string
}

// ModelConstraints is a slice of ModelConstraint.
type ModelConstraints []ModelConstraint

// NewModelConstraintIndex returns the next unique constraint index.
// This function can be used to create a unique index for a custom
// constraint.
func NewModelConstraintIndex() int {
	connect.Connect(con, &newModelConstraintIndex)
	return newModelConstraintIndex()
}
