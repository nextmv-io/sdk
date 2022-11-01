package dataframe

import "fmt"

// Aggregation defines how to aggregate rows of a group of rows in a Groups
// instance.
type Aggregation interface {
	fmt.Stringer

	// Column return the column the aggregation will be applied to.
	Column() Column

	// As return the column  to be used to identify the newly created column
	// containing the aggregated value.
	As() Column
}

// Aggregations is the slice of Aggregation instances.
type Aggregations []Aggregation

// NumericAggregations defines the possible aggregations which can be applied on
// columns of type Float and Int.
type NumericAggregations interface {
	// NewMaximum creates an aggregation which reports the maximum value using
	// name as.
	NewMaximum(as string) Aggregation
	// NewMinimum creates an aggregation which reports the minimum value using
	// name as.
	NewMinimum(as string) Aggregation
	// NewSum creates an aggregation which reports the sum of values using
	// name as.
	NewSum(as string) Aggregation
}
