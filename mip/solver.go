package mip

import "time"

// Solver exposes an API to run a MIP solver
//
//			model := mip.NewModel()
//
//	 	// build the model
//
//			provider := "my_favorite_solver"
//
//			solver, err := NewSolver(provider, model)
//
//			solution, err := solver.Solve(mip.DefaultSolverOptions())
type Solver interface {
	// Solve is the entrypoint to solve the model associated with
	// the invoking solver. Returns a solution when the invoking solver
	// reaches a conclusion.
	Solve(options SolveOptions) (Solution, error)
}

// SolveOptions interface to options for back-end solver.
type SolveOptions interface {
	// FloatParameters returns all float parameter settings.
	FloatParameters() FloatSolverParameterSettings
	// GetFloatParameter returns value set for parameter, returns error if no
	// such parameter has been set.
	GetFloatParameter(parameter SolverParameter) (float64, error)
	// GetIntParameter returns value set for parameter, returns error if no
	// such parameter has been set.
	GetIntParameter(parameter SolverParameter) (int64, error)
	// GetStringParameter returns value set for parameter, returns error if no
	// such parameter has been set.
	GetStringParameter(parameter SolverParameter) (string, error)
	// IntParameters returns all int parameter settings.
	IntParameters() IntSolverParameterSettings
	// MaximumDuration returns maximum duration of a Solver.Solve invocation.
	MaximumDuration() time.Duration
	// MIPGapAbsolute returns the absolute gap stopping value. If the problem
	// is an integer problem the solver will stop if the gap between the relaxed
	// problem and the best found integer problem is less than this value.
	MIPGapAbsolute() float64
	// MIPGapRelative returns the relative gap stopping value. If the problem
	// is an integer problem the solver will stop if the relative gap between
	// the relaxed problem and the best found integer problem is less than
	// this value.
	MIPGapRelative() float64
	// SetFloatParameter specifies the value to use for parameter, this is
	// back-end-solver specific.
	SetFloatParameter(parameter SolverParameter, value float64)
	// SetIntParameter specifies the value to use for parameter, this is
	// back-end-solver specific.
	SetIntParameter(parameter SolverParameter, value int64)
	// SetMaximumDuration specifies the maximum duration of a Solver.Solve
	// invocation.
	SetMaximumDuration(duration time.Duration) error
	// SetMIPGapAbsolute specifies the absolute gap stopping value, only
	// used in case the problem is an integer solution, raises an error if
	// value is not strictly positive (> 0).
	SetMIPGapAbsolute(value float64) error
	// SetMIPGapRelative specifies the relative gap stopping value, only
	// used in case the problem is an integer solution, raises an error if
	// value is less than zero or larger or equal to one.
	SetMIPGapRelative(value float64) error
	// SetStringParameter specifies the value to use for parameter, this is
	// back-end-solver specific.
	SetStringParameter(parameter SolverParameter, value string)
	// SetVerbosity specifies the verbosity level of the underlying
	// back-end solver. Forwards output to std out.
	SetVerbosity(verbosity Verbosity)
	// StringParameters returns all string parameter settings
	StringParameters() StringSolverParameterSettings
	// Verbosity returns the configured verbosity level of the
	// underlying back-end solver.
	Verbosity() Verbosity
}

// Verbosity specifies the level of verbosity of the back-end solver.
type Verbosity int

const (
	// OFF logs nothing.
	OFF Verbosity = iota
	// LOW logs essentials, depends on the back-end solver.
	LOW
	// MEDIUM logs essentials plus high level events,
	// depends on the back-end solver.
	MEDIUM
	// HIGH logs everything the underlying logs,
	// depends on the back-end solver.
	HIGH
)

// SolverParameter identifier for parameters in the back-end solver.
type SolverParameter int

// SolverProvider identifier for a back-end solver.
type SolverProvider string

// GetSolverParameter interface to retrieve a solver parameter.
type GetSolverParameter interface {
	SolverParameter() SolverParameter
}

// FloatSolverParameterSetting interface for setting of type float64.
type FloatSolverParameterSetting interface {
	GetSolverParameter
	Value() float64
}

// FloatSolverParameterSettings slice of FloatSolverParameterSetting.
type FloatSolverParameterSettings []FloatSolverParameterSetting

// IntSolverParameterSetting interface for setting of type int64.
type IntSolverParameterSetting interface {
	GetSolverParameter
	Value() int64
}

// IntSolverParameterSettings slice of IntSolverParameterSetting.
type IntSolverParameterSettings []IntSolverParameterSetting

// StringSolverParameterSetting interface for setting of type string.
type StringSolverParameterSetting interface {
	GetSolverParameter
	Value() string
}

// StringSolverParameterSettings slice of StringSolverParameterSetting.
type StringSolverParameterSettings []StringSolverParameterSetting
