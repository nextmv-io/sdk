package nextroute

// Copier is the interface that all objects that can be copied must implement.
type Copier interface {
	// Copy returns a copy of the object.
	Copy() Copier
}

// ConstraintDataUpdater is the interface than can be used by a constraint if
// it wants to store data with each stop in a solution.
type ConstraintDataUpdater interface {
	// Update is called when a stop is added to a solution. The solutionStop
	// s has all it's expression values set and this function can use them to
	// update the constraint data for the stop. The data returned can be used
	// by the estimate function and can be retrieved by the
	// SolutionStop.ConstraintValue function.
	Update(s SolutionStop) Copier
}

// RegisteredModelExpressions is the interface that exposes the expressions
// that should be registered with the model.
type RegisteredModelExpressions interface {
	// ModelExpressions registers the expressions that are registered with the
	// model.
	ModelExpressions() ModelExpressions
}

// ModelConstraint is the interface that all constraints must implement.
// Constraints are used to estimate if a move is allowed and can be used to
// check if a solution is valid after a move is executed or plan clusters have
// been unplanned.
type ModelConstraint interface {
	RegisteredModelExpressions
	// CheckedAt returns when the constraint should be checked. A constraint can
	// be checked at each stop, each vehicle or each solution. If the constraint
	// is never checked it relies on its estimate of allowed moves to be
	// correct. Estimates are only used when planning plan-clusters so if
	// un-planning plan clusters can result in a solution that violates the
	// constraint then the constraint must be checked at each solution or stop
	// or vehicle.
	CheckedAt() CheckedAt

	// DoesStopHaveViolations returns true if the stop violates the constraint.
	// The stop is not allowed to be nil. The stop must be part of the solution.
	// This method is only called if CheckedAt returns AtEachStop.
	DoesStopHaveViolations(stop SolutionStop) bool
	// DoesVehicleHaveViolations returns true if the vehicle violates the
	// constraint. The vehicle is not allowed to be nil. The vehicle must be
	// part of the solution. This method is only called if CheckedAt returns
	// AtEachVehicle.
	DoesVehicleHaveViolations(vehicle SolutionVehicle) bool
	// DoesSolutionHaveViolations returns true if the solution violates the
	// constraint. The solution is not allowed to be nil. This method is only
	// called if CheckedAt returns AtEachSolution.
	DoesSolutionHaveViolations(solution Solution) bool

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
