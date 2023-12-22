package dataframe

// Groups contains groups of rows produced by DataFrame.GroupBy function.
//
// Deprecated: This package is deprecated and will be removed in the next major release.
type Groups interface {
	// Aggregate applies the given aggregations to all row groups in the
	// Groups and returns DataFrame instance where each row corresponds
	// to each group.
	//
	// Deprecated: This package is deprecated and will be removed in the next major release.
	Aggregate(aggregations ...Aggregation) DataFrame

	// DataFrames returns a slice of DataFrame where each frame represents
	// the content of one group.
	//
	// Deprecated: This package is deprecated and will be removed in the next major release.
	DataFrames() DataFrames
}
