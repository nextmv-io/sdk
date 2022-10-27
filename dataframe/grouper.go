package dataframe

type Groups interface {
	Aggregate(aggregations Aggregations) (DataFrame, error)

	DataFrames() ([]DataFrame, error)
}
