package mip

import "time"

// Limits on the MIP solver, such as how long it can run for.
type Limits struct {
	Duration time.Duration `json:"duration" usage:"maximum duration of the solver. A duration limit of 0 is treated as infinity" default:"30s"`
}
