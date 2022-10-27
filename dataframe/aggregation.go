package dataframe

import "fmt"

type AggregationType string

const (
	Maximum = "maximum"
	Minimum = "minimum"
	Sum     = "sum"
)

type Aggregation interface {
	fmt.Stringer

	AggregationType() AggregationType

	Column() Column

	As() string
}

type Aggregations []Aggregation

type BoolAggregations interface {
	Maximum(string) Aggregation
	Minimum(string) Aggregation
}

type NumericAggregations interface {
	Maximum(string) Aggregation
	Minimum(string) Aggregation
	Sum(string) Aggregation
}

type StringAggregations interface {
}
