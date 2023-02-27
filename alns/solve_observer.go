package alns

// SolveObserver defines the interface for a solve observer.
type SolveObserver[T baseSolution[T]] interface {
	// OnEnd is called when the solver has ended.
	OnEnd(SolveInformation[T])
	// OnExecute is called when a solve-operator is about to be executed.
	OnExecute(SolveInformation[T])
	// OnExecuted is called when a solve-operator has been executed.
	OnExecuted(SolveInformation[T])
	// OnImprovement is called when a solve-operator has resulted in an
	// improvement compared to the best solution.
	OnImprovement(SolveInformation[T])
	// OnIterated is called when a new iteration has been completed.
	OnIterated(SolveInformation[T])
	// OnIteration is called when a new iteration is started.
	OnIteration(SolveInformation[T])
}

// SolveObservers defines a slice of solve observers.
type SolveObservers[T baseSolution[T]] []SolveObserver[T]

// SolveObserved defines the interface for a solve-observed.
type SolveObserved[T baseSolution[T]] interface {
	SolveObserver[T]

	// AddSolverObserver adds a solve observer to the solve-observed.
	AddSolverObserver(observer SolveObserver[T])
}

type solveObservedImpl[T baseSolution[T]] struct {
	observers SolveObservers[T]
}

func (s *solveObservedImpl[T]) AddSolverObserver(
	observer SolveObserver[T],
) {
	s.observers = append(s.observers, observer)
}

func (s *solveObservedImpl[T]) OnIteration(
	solveInformation SolveInformation[T],
) {
	for _, observer := range s.observers {
		observer.OnIteration(solveInformation)
	}
}

func (s *solveObservedImpl[T]) OnIterated(
	solveInformation SolveInformation[T],
) {
	for _, observer := range s.observers {
		observer.OnIterated(solveInformation)
	}
}

func (s *solveObservedImpl[T]) OnExecute(
	solveInformation SolveInformation[T],
) {
	for _, observer := range s.observers {
		observer.OnExecute(solveInformation)
	}
}

func (s *solveObservedImpl[T]) OnExecuted(
	solveInformation SolveInformation[T],
) {
	for _, observer := range s.observers {
		observer.OnExecuted(solveInformation)
	}
}

func (s *solveObservedImpl[T]) OnImprovement(
	solveInformation SolveInformation[T],
) {
	for _, observer := range s.observers {
		observer.OnImprovement(solveInformation)
	}
}

func (s *solveObservedImpl[T]) OnEnd(
	solveInformation SolveInformation[T],
) {
	for _, observer := range s.observers {
		observer.OnEnd(solveInformation)
	}
}
