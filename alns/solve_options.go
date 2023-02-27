package alns

import "time"

// NewSolveOptions returns a new instance of SolveOptions. The default
// options are:
//   - iterations: 1
//   - maximumDuration: 1 hour
func NewSolveOptions() SolveOptions {
	return SolveOptions{
		Iterations:      1,
		MaximumDuration: time.Hour,
	}
}

type SolveOptions struct {
	Iterations      int           `json:"iterations" default:"1"`
	MaximumDuration time.Duration `json:"maximum_duration" default:"1h"`
}
