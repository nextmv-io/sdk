package dataframe

type DataFrame interface {
	Column(string) (AnyColumn, error)
	Columns() AnyColumns

	Bools(string) (BoolColumn, error)
	Ints(string) (IntColumn, error)
	Floats(string) (FloatColumn, error)
	Strings(string) (StringColumn, error)

	// Distinct(columns ...Column) DataFrame

	Filter(filter Filter) (DataFrame, error)

	GroupBy(columns Columns) Groups

	Len() int

	Select(columns ...Column) DataFrame
}
