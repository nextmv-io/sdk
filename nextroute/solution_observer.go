package nextroute

// SolutionObserver is an interface that can be implemented to observe the
// solution manipulation process.
type SolutionObserver interface {
	// OnNewSolution is called when a new solution is going to be created.
	OnNewSolution(model Model)
	// OnNewSolutionCreated is called when a new solution has been created.
	OnNewSolutionCreated(solution Solution)

	// OnCopySolution is called when a solution is going to be copied.
	OnCopySolution(solution Solution)
	// OnCopiedSolution is called when a solution has been copied.
	OnCopiedSolution(solution Solution)

	// OnCheckConstraint is called when a constraint is going to be checked.
	OnCheckConstraint(
		constraint ModelConstraint,
		violation CheckedAt,
	)
	// OnCheckedConstraint is called when a constraint has been checked.
	OnCheckedConstraint(
		constraint ModelConstraint,
		feasible bool,
	)
	// OnEstimateIsViolated is called when the delta constraint is going to be
	// estimated if it will be violated
	OnEstimateIsViolated(
		constraint ModelConstraint,
	)
	// OnEstimatedIsViolated is called when the delta constraint score
	// has been estimated.
	OnEstimatedIsViolated(
		constraint ModelConstraint,
		isViolated bool,
		planPositionsHint StopPositionsHint,
	)
	// OnEstimateDeltaObjectiveScore is called when the delta objective score is
	// going to be estimated.
	OnEstimateDeltaObjectiveScore()
	// OnEstimatedDeltaObjectiveScore is called when the delta objective score
	// has been estimated.
	OnEstimatedDeltaObjectiveScore(
		estimate float64,
	)
	// OnBestMove is called when the solution is asked for it's best move.
	OnBestMove(solution Solution)
	// OnBestMoveFound is called when the solution has found it's best move.
	OnBestMoveFound(move SolutionMove)

	// OnPlan is called when a move is going to be planned.
	OnPlan(move SolutionMove)
	// OnPlanFailed is called when a move has failed to be planned.
	OnPlanFailed(move SolutionMove)
	// OnPlanSucceeded is called when a move has succeeded to be planned.
	OnPlanSucceeded(move SolutionMove)
}

// SolutionObservers is a slice of SolutionObserver.
type SolutionObservers []SolutionObserver

// SolutionUnPlanObserver is an interface that can be implemented to observe the
// plan units un-planning process.
type SolutionUnPlanObserver interface {
	// OnUnPlan is called when a planUnit is going to be un-planned.
	OnUnPlan(planUnit SolutionPlanStopsUnit)
	// OnUnPlanFailed is called when a planUnit has failed to be un-planned.
	OnUnPlanFailed(planUnit SolutionPlanStopsUnit)
	// OnUnPlanSucceeded is called when a planUnit has succeeded to be un-planned.
	OnUnPlanSucceeded(planUnit SolutionPlanStopsUnit)
}

// OnEstimatedIsViolatedWithStopObserver is called when the delta constraint score has
// been estimated. The stop is provided to check on which stop the constraint
// is violated or not.
type OnEstimatedIsViolatedWithStopObserver interface {
	// OnEstimatedIsViolatedWithStop is called when the delta constraint score has
	// been estimated.
	OnEstimatedIsViolatedWithStop(
		constraint ModelConstraint,
		isViolated bool,
		planPositionsHint StopPositionsHint,
		stop ModelStop,
	)
}

// OnConstraintCheckObserver is an interface to observe the constraint check.
type OnConstraintCheckObserver interface {
	// OnCheckedConstraintWithStop is called when a stop constraint has been checked.
	OnCheckedConstraintWithStop(
		constraint ModelConstraint,
		feasible bool,
		stop ModelStop,
	)

	// OnCheckedConstraintWithVehicle is called when a vehicle constraint has been checked.
	OnCheckedConstraintWithVehicle(
		constraint ModelConstraint,
		feasible bool,
		vehicle SolutionVehicle,
	)
}

// OnPlanFailedWithConstraintObserver is an interface to observe the constraint check.
// The stop is provided to check on which stop the constraint is violated.
type OnPlanFailedWithConstraintObserver interface {
	// OnPlanFailedWithConstraint is called when a move has failed to be planned.
	OnPlanFailedWithConstraint(
		stop ModelStop,
		constraint ModelConstraint,
	)
}

// OnPlanFailedWithConstraintObservers is a slice of OnPlanFailedWithConstraintObserver.
type OnPlanFailedWithConstraintObservers []OnPlanFailedWithConstraintObserver

// OnConstraintCheckObservers is a slice of OnConstraintCheckObserver.
type OnConstraintCheckObservers []OnConstraintCheckObserver

// OnEstimatedIsViolatedWithStopObservers is a slice of OnEstimatedIsViolatedWithStopObserver.
type OnEstimatedIsViolatedWithStopObservers []OnEstimatedIsViolatedWithStopObserver

// SolutionUnPlanObservers is a slice of SolutionUnPlanObserver.
type SolutionUnPlanObservers []SolutionUnPlanObserver

// SolutionObserved is an interface that can be implemented to observe the
// solution manipulation process.
type SolutionObserved interface {
	SolutionObserver
	SolutionUnPlanObserver
	OnEstimatedIsViolatedWithStopObserver
	OnConstraintCheckObserver
	OnPlanFailedWithConstraintObserver

	// AddSolutionObserver adds the given solution observer to the solution
	// observed.
	AddSolutionObserver(observer SolutionObserver)

	// AddSolutionUnPlanObserver adds the given solution un-plan observer to the
	// solution observed.
	AddSolutionUnPlanObserver(observer SolutionUnPlanObserver)

	// AddOnConstraintCheckObserver adds the given observer to the solution.
	AddOnConstraintCheckObserver(observer OnConstraintCheckObserver)

	// AddOnPlanFailedWithConstraintObserver adds the given observer to the solution.
	AddOnPlanFailedWithConstraintObserver(observer OnPlanFailedWithConstraintObserver)

	// AddOnEstimatedIsViolatedWithStopObserver adds the given observer to the solution.
	AddOnEstimatedIsViolatedWithStopObserver(observer OnEstimatedIsViolatedWithStopObserver)

	// RemoveSolutionObserver remove the given solution observer from the
	// solution observed.
	RemoveSolutionObserver(observer SolutionObserver)

	// RemoveSolutionUnPlanObserver remove the given solution un-plan observer
	// from the solution observed.
	RemoveSolutionUnPlanObserver(observer SolutionUnPlanObserver)

	// SolutionObservers returns the solution observers.
	SolutionObservers() SolutionObservers

	// SolutionUnPlanObservers returns the solution un-plan observers.
	SolutionUnPlanObservers() SolutionUnPlanObservers
}
