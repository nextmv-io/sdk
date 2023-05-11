// Package schema provides the schema for the output.
package schema

import (
	"github.com/nextmv-io/sdk"
	"github.com/nextmv-io/sdk/run/statistics"
)

// Version of the sdk.
type Version struct {
	Sdk string `json:"sdk"`
}

// Output adds Output information by wrapping the solutions.
type Output struct {
	Statistics statistics.Statistics `json:"statistics,omitempty"`
	Options    any                   `json:"options,omitempty"`
	Solution   any                   `json:"solution"`
	Version    Version               `json:"version,omitempty"`
}

// NewOutput creates a new Output.
func NewOutput(solution, options any) Output {
	return Output{
		Solution: solution,
		Options:  options,
		Version: Version{
			Sdk: sdk.VERSION,
		},
	}
}
