// Package statistics provides the schema for the statistics section of the
// output.
package statistics

// NewStatistics creates new statistics.
func NewStatistics() *Statistics {
	return &Statistics{
		Schema: "v1",
	}
}

// Statistics is the structure of the statistics section of the output.
type Statistics struct {
	Schema     string      `json:"schema,omitempty"`
	Run        *Run        `json:"run,omitempty"`
	Result     *Result     `json:"result,omitempty"`
	SeriesData *SeriesData `json:"series_data,omitempty"`
}

// Run is the structure of the run section of the statistics.
type Run struct {
	Duration *float64 `json:"duration,omitempty"`
	Custom   any      `json:"custom,omitempty"`
}

// Result is the structure of the result section of the statistics.
type Result struct {
	Duration   *float64 `json:"duration,omitempty"`
	Iterations *int     `json:"iterations,omitempty"`
	Value      *Float64 `json:"value,omitempty"`
	Custom     any      `json:"custom,omitempty"`
}

// SeriesData is the structure of the series section of the statistics.
type SeriesData struct {
	Value  Series   `json:"value,omitempty"`
	Custom []Series `json:"custom,omitempty"`
}

// Series is the structure of a time series.
type Series struct {
	Name       string      `json:"name,omitempty"`
	DataPoints []DataPoint `json:"data_points,omitempty"`
}

// DataPoint is the structure of a time series datum.
type DataPoint struct {
	X Float64 `json:"x"`
	Y Float64 `json:"y"`
}
