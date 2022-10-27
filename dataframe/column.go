package dataframe

import "fmt"

type DataType string

const (
	Bool   DataType = "bool"
	Int             = "int"
	Float           = "float"
	String          = "string"
)

type Column interface {
	fmt.Stringer

	Name() string

	DataType() DataType
}

type Columns []Column

type AnyColumn interface {
	Column

	ToBoolColumn() (BoolColumn, error)
	ToIntColumn() (IntColumn, error)
	ToFloatColumn() (FloatColumn, error)
	ToStringColumn() (StringColumn, error)
}

type AnyColumns []AnyColumn

type BoolColumn interface {
	Column
	BoolAggregations

	IsFalse() Filter
	IsTrue() Filter

	ItemAt(i int) bool

	Slice() []bool
}

type FloatColumn interface {
	Column
	NumericAggregations

	IsInRange(min, max float64) Filter

	ItemAt(i int) float64

	Slice() []float64
}

type IntColumn interface {
	Column
	NumericAggregations

	IsInRange(min, max int) Filter

	ItemAt(i int) int

	Slice() []int
}

type StringColumn interface {
	Column
	StringAggregations

	Equals(string) Filter

	ItemAt(i int) *string

	Slice() []*string
}
