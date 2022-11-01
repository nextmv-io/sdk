package dataframe

import (
	"fmt"

	"github.com/nextmv-io/sdk/connect"
)

// Bools returns a BoolColumn identified by name.
func Bools(name string) BoolColumn {
	connect.Connect(con, &newBoolColumn)
	return newBoolColumn(name)
}

// Floats returns a FloatColumn identified by name.
func Floats(name string) FloatColumn {
	connect.Connect(con, &newFloatColumn)
	return newFloatColumn(name)
}

// Ints returns a IntColumn identified by name.
func Ints(name string) IntColumn {
	connect.Connect(con, &newIntColumn)
	return newIntColumn(name)
}

// Strings returns a StringColumn identified by name.
func Strings(name string) StringColumn {
	connect.Connect(con, &newStringColumn)
	return newStringColumn(name)
}

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

// BoolColumn is the typed column of type Bool.
type BoolColumn interface {
	Column

	// NewIsFalse creates a filter to filter all rows having value false.
	NewIsFalse() Filter
	// NewIsTrue creates a filter to filter all rows having value true.
	NewIsTrue() Filter

	// Row return the value in row i for dataframe df,
	// panics if out of bound.
	Row(df DataFrame, i int) bool

	// Rows returns all the values in the column for dataframe df.
	Rows(df DataFrame) []bool
}

// FloatColumn is the typed column of type Float.
type FloatColumn interface {
	Column
	NumericAggregations

	// NewIsInRange creates a filter to filter all rows within range [min, max].
	NewIsInRange(min, max float64) Filter

	// Row return the value in row i, panics if out of bound.
	Row(df DataFrame, i int) float64

	// Rows returns all the values in the column.
	Rows(df DataFrame) []float64
}

// IntColumn is the typed column of type Int.
type IntColumn interface {
	Column
	NumericAggregations

	// NewIsInRange creates a filter to filter all value within range [min, max].
	NewIsInRange(min, max int) Filter

	// Row return the value in row i, panics if out of bound.
	Row(df DataFrame, i int) int

	// Rows returns all the values in the column.
	Rows(df DataFrame) []int
}

// StringColumn is the typed column of type String.
type StringColumn interface {
	Column

	// NewEquals creates a filter to filter all rows having value value.
	NewEquals(value string) Filter

	// Row return the value in row i, panics if out of bound.
	Row(df DataFrame, i int) *string

	// Rows returns all the values in the column.
	Rows(df DataFrame) []*string
}

var (
	newBoolColumn   func(string) BoolColumn
	newFloatColumn  func(string) FloatColumn
	newIntColumn    func(string) IntColumn
	newStringColumn func(string) StringColumn
)
