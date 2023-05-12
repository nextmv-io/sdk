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
	Version    Version                `json:"version,omitempty"`
	Options    any                    `json:"options,omitempty"`
	Solutions  []any                  `json:"solutions,omitempty"`
	Statistics *statistics.Statistics `json:"statistics,omitempty"`
}

// NewOutput creates a new Output.
func NewOutput(options any, solutions ...any) Output {
	return Output{
		Solutions: solutions,
		Options:   options,
		Version: Version{
			Sdk: sdk.VERSION,
		},
	}
}
