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
	OnBestMoveFound(move Move)

	// OnPlan is called when a move is going to be planned.
	OnPlan(move Move)
	// OnPlanFailed is called when a move has failed to be planned.
	OnPlanFailed(move Move)
	// OnPlanSucceeded is called when a move has succeeded to be planned.
	OnPlanSucceeded(move Move)
}

// SolutionObservers is a slice of SolutionObserver.
type SolutionObservers []SolutionObserver

// SolutionObserved is an interface that can be implemented to observe the
// solution manipulation process.
type SolutionObserved interface {
	SolutionObserver
	// AddSolutionObserver adds the given solution observer to the solution
	// observed.
	AddSolutionObserver(observer SolutionObserver)

	// SolutionObservers returns the solution observers.
	SolutionObservers() SolutionObservers
}
