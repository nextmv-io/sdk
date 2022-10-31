package dataframe

import "fmt"

// AggregationType defines the types of aggregations available in DataFrame.
type AggregationType string

// Types of aggregations in DataFrame.
const (
	// Maximum aggregation takes the maximum of all values.
	Maximum = "maximum"
	// Minimum aggregation takes the minimum of all values.
	Minimum = "minimum"
	// Sum aggregation takes the sum of all values.
	Sum = "sum"
)

// Aggregation defines how to aggregate rows of a group of rows in a Groups
// instance.
type Aggregation interface {
	fmt.Stringer

	// AggregationType return the type of aggregation
	AggregationType() AggregationType

	// Column return the column the aggregation will be applied to.
	Column() Column

	// As return the name to be used to identify the newly created column
	// containing the aggregated value.
	As() string
}

// Aggregations is the slice of Aggregation instances.
type Aggregations []Aggregation

// BoolAggregations defines the possible aggregations which can be applied on
// columns of type Bool
type BoolAggregations interface {
	// NewMaximum creates an aggregation which reports the maximum value using
	// name as. True is considered larger than false.
	NewMaximum(as string) Aggregation
	// NewMinimum creates an aggregation which reports the minimum value using
	// name as. True is considered larger than false.
	NewMinimum(as string) Aggregation
}

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
