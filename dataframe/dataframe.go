package dataframe

// DataFrame is an immutable data frame that support filtering, aggregation and
// data manipulation.
//
// Deprecated: This package is deprecated and will be removed in a future.
type DataFrame interface {
	// Column returns a column identified by name, panics if not present.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	Column(name string) Column
	// Columns returns all columns present in the dataframe.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	Columns() Columns

	// Distinct returns a new DataFrame that only contains unique rows with
	// respect to the specified columns. If no columns are given Distinct will
	// return rows where all columns are unique.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	Distinct(columns ...Column) DataFrame

	// Filter returns a new filtered DataFrame according to the filter.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	Filter(filter Filter) DataFrame

	// GroupBy groups rows together for which the values of specified columns
	// are the same.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	GroupBy(columns ...Column) Groups

	// HasColumn reports if a columns with name is present in the dataframe.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	HasColumn(name string) bool

	// AreBools returns true if column by name is of type Bool, otherwise false.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	AreBools(name string) bool
	// AreInts returns true if column by name is of type Int, otherwise false.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	AreInts(name string) bool
	// AreFloats returns true if column by name is of type floats, otherwise
	// false.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	AreFloats(name string) bool
	// AreStrings returns true if column by name is of type String, otherwise
	// false.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	AreStrings(name string) bool

	// Len returns the number of rows in the dataframe.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	Len() int

	// Select returns a new dataframe containing only the specified columns.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	Select(columns ...Column) DataFrame
}

// DataFrames is the slice of DataFrame instances.
//
// Deprecated: This package is deprecated and will be removed in a future.
type DataFrames []DataFrame
