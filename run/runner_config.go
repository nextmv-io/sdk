package run

import "errors"

// RunnerConfig is the configuration of a runner.
type RunnerConfig struct {
	Runner struct {
		Input struct {
			Path string // Path to input file
		}
		HTTP struct {
			Address     string // HTTP address (e.g. ":8000")
			Certificate string // TLS certificate file
			Key         string // TLS key file
		}
		Profile struct { // CLI only
			CPU    string // CPU profile location
			Memory string // Memory profile location
		}
		Output struct {
			Path      string    // Path to output file
			Quiet     bool      // Only output solutions
			Stream    bool      // Stream solutions (CLI only)
			Solutions Solutions // All solutions or last (best) solution
		}
	}
}

// Solutions can be all or last.
type Solutions int

// Constants for setting all or last solution output.
const (
	All Solutions = iota
	Last
)

func (s Solutions) String() string {
	if s == Last {
		return "last"
	}
	return "all"
}

// ParseSolutions converts "all" to All and "last" to Last.
func ParseSolutions(s string) (Solutions, error) {
	switch s {
	case "all":
		return All, nil
	case "last":
		return Last, nil
	default:
		return Last, errors.New(`solutions must be "all" or "last"`)
	}
}
