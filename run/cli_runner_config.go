package run

import "errors"

// CPUProfiler is the interface for profiling CPU usage.
type CPUProfiler interface {
	CPUProfilePath() string
}

// MemoryProfiler is the interface for profiling memory usage.
type MemoryProfiler interface {
	MemoryProfilePath() string
}

// CliRunnerConfig is the configuration of the  CliRunner.
type CliRunnerConfig struct {
	Runner struct {
		Input struct {
			Path string // Path to input file
		}
		Profile struct {
			CPU    string // CPU profile location
			Memory string // Memory profile location
		}
		Output struct {
			Path      string // Path to output file
			Solutions string // All solutions or last (best) solution
			Quiet     bool   // Only output solutions
		}
	}
}

// CPUProfilePath returns the CPU profile path.
func (c CliRunnerConfig) CPUProfilePath() string {
	return c.Runner.Profile.CPU
}

// MemoryProfilePath returns the memory profile path.
func (c CliRunnerConfig) MemoryProfilePath() string {
	return c.Runner.Profile.Memory
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
	case "":
		return All, nil
	case "all":
		return All, nil
	case "last":
		return Last, nil
	default:
		return Last, errors.New(`solutions must be "all" or "last"`)
	}
}
