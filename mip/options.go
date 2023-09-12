package mip

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// SolveOptions are options that can be cionfigured for any solver.
type SolveOptions struct {
	// Duration is the maximum duration of the solver. A duration limit of 0 is
	// treated as infinity.
	Duration time.Duration `json:"duration" usage:"Maximum duration of the solver." default:"30s"`
	// Verbosity of the solver in the console.
	Verbosity Verbosity `json:"verbosity" usage:"{off, low, medium, high} Verbosity of the solver in the console." default:"off"`
	// Mip-specific options.
	Mip MipOptions `json:"mip" usage:"Options specific to MIP problems. Linear problems do not use these options."`
	// Control options for the specific solver.
	Control ControlOptions `json:"control" usage:"Options to control a specific solver, as defined by the provider."`
}

// MipOptions are options specific to MIP problems. LP problems do not use
// these options.
type MipOptions struct {
	// Gap stopping criteria.
	Gap GapOptions `json:"gap" usage:"Gap stopping criteria."`
}

// GapOptions specifies the gap stopping criteria.
type GapOptions struct {
	// Absolute gap.
	Absolute float64 `json:"absolute" usage:"Absolute gap stopping value. If the problem is an integer problem the solver will stop if the gap between the relaxed problem and the best found integer problem is less than this value." default:"0.000001"`
	// Relative gap.
	Relative float64 `json:"relative" usage:"Relative gap stopping value. If the problem is an integer problem the solver will stop if the relative gap between the relaxed problem and the best found integer problem is less than this value." default:"0.0001"`
}

// Verbosity specifies the level of verbosity of the back-end solver.
type Verbosity string

const (
	// Off logs nothing.
	Off Verbosity = "off"
	// Low logs essentials, depends on the back-end solver.
	Low Verbosity = "low"
	// Medium logs essentials plus high level events,
	// depends on the back-end solver.
	Medium Verbosity = "medium"
	// High logs everything the underlying logs,
	// depends on the back-end solver.
	High Verbosity = "high"
)

// ControlOptions allow the user to define solver-specific parameters. The
// parameters' names and types must be known beforehand so that the correct
// type is used.
type ControlOptions struct {
	Bool   string `json:"bool" usage:"List of solver-specific control options (configurations) with bool values. Example: \"name1=value1,name2=value2\", where value1 and value2 are bool values."`
	Float  string `json:"float" usage:"List of solver-specific control options (configurations) with float values. Example: \"name1=value1,name2=value2\", where value1 and value2 are float values."`
	Int    string `json:"int" usage:"List of solver-specific control options (configurations) with int values. Example: \"name1=value1,name2=value2\", where value1 and value2 are int values."`
	String string `json:"string" usage:"List of solver-specific control options (configurations) with string values. Example: \"name1=value1,name2=value2\", where value1 and value2 are string values."`
}

// MarshalJSON implements the [json.Marshaler] interface.
func (controlOptions ControlOptions) MarshalJSON() ([]byte, error) {
	v, err := controlOptions.ToTyped()
	if err != nil {
		return nil, err
	}

	return json.Marshal(v)
}

// TypedControlOptions is the typed equivalent to [ControlOptions] for
// configuring a solver's parameters.
type TypedControlOptions struct {
	Bool   []TypedControlOption[bool]    `json:"bool"`
	Float  []TypedControlOption[float64] `json:"float"`
	Int    []TypedControlOption[int]     `json:"int"`
	String []TypedControlOption[string]  `json:"string"`
}

// TypedControlOption defines a generic way to specify a control parameter for a
// solver.
type TypedControlOption[T string | float64 | int | bool] struct {
	// Name of the option.
	Name string `json:"name"`
	// Value for the option. The value's type is defined by how the control
	// option is instantiated.
	Value T `json:"value"`
}

// ToTyped converts the string-based control options into fully typed options.
func (controlOptions ControlOptions) ToTyped() (*TypedControlOptions, error) {
	typedControlOptions := TypedControlOptions{
		Bool:   make([]TypedControlOption[bool], 0),
		Float:  make([]TypedControlOption[float64], 0),
		Int:    make([]TypedControlOption[int], 0),
		String: make([]TypedControlOption[string], 0),
	}

	if controlOptions.Bool != "" {
		untypedOptions := strings.Split(controlOptions.Bool, ",")
		typedOptions := make([]TypedControlOption[bool], len(untypedOptions))
		for ix, untypedOption := range untypedOptions {
			name, value, err := extractNameValue(untypedOption)
			if err != nil {
				return nil, err
			}

			typedValue, err := strconv.ParseBool(value)
			if err != nil {
				return nil, fmt.Errorf("option %s with non-valid bool value %v: %w", name, value, err)
			}

			typedOptions[ix] = TypedControlOption[bool]{
				Name:  name,
				Value: typedValue,
			}
		}

		typedControlOptions.Bool = typedOptions
	}

	if controlOptions.Float != "" {
		untypedOptions := strings.Split(controlOptions.Float, ",")
		typedOptions := make([]TypedControlOption[float64], len(untypedOptions))
		for ix, untypedOption := range untypedOptions {
			name, value, err := extractNameValue(untypedOption)
			if err != nil {
				return nil, err
			}

			typedValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, fmt.Errorf("option %s with non-valid float64 value %v: %w", name, value, err)
			}

			typedOptions[ix] = TypedControlOption[float64]{
				Name:  name,
				Value: typedValue,
			}
		}

		typedControlOptions.Float = typedOptions
	}

	if controlOptions.Int != "" {
		untypedOptions := strings.Split(controlOptions.Int, ",")
		typedOptions := make([]TypedControlOption[int], len(untypedOptions))
		for ix, untypedOption := range untypedOptions {
			name, value, err := extractNameValue(untypedOption)
			if err != nil {
				return nil, err
			}

			typedValue, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("option %s with non-valid int value %v: %w", name, value, err)
			}

			typedOptions[ix] = TypedControlOption[int]{
				Name:  name,
				Value: typedValue,
			}
		}

		typedControlOptions.Int = typedOptions
	}

	if controlOptions.String != "" {
		untypedOptions := strings.Split(controlOptions.String, ",")
		typedOptions := make([]TypedControlOption[string], len(untypedOptions))
		for ix, untypedOption := range untypedOptions {
			name, value, err := extractNameValue(untypedOption)
			if err != nil {
				return nil, err
			}

			typedOptions[ix] = TypedControlOption[string]{
				Name:  name,
				Value: value,
			}
		}

		typedControlOptions.String = typedOptions
	}

	return &typedControlOptions, nil
}

func extractNameValue(option string) (name string, value string, err error) {
	splitOption := strings.Split(option, "=")
	if len(splitOption) != 2 {
		return name, value,
			fmt.Errorf("option %s with unexpected format, want \"name=value\"", option)
	}

	name, value = splitOption[0], splitOption[1]

	return
}
