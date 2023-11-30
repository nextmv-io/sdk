package dataframe

import "fmt"

// Aggregation defines how to aggregate rows of a group of rows in a Groups
// instance.
//
// Deprecated: This package is deprecated and will be removed in a future.
type Aggregation interface {
	fmt.Stringer

	// Column returns the column the aggregation will be applied to.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	Column() Column

	// As returns the column to be used to identify the newly created column.
	// containing the aggregated value.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	As() Column
}

// Aggregations is the slice of Aggregation instances.
//
// Deprecated: This package is deprecated and will be removed in a future.
type Aggregations []Aggregation

// NumericAggregations defines the possible aggregations which can be applied on
// columns of type Float and Int.
//
// Deprecated: This package is deprecated and will be removed in a future.
type NumericAggregations interface {
	// Max creates an aggregation which reports the maximum value using
	// name as.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	Max(as string) Aggregation
	// Min creates an aggregation which reports the minimum value using
	// name as.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	Min(as string) Aggregation
	// Sum creates an aggregation which reports the sum of values using
	// name as.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	Sum(as string) Aggregation
}
