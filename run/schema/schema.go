// Package schema provides the schema for the output.
package schema

import (
	"runtime/debug"
	"strings"

	"github.com/nextmv-io/sdk/run/statistics"
)

// Version info of the output.
type Version map[string]string

// Output adds Output information by wrapping the solutions.
type Output struct {
	Version    Version                `json:"version,omitempty"`
	Options    any                    `json:"options,omitempty"`
	Solutions  []any                  `json:"solutions,omitempty"`
	Statistics *statistics.Statistics `json:"statistics,omitempty"`
}

// NewOutput creates a new Output.
func NewOutput[Solution any](options any, solutions ...Solution) Output {
	// convert solutions to any
	solutionsAny := make([]any, len(solutions))
	for i, solution := range solutions {
		solutionsAny[i] = solution
	}
	// collect known dependency versions
	dependencies := collectKnownDependencies()
	return Output{
		Solutions: solutionsAny,
		Options:   options,
		Version:   dependencies,
	}
}

// knownDependencies is a list of known dependencies that we want to put in the
// version of the output.
var knownDependencies = []struct {
	name string
	path string
}{
	{name: "sdk", path: "github.com/nextmv-io/sdk"},
	{name: "nextroute", path: "github.com/nextmv-io/nextroute"},
	{name: "go-mip", path: "github.com/nextmv-io/go-mip"},
	{name: "go-highs", path: "github.com/nextmv-io/go-highs"},
	{name: "go-xpress", path: "github.com/nextmv-io/go-xpress"},
}

func collectKnownDependencies() Version {
	// We use the debug.ReadBuildInfo to get the version of the dependencies.
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		// If this happens, we're not running in a module context. We cannot
		// provide the version of the dependencies.
		return map[string]string{}
	}

	// Search all dependencies for known ones and collect their versions.
	deps := map[string]string{}
	for _, dep := range bi.Deps {
		for _, knownDep := range knownDependencies {
			if strings.HasPrefix(dep.Path, knownDep.path) {
				deps[knownDep.name] = dep.Version
			}
		}
	}
	return deps
}
