package dataframe

// DataFrame is an immutable data frame that support filtering, aggregation and
// data manipulation.
type DataFrame interface {
	// Column returns a column identified by name, panics if not present.
	Column(name string) Column
	// Columns returns all columns present in the dataframe.
	Columns() Columns

	// Distinct returns a new DataFrame that only contains unique rows with
	// respect to the specified columns. If no columns are given Distinct will
	// return rows where allow columns are unique.
	Distinct(columns ...Column) DataFrame

	// Filter returns a new filtered DataFrame according to the filter.
	Filter(filter Filter) DataFrame

	// GroupBy groups rows together for which the values of specified columns
	// are the same.
	GroupBy(columns ...Column) Groups

	// HasColumn reports if a columns with name is present in the dataframe.
	HasColumn(name string) bool

	// AreBools returns true if column by name is of type Bool, otherwise false.
	AreBools(name string) bool
	// AreInts returns true if column by name is of type Int, otherwise false.
	AreInts(name string) bool
	// AreFloats returns true if column by name is of type floats, otherwise
	// false.
	AreFloats(name string) bool
	// AreStrings returns true if column by name is of type String, otherwise
	// false.
	AreStrings(name string) bool

	// Len returns the number of rows in the dataframe.
	Len() int

	// Select returns a new dataframe containing only the specified columns.
	Select(columns ...Column) DataFrame
}

// DataFrames is the slice of DataFrame instances.
type DataFrames []DataFrame
