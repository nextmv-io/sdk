package dataframe

// DataFrame is an immutable data frame that support filtering, aggregation and
// data manipulation.
type DataFrame interface {
	// Bools returns a BoolColumn identified by name, panics if not present.
	Bools(name string) BoolColumn

	// Column returns a column identified by name, panics if not present.
	Column(name string) AnyColumn
	// Columns returns all columns present in the dataframe.
	Columns() AnyColumns

	// Distinct returns a new DataFrame that only contains unique rows with
	// respect to the specified columns. If no columns are given Distinct will
	// return rows where allow columns are unique.
	Distinct(columns ...Column) DataFrame

	// Filter returns a new filtered DataFrame according to the filter.
	Filter(filter Filter) DataFrame
	// Floats returns a FloatColumn identified by name, panics if not present.
	Floats(name string) FloatColumn

	// GroupBy groups rows together for which the values of specified columns
	// are the same.
	GroupBy(columns Columns) Groups

	// HasColumn reports if a columns with name is present in the dataframe.
	HasColumn(name string) bool

	// Ints returns a IntColumn identified by name, panics if not present.
	Ints(name string) IntColumn
	// IsBools returns true if column by name is of type Bool, otherwise false.
	IsBools(name string) bool
	// IsInts returns true if column by name is of type Int, otherwise false.
	IsInts(name string) bool
	// IsFloats returns true if column by name is of type floats, otherwise
	// false.
	IsFloats(name string) bool
	// IsStrings returns true if column by name is of type String, otherwise
	// false.
	IsStrings(name string) bool

	// Len returns the number of rows in the dataframe.
	Len() int

	// Select returns a new dataframe containing only the specified columns.
	Select(columns ...Column) DataFrame
	// Strings returns a StringColumn identified by name, panics if not present.
	Strings(name string) StringColumn
}

// DataFrames is the slice of DataFrame instances.
type DataFrames []DataFrame
