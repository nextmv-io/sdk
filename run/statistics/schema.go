// Package statistics provides the schema for the statistics section of the
// output.
package statistics

// Statistics is the structure of the statistics section of the output.
type Statistics struct {
	Result     *Result     `json:"result,omitempty"`
	Run        *Run        `json:"run,omitempty"`
	SeriesData *SeriesData `json:"series_data,omitempty"`
	Schema     *string     `json:"schema,omitempty"`
}

// Run is the structure of the run section of the statistics.
type Run struct {
	Duration *float64 `json:"duration,omitempty"`
	Custom   any      `json:"custom,omitempty"`
}

// Result is the structure of the result section of the statistics.
type Result struct {
	Duration *float64 `json:"duration,omitempty"`
	Value    *Float64 `json:"value,omitempty"`
	Custom   any      `json:"custom,omitempty"`
}

// SeriesData is the structure of the series section of the statistics.
type SeriesData struct {
	Custom []Series `json:"custom,omitempty"`
	Value  Series   `json:"value,omitempty"`
}

// Series is the structure of a time series.
type Series struct {
	Name       string      `json:"name,omitempty"`
	DataPoints []DataPoint `json:"data_points,omitempty"`
}

// DataPoint is the structure of a time series datum.
type DataPoint struct {
	Y Float64 `json:"y"`
	X Float64 `json:"x"`
}
