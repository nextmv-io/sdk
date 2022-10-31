package dataframe

import "fmt"

// DataType defines the types of colums available in DataFrame.
type DataType string

// Types of data in DataFrame.
const (
	// Bool type representing boolean true and false values.
	Bool DataType = "bool"
	// Int type representing int values.
	Int = "int"
	// Float type representing float64 values.
	Float = "float"
	// String type representing string values.
	String = "string"
)

// Column is a single column in a DataFrame instance. It is identified by its
// name and has a DataType.
type Column interface {
	fmt.Stringer

	// Name returns the name of the column, the name is the unique identifier
	// of the column within a DataFrame instance.
	Name() string

	// DataType returns the type of the column.
	DataType() DataType
}

// Columns is the slice of Column instances.
type Columns []Column

// AnyColumn is the interface to cast a column to a typed column.
type AnyColumn interface {
	Column

	// ToBoolColumn return the column as a BoolColumn, panics if the column
	// is not of type Bool.
	ToBoolColumn() BoolColumn
	// ToIntColumn return the column as a IntColumn, panics if the column
	// is not of type Int.
	ToIntColumn() IntColumn
	// ToFloatColumn return the column as a FloatColumn, panics if the column
	// is not of type Float.
	ToFloatColumn() FloatColumn
	// ToStringColumn return the column as a StringColumn, panics if the column
	// is not of type String.
	ToStringColumn() StringColumn
}

// AnyColumns is the slice of AnyColumn instances.
type AnyColumns []AnyColumn

// BoolColumn is the typed column of type Bool
type BoolColumn interface {
	Column
	BoolAggregations

	// NewIsFalse creates a filter to filter all rows having value false.
	NewIsFalse() Filter
	// NewIsTrue creates a filter to filter all rows having value true.
	NewIsTrue() Filter

	// ItemAt return the value in row i, panics if out of bound.
	ItemAt(i int) bool

	// Slice returns all the values in the column.
	Slice() []bool
}

// FloatColumn is the typed column of type Float
type FloatColumn interface {
	Column
	NumericAggregations

	// NewIsInRange creates a filter to filter all rows within range [min, max].
	NewIsInRange(min, max float64) Filter

	// ItemAt return the value in row i, panics if out of bound.
	ItemAt(i int) float64

	// Slice returns all the values in the column.
	Slice() []float64
}

// IntColumn is the typed column of type Int
type IntColumn interface {
	Column
	NumericAggregations

	// NewIsInRange creates a filter to filter all value within range [min, max].
	NewIsInRange(min, max int) Filter

	// ItemAt return the value in row i, panics if out of bound.
	ItemAt(i int) int

	// Slice returns all the values in the column.
	Slice() []int
}

// StringColumn is the typed column of type String
type StringColumn interface {
	Column

	// NewEquals creates a filter to filter all rows having value value.
	NewEquals(value string) Filter

	// ItemAt return the value in row i, panics if out of bound.
	ItemAt(i int) *string

	// Slice returns all the values in the column.
	Slice() []*string
}
