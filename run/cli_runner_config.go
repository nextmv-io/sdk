package run

import "errors"

// CPUProfiler is the interface a runner configuration can implement to return
// the CPU profile path.
type CPUProfiler interface {
	CPUProfilePath() string
}

// MemoryProfiler is the interface a runner configuration can implement to
// return the memory profile path.
type MemoryProfiler interface {
	MemoryProfilePath() string
}

// OutputPather is the interface a runner configuration can implement to return
// the output path.
type OutputPather interface {
	OutputPath() string
}

// SolutionLimiter is the interface a runner configuration can implement to
// control whether all or only the last solution is returned.
type SolutionLimiter interface {
	Solutions() (Solutions, error)
}

// CLIRunnerConfig is the configuration of the  CliRunner.
type CLIRunnerConfig struct {
	Runner struct {
		Input struct {
			Path string `usage:"The input file path"`
		}
		Profile struct {
			CPU    string `usage:"The CPU profile file path"`
			Memory string `usage:"The memory profile file path"`
		}
		Output struct {
			Path      string `usage:"The output file path"`
			Solutions string `default:"last" usage:"{all, last}"`
		}
	}
}

// OutputPath returns the output path.
func (c CLIRunnerConfig) OutputPath() string {
	return c.Runner.Output.Path
}

// CPUProfilePath returns the CPU profile path.
func (c CLIRunnerConfig) CPUProfilePath() string {
	return c.Runner.Profile.CPU
}

// MemoryProfilePath returns the memory profile path.
func (c CLIRunnerConfig) MemoryProfilePath() string {
	return c.Runner.Profile.Memory
}

// Solutions returns the configured solutions.
func (c CLIRunnerConfig) Solutions() (Solutions, error) {
	return ParseSolutions(c.Runner.Output.Solutions)
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
		return Last, nil
	case "all":
		return All, nil
	case "last":
		return Last, nil
	default:
		return Last, errors.New(`solutions must be "all" or "last"`)
	}
}
