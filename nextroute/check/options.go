package check

import "time"

// Options are the options for a check.
type Options struct {
	Duration  time.Duration `json:"duration" usage:"maximum duration of the check" default:"30s"`
	Verbosity string        `json:"verbosity" usage:"verbosity of the check, current available options are 'off' , 'low', 'medium', 'high', 'very_high'" default:"off"`
}