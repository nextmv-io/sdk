package nextroute

import "github.com/nextmv-io/sdk/connect"

// SolutionStopGenerator is an iterator of solution stops.
type SolutionStopGenerator interface {
	Next() SolutionStop
}

// NewSolutionStopGenerator return a solution stop iterator of a move.
// If startAtFirst is true, the first stop will be first stop of the vehicle.
// If endAtLast is true, the last stop will be the last stop of the vehicle.
//
// For example adding sequence A, B in a sequence 1 -> 2 -> 3 -> 4 -> 5 -> 6
// where A goes before 4 and B goes before 5 will generate the following
// solution stops: 3 -> A -> 4 -> B -> 5
// If startsAtFirst is true, the solution stops will start with 1:
// 1 -> 2 -> 3 -> A -> 4 -> B -> 5
// If endAtLast is also true, the solution stops will end with 6:
// 1 -> 2 -> 3 -> A -> 4 -> B -> 5 -> 6.
//
// For example:
//
//	   generator := NewSolutionStopGenerator(move, false, true)
//
//		  for solutionStop := generator.Next(); solutionStop != nil; solutionStop = generator.Next() {
//			  // Do something with solutionStop
//	   }
func NewSolutionStopGenerator() SolutionStopGenerator {
	connect.Connect(con, &newSolutionStopGenerator)
	return newSolutionStopGenerator()
}
