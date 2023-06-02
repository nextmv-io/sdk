package nextroute

// SolutionMoveFilter is a filter for SolutionMoves. A SolutionMoveFilter
// filters out SolutionMoves that should not be evaluated.
type SolutionMoveFilter interface {
	// FilterMove returns true if the move should be filtered out.
	FilterMove(move Move) bool
}

// SolutionMoveFilters is a slice of SolutionMoveFilter.
type SolutionMoveFilters []SolutionMoveFilter
