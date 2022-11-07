package dataframe

import "fmt"

// Aggregation defines how to aggregate rows of a group of rows in a Groups
// instance.
type Aggregation interface {
	fmt.Stringer

	// Column returns the column the aggregation will be applied to.
	Column() Column

	// As returns the column to be used to identify the newly created column.
	// containing the aggregated value.
	As() Column
}

// Aggregations is the slice of Aggregation instances.
type Aggregations []Aggregation

// NumericAggregations defines the possible aggregations which can be applied on
// columns of type Float and Int.
type NumericAggregations interface {
	// Max creates an aggregation which reports the maximum value using
	// name as.
	Max(as string) Aggregation
	// Min creates an aggregation which reports the minimum value using
	// name as.
	Min(as string) Aggregation
	// Sum creates an aggregation which reports the sum of values using
	// name as.
	Sum(as string) Aggregation
}
