package nextroute

import (
	"math/rand"
	"time"

	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/nextroute/schema"
)

type RestartPolicy interface {
	RestartSolution(*rand.Rand, Solution) Solution
	RestartIterations() int
}

type randomRestartPolicy struct {
	iterations int
}

func (r *randomRestartPolicy) RestartSolution(random *rand.Rand, s Solution) Solution {
	solution := s.Copy()

	vehicles := solution.Vehicles()

	unplannedPlanClusters := solution.UnplannedPlanClusters()

	random.Shuffle(len(unplannedPlanClusters), func(i, j int) {
		unplannedPlanClusters[i], unplannedPlanClusters[j] = unplannedPlanClusters[j], unplannedPlanClusters[i]
	})

	for _, vehicle := range vehicles {
		for idx, unplannedPlanCluster := range unplannedPlanClusters {
			if !unplannedPlanCluster.IsPlanned() {
				move := vehicle.BestMove(unplannedPlanCluster)
				if move.IsExecutable() {
					result, err := move.Execute()
					if err != nil {
						panic(err)
					}
					if result {
						unplannedPlanClusters = append(
							unplannedPlanClusters[:idx],
							unplannedPlanClusters[idx+1:]...,
						)
						break
					}
				}
			}
		}
	}
	for _, unplannedPlanCluster := range unplannedPlanClusters {
		if !unplannedPlanCluster.IsPlanned() {
			move := solution.BestMove(unplannedPlanCluster)
			if move.IsExecutable() {
				_, err := move.Execute()
				if err != nil {
					panic(err)
				}
			}
		}
	}
	return solution
}

func (r *randomRestartPolicy) RestartIterations() int {
	return r.iterations
}

func NewRandomRestartPolicy(iterations int) RestartPolicy {
	return &randomRestartPolicy{
		iterations: iterations,
	}
}

type SolveOptions struct {
	Iterations      int           `json:"iterations"  usage:"number of iterations"`
	MaximumDuration time.Duration `json:"maximum_duration"  usage:"maximum duration of solver in seconds"`
}

// Solver is the interface for a solver.
type Solver interface {
	// Solve solves the problem usint the solve-options.
	Solve(solveOptions SolveOptions) (Solution, error)
	// SolverOptions returns the solver-options used to create the solver. The
	// returned options are a copy of the options used to create the solver.
	// They can be used to create a new solver and changes will have no effect
	// on this invoked solver.
	SolverOptions() SolverOptions

	SetStartSolution(solution Solution)

	SetRestartPolicy(restartPolicy RestartPolicy)

	// Progression returns the progression of the solver.
	Progression() []schema.JsonObjectiveElapsed
}

type IntParameterOptions struct {
	StartValue               int  `json:"start_value"  usage:"start value"`
	DeltaAfterIterations     int  `json:"delta_after_iterations"  usage:"delta after each iterations"`
	Delta                    int  `json:"delta"  usage:"delta"`
	MinValue                 int  `json:"min_value"  usage:"min value of parameter"`
	MaxValue                 int  `json:"max_value"  usage:"max value of parameter"`
	SnapBackAfterImprovement bool `json:"snap_back_after_improvement"  usage:"snap back to start value after improvement of best solution"`
	Zigzag                   bool `json:"zigzag"  usage:"zigzag between min and max value lik a jig saw"`
}

type SolverOptions struct {
	Unplan  IntParameterOptions `json:"unplan"  usage:"unplan parameter"`
	Plan    IntParameterOptions `json:"plan"  usage:"plan parameter"`
	Restart IntParameterOptions `json:"restart"  usage:"restart parameter"`
}

// SolverFactory is the interface for a solver-factory.
type SolverFactory interface {
	// NewSolver creates a new solver.
	NewSolver(model Model) (Solver, error)
}

func NewSolverFactory() SolverFactory {
	connect.Connect(con, &newSolverFactory)
	return newSolverFactory()
}

func NewSolver(
	solution Solution,
	options SolverOptions,
) (Solver, error) {
	connect.Connect(con, &newSolver)
	return newSolver(solution, options)
}
