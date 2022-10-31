package dataframe

// Groups contains groups of rows produced by DataFrame.GroupBy function.
type Groups interface {
	// Aggregate applies the given aggregations to all row groups in the
	// Groups.
	Aggregate(aggregations Aggregations) DataFrame

	// DataFrames returns a slice of DataFrame where each frame represents
	// the content of one group.
	DataFrames() DataFrames
}
