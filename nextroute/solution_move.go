package nextroute

import "github.com/nextmv-io/sdk/connect"

// NewEmptyMove returns a new empty move. An empty move is a move that does not
// change the solution and it is not executable.
func NewEmptyMove() Move {
	connect.Connect(con, &newEmptyMove)
	return newEmptyMove()
}

// Move is a move in a solution. A move is a change in the solution that can be
// executed. A move can be executed if it is executable.
type Move interface {
	// Execute executes the move. Returns true if the move was executed
	// successfully, false if the move was not executed successfully. A
	// move is not successful if it did not result in a change in the
	// solution without violating any hard constraints. A move can be
	// marked executable even if it is not successful in executing.
	Execute() (bool, error)

	// IsExecutable returns true if the move is executable, false if the
	// move is not executable. A move is executable if the estimates believe
	// the move will result in a change in the solution without violating
	// any hard constraints.
	IsExecutable() bool

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
}

// StopPosition is the definition of the change in the solution for a
// specific stop. The change is defined by a BeforeStop and a Stop. The
// BeforeStop is a stop which is already part of the solution (it is planned)
// and the Stop is a stop which is not yet part of the solution (it is not
// planned). A stop position states that the stop should be moved from the
// unplanned set to the planned set by positioning it directly before the
// BeforeStop.
type StopPosition interface {
	// BeforeStop returns the stop which is already part of the solution.
	BeforeStop() SolutionStop
	// Stop returns the stop which is not yet part of the solution.
	Stop() SolutionStop
}

// StopPositions is a list of stop positions.
type StopPositions []StopPosition
