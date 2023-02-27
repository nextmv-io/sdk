package alns

import (
	"math/rand"
	"time"
)

// Solver is the interface for the Adaptive Local Neighborhood Search algorithm.
type Solver[T baseSolution[T]] interface {
	SolveObserved[T]

	// BestSolution returns the best solution found so far.
	BestSolution() T

	// Restart sets the current work solution to the best solution.
	Restart()

	// Solve starts the solving process using the given options.
	Solve(solveOptions SolveOptions) (T, error)
	// SolveOperators returns the solve-operators used by the solver.
	SolveOperators() SolveOperators[T]

	// WorkSolution returns the current work solution.
	WorkSolution() T

	// Random returns the random number generator used by the solver.
	Random() *rand.Rand
	// Register registers a parameter with the solver. Any parameter used by
	// a solve-operator must be registered with the solver.
	Register(parameter SolveParameter[T])
}

// NewSolver creates a new solver for the given work solution.
func NewSolver[T baseSolution[T]](
	workSolution T,
	options ...OptionSolver[T],
) (Solver[T], error) {
	solver := &solveImpl[T]{
		workSolution: workSolution.Copy(),
		bestSolution: workSolution.Copy(),
		random:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	for _, option := range options {
		if err := option(solver); err != nil {
			return nil, err
		}
	}
	return solver, nil
}

// OptionSolver is a function that configures a solver.
type OptionSolver[T baseSolution[T]] func(s Solver[T]) error

// AddSolveOperator adds a solve-operator to the solver.
func AddSolveOperator[T baseSolution[T]](
	solveOperator SolveOperator[T],
) OptionSolver[T] {
	return func(s Solver[T]) error {
		s.(*solveImpl[T]).solveOperators = append(
			s.(*solveImpl[T]).solveOperators,
			solveOperator,
		)
		for _, parameter := range solveOperator.SolveParameters() {
			s.Register(parameter)
		}
		return nil
	}
}

// AddObserver adds an observer to the solver.
func AddObserver[T baseSolution[T]](
	observer SolveObserver[T],
) OptionSolver[T] {
	return func(s Solver[T]) error {
		s.(*solveImpl[T]).observers = append(
			s.(*solveImpl[T]).observers,
			observer,
		)
		return nil
	}
}

// RandomSource sets the random number generator used by the solver.
func RandomSource[T baseSolution[T]](
	random *rand.Rand,
) OptionSolver[T] {
	return func(s Solver[T]) error {
		s.(*solveImpl[T]).random = random
		return nil
	}
}

// Seed sets the seed of the random number generator used by the solver.
func Seed[T baseSolution[T]](seed int64) OptionSolver[T] {
	return func(s Solver[T]) error {
		s.(*solveImpl[T]).random.Seed(seed)
		return nil
	}
}

type solveImpl[T Solution[T]] struct {
	solveObservedImpl[T]
	workSolution   Solution[T]
	bestSolution   Solution[T]
	solveOperators SolveOperators[T]
	parameters     SolveParameters[T]
	random         *rand.Rand
}

func (s *solveImpl[T]) Random() *rand.Rand {
	return s.random
}

func (s *solveImpl[T]) AddSolveObserver(observer SolveObserver[T]) {
	s.observers = append(s.observers, observer)
}

func (s *solveImpl[T]) Register(parameter SolveParameter[T]) {
	s.parameters = append(s.parameters, parameter)
}

func (s *solveImpl[T]) SolveOperators() SolveOperators[T] {
	solveOperators := make(SolveOperators[T], len(s.solveOperators))
	copy(solveOperators, s.solveOperators)
	return solveOperators
}

func (s *solveImpl[T]) Restart() {
	s.workSolution = s.bestSolution.Copy()
}

func (s *solveImpl[T]) BestSolution() T {
	return s.bestSolution.(T)
}

func (s *solveImpl[T]) WorkSolution() T {
	return s.workSolution.(T)
}

func (s *solveImpl[T]) invoke(
	solveOperator SolveOperator[T],
	solveInformation *solveInformationImpl[T],
) {
	solveInformation.solveOperators = append(
		solveInformation.solveOperators,
		solveOperator,
	)
	if s.Random().Float64() <= solveOperator.Probability() {
		s.OnExecute(solveInformation)
		solveOperator.Execute(solveInformation)
		s.OnExecuted(solveInformation)

		if solveOperator.CanResultInImprovement() {
			delta := s.workSolution.Score() - s.bestSolution.Score()

			if delta < 0.0 {
				solveInformation.deltaScore += delta
				s.OnImprovement(solveInformation)
				s.bestSolution = s.workSolution.Copy()
				for _, st := range s.SolveOperators() {
					if interested, ok := st.(InterestedInBetterSolution[T]); ok {
						interested.OnBetterSolution(solveInformation)
					}
				}
			}
		}
	}
}

func (s *solveImpl[T]) Solve(
	solveOptions SolveOptions,
) (T, error) {
	start := time.Now()

	solveInformation := &solveInformationImpl[T]{
		iteration:      0,
		solver:         s,
		solveOperators: make(SolveOperators[T], 0, len(s.solveOperators)),
		start:          time.Now(),
	}

	for _, solveOperator := range s.SolveOperators() {
		if interested, ok := solveOperator.(InterestedInStartSolve[T]); ok {
			interested.OnStartSolve(solveInformation)
		}
	}
	for iteration := 0; iteration < solveOptions.Iterations &&
		time.Since(start) < solveOptions.MaximumDuration; iteration++ {
		solveInformation.iteration = iteration
		solveInformation.deltaScore = 0.0
		solveInformation.solveOperators = make(
			SolveOperators[T],
			0,
			len(s.solveOperators),
		)

		s.OnIteration(solveInformation)

		for _, solveOperator := range s.SolveOperators() {
			s.invoke(solveOperator, solveInformation)
		}

		for _, parameter := range s.parameters {
			parameter.Update(solveInformation)
		}

		s.OnIterated(solveInformation)
	}
	s.OnEnd(solveInformation)

	return s.bestSolution.(T), nil
}
