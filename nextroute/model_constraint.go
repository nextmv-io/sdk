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

// ModelConstraint is the interface that all constraints must implement.
// Constraints are used to estimate if a move is allowed and can be used to
// check if a solution is valid after a move is executed or plan clusters have
// been unplanned.
type ModelConstraint interface {
	// CanBeViolated returns true if the constraint can be violated. If the
	// constraint can not be violated no solutions can be created that violate
	// the constraint.
	CanBeViolated() bool
	// CanNotBeViolated returns true if the constraint can not be violated. See
	// CanBeViolated for more information.
	CanNotBeViolated() bool
	// CheckedAt returns when the constraint should be checked. A constraint can
	// be checked at each stop, each vehicle or each solution. If the constraint
	// is never checked it relies on its estimate of allowed moves to be
	// correct.
	CheckedAt() CheckedAt
	// ConstraintViolationToScore converts a constraint violation to a score.
	// A score is a positive number that is added to the solution score if the
	// constraint is violated. The violation is multiplied by the factor plus
	// offset if violation is greater than zero otherwise zero is returned.
	ConstraintViolationToScore(violation float64) float64

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

	// EstimateDeltaScore estimates the delta score if the stop is moved to the
	// new position described in stopPositions. The stopPositions is not
	// allowed to be nil. Should be a pure function.
	EstimateDeltaScore(
		stopPositions StopPositions,
	) (estimateDeltaScore float64, stopPositionsHint StopPositionsHint)

	// Factor returns the factor that is used to convert a constraint violation
	// to a score.
	Factor() float64

	// Index returns the index of the constraint. The index is unique for each
	// constraint. The index is used to identify the constraint in the
	// solution.
	Index() int

	// ModelExpressions returns the expressions that are used by the constraint
	// whose values are defined by the solution. The solution will calculate
	// the values and cumulative values of the expressions for each stop. These
	// values can be used to implement the constraint.
	ModelExpressions() ModelExpressions

	// Name returns the name of the constraint.
	Name() string

	// Offset returns the offset that is used to convert a constraint violation
	// to a score.
	Offset() float64

	// Score returns the score of the constraint violation. The score
	// is a positive number that is added to the solution score if the
	// constraint is violated. The violation is multiplied by the factor plus
	// offset if violation is greater than zero otherwise zero is returned.
	Score(solution Solution) float64
}

// ModelConstraints is a slice of ModelConstraint.
type ModelConstraints []ModelConstraint
