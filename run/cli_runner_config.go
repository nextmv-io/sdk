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

// OutputPather is the interface for getting the output path.
type OutputPather interface {
	OutputPath() string
}

// Quieter is the interface for getting the quiet flag.
type Quieter interface {
	Quiet() bool
}

// SolutionLimiter is the interface for getting the configured solutions.
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
			Solutions string `default:"all" usage:"Return all or last solution"`
			Quiet     bool   `default:"false" usage:"Do not return statistics"`
		}
	}
}

// OutputPath returns the output path.
func (c CLIRunnerConfig) OutputPath() string {
	return c.Runner.Output.Path
}

// Quiet returns the quiet flag.
func (c CLIRunnerConfig) Quiet() bool {
	return c.Runner.Output.Quiet
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
		return All, nil
	case "all":
		return All, nil
	case "last":
		return Last, nil
	default:
		return Last, errors.New(`solutions must be "all" or "last"`)
	}
}
