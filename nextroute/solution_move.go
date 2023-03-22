package nextroute

import (
	"context"

	"github.com/nextmv-io/sdk/connect"
)

// NewNotExecutableMove returns a new empty move. An empty move is a move that does not
// change the solution and it is not executable.
func NewNotExecutableMove() Move {
	connect.Connect(con, &newNotExecutableMove)
	return newNotExecutableMove()
}

// Move is a move in a solution. A move is a change in the solution that can be
// executed. A move can be executed if it is executable.
type Move interface {
	// AfterStop returns the stop after which the first to be planned stop
	// is supposed to be planned. The first to be planned stop is the first
	// stop in the stop positions which is the first stop in the cluster.
	AfterStop() SolutionStop

	// BeforeStop returns the stop before which the last to be planned stop
	// is supposed to be planned. The last to be planned stop is the last
	// stop in the stop positions which is the last stop in the cluster.
	BeforeStop() SolutionStop

	// Execute executes the move. Returns true if the move was executed
	// successfully, false if the move was not executed successfully. A
	// move is not successful if it did not result in a change in the
	// solution without violating any hard constraints. A move can be
	// marked executable even if it is not successful in executing.
	Execute(context.Context) (bool, error)

	// IsExecutable returns true if the move is executable, false if the
	// move is not executable. A move is executable if the estimates believe
	// the move will result in a change in the solution without violating
	// any hard constraints.
	IsExecutable() bool

	// IsImprovement returns true if the move is executable and the move
	// has a value less than zero, false if the move is not executable or
	// the move has a value of zero or greater than zero.
	IsImprovement() bool

	// PlanCluster returns the plan cluster that is affected by the move.
	PlanCluster() SolutionPlanCluster

	// StopPositions returns the stop positions that define the move and
	// how it will change the solution.
	StopPositions() StopPositions

	// TakeBest returns the best move between the given move and the
	// current move. The best move is the move with the lowest score. If
	// the scores are equal, a random uniform distribution is used to
	// determine the move to use.
	TakeBest(that Move) Move

	// Value returns the score of the move. The score is the difference
	// between the score of the solution before the move and the score of
	// the solution after the move. The score is based on the estimates and
	// the actual score of the solution after the move should be retrieved
	// using Solution.Score after the move has been executed.
	Value() float64

	// ValueSeen returns the number of times the value of this move has been
	// seen by the estimates. A tie-breaker is a mechanism used to resolve
	// situations where multiple moves have the same value. In cases where the
	// same value is seen multiple times, a tie-breaker is applied to ensure
	// that each option has an equal chance of being selected.
	ValueSeen() int

	// Vehicle returns the vehicle, if known, that is affected by the move. If
	// not known, nil is returned.
	Vehicle() SolutionVehicle
}

// StopPosition is the definition of the change in the solution for a
// specific stop. The change is defined by a BeforeStop and a Stop. The
// BeforeStop is a stop which is already part of the solution (it is planned)
// and the Stop is a stop which is not yet part of the solution (it is not
// planned). A stop position states that the stop should be moved from the
// unplanned set to the planned set by positioning it directly before the
// BeforeStop.
type StopPosition interface {
	// AfterStop returns the stop after which Stop will be inserted.
	AfterStop() SolutionStop

	// BeforeStop returns the stop which is already part of the solution.
	BeforeStop() SolutionStop
	// Stop returns the stop which is not yet part of the solution.
	Stop() SolutionStop
}

// StopPositions is a list of stop positions.
type StopPositions []StopPosition
