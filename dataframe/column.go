package dataframe

import (
	"fmt"

	"github.com/nextmv-io/sdk/connect"
)

// Bools returns a BoolColumn identified by name.
//
// Deprecated: This package is deprecated and will be removed in the future.
func Bools(name string) BoolColumn {
	connect.Connect(con, &newBoolColumn)
	return newBoolColumn(name)
}

// Floats returns a FloatColumn identified by name.
//
// Deprecated: This package is deprecated and will be removed in the future.
func Floats(name string) FloatColumn {
	connect.Connect(con, &newFloatColumn)
	return newFloatColumn(name)
}

// Ints returns a IntColumn identified by name.
//
// Deprecated: This package is deprecated and will be removed in the future.
func Ints(name string) IntColumn {
	connect.Connect(con, &newIntColumn)
	return newIntColumn(name)
}

// Strings returns a StringColumn identified by name.
//
// Deprecated: This package is deprecated and will be removed in the future.
func Strings(name string) StringColumn {
	connect.Connect(con, &newStringColumn)
	return newStringColumn(name)
}

// DataType defines the types of colums available in DataFrame.
//
// Deprecated: This package is deprecated and will be removed in the future.
type DataType string

// Types of data in DataFrame.
//
// Deprecated: This package is deprecated and will be removed in the future.
const (
	// Bool type representing boolean true and false values.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Bool DataType = "bool"
	// Int type representing int values.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Int = "int"
	// Float type representing float64 values.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Float = "float"
	// String type representing string values.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	String = "string"
)

// Column is a single column in a DataFrame instance. It is identified by its
// name and has a DataType.
//
// Deprecated: This package is deprecated and will be removed in the future.
type Column interface {
	fmt.Stringer

	// Name returns the name of the column, the name is the unique identifier
	// of the column within a DataFrame instance.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Name() string

	// DataType returns the type of the column.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	DataType() DataType
}

// Columns is the slice of Column instances.
//
// Deprecated: This package is deprecated and will be removed in the future.
type Columns []Column

// BoolColumn is the typed column of type Bool.
//
// Deprecated: This package is deprecated and will be removed in the future.
type BoolColumn interface {
	Column

	// IsFalse creates a filter to filter all rows having value false.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	IsFalse() Filter
	// IsTrue creates a filter to filter all rows having value true.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	IsTrue() Filter

	// Value return the value at row for dataframe df,
	// panics if out of bound.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Value(df DataFrame, row int) bool

	// Values returns all the values in the column for dataframe df.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Values(df DataFrame) []bool
}

// FloatColumn is the typed column of type Float.
//
// Deprecated: This package is deprecated and will be removed in the future.
type FloatColumn interface {
	Column
	NumericAggregations

	// IsInRange creates a filter to filter all rows within range [min, max].
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	IsInRange(min, max float64) Filter

	// Value return the value at row, panics if out of bound.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Value(df DataFrame, row int) float64

	// Values returns all the values in the column.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Values(df DataFrame) []float64
}

// IntColumn is the typed column of type Int.
//
// Deprecated: This package is deprecated and will be removed in the future.
type IntColumn interface {
	Column
	NumericAggregations

	// IsInRange creates a filter to filter all value within range [min, max].
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	IsInRange(min, max int) Filter

	// Value return the value at row, panics if out of bound.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Value(df DataFrame, row int) int

	// Values returns all the values in the column.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Values(df DataFrame) []int
}

// StringColumn is the typed column of type String.
//
// Deprecated: This package is deprecated and will be removed in the future.
type StringColumn interface {
	Column

	// Equals creates a filter to filter all rows having value value.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Equals(value string) Filter

	// Value return the value at row, panics if out of bound.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Value(df DataFrame, row int) string

	// Values returns all the values in the column.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Values(df DataFrame) []string
}

var (
	newBoolColumn   func(string) BoolColumn
	newFloatColumn  func(string) FloatColumn
	newIntColumn    func(string) IntColumn
	newStringColumn func(string) StringColumn
)
