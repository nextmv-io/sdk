package mip

import "time"

// Solution contains the results of a Solver.Solve invocation.
type Solution interface {
	// HasValues returns true if the solver was able to associate values with
	// variables.
	HasValues() bool
	// IsInfeasible returns true if the solver has proven that the model
	// defines an infeasible solution, otherwise returns false.
	IsInfeasible() bool
	// IsNumericalFailure returns true if the solver encountered a numerical
	// failure, otherwise returns false. Numerical failures can have different
	// causes and depend on the underlying solver provider.
	IsNumericalFailure() bool
	// IsOptimal returns true if the solver has proven that the solution
	// is one of the optimal solutions, otherwise returns false.
	IsOptimal() bool
	// IsSubOptimal returns true if the solver sub-optimal conform the
	// model of the underlying solver, otherwise false.
	IsSubOptimal() bool
	// IsTimeOut returns true if the solver returned due to a time limit
	// before reaching a conclusion, otherwise returns true.
	IsTimeOut() bool
	// IsUnbounded returns true if the solver proved the solution is
	// unbounded, otherwise returns false. An unbounded solution is
	// a solution that can be improved by changing a variable in a
	// direction it is not limited by bounds.
	IsUnbounded() bool
	// ObjectiveValue return the value of the objective, the value should only
	// be used if HasValues returns true. Returns 0.0 if HasValues is false.
	ObjectiveValue() float64
	// RunTime returns the duration it took for the Solver.Solve to return
	// this solution
	RunTime() time.Duration
	// Value returns the value the solver has associated with the variable
	// in the invoking solution. the value should only be used if HasValues
	// returns true. Returns math.MaxFloat64 if HasValues is false.
	Value(variable Var) float64
}
