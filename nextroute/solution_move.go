package nextroute

import (
	"context"

	"github.com/nextmv-io/sdk/connect"
)

// NewNotExecutableMove returns a new empty move. An empty move is a move that
// does not change the solution, and it is marked as not executable.
func NewNotExecutableMove() SolutionMove {
	connect.Connect(con, &newNotExecutableMove)
	return newNotExecutableMove()
}

// SolutionMove is a move in a solution.
type SolutionMove interface {
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

	// PlanUnit returns the [SolutionPlanUnit] that is affected by the move.
	PlanUnit() SolutionPlanUnit

	// TakeBest returns the best move between the given move and the
	// current move. The best move is the move with the lowest score. If
	// the scores are equal, a random uniform distribution is used to
	// determine the move to use.
	TakeBest(that SolutionMove) SolutionMove

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

	// IncrementValueSeen increments the number of times the value of this move
	// has been seen by the estimates and returns the move. A tie-breaker is a
	// mechanism used to resolve situations where multiple moves have the same
	// value. In cases where the same value is seen multiple times, a
	// tie-breaker is applied to ensure that each option has an equal chance of
	// being selected.
	IncrementValueSeen(inc int) SolutionMove
}

// SolutionMoves is a slice of SolutionMove.
type SolutionMoves []SolutionMove
