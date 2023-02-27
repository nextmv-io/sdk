package nextroute

import (
	"github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/connect"
)

// SolveInformation holds the information about the current solve run.
type SolveInformation = alns.SolveInformation[Solution]

// SolveObserver observes a solve run.
type SolveObserver = alns.SolveObserver[Solution]

// SolveOperator is a solve-operator.
type SolveOperator = alns.SolveOperator[Solution]

// SolveOperatorImpl is the implementation of a solve-operator implementation.
type SolveOperatorImpl = alns.SolveOperatorImpl[Solution]

// SolveParameter is a solve-parameter.
type SolveParameter = alns.SolveParameter[Solution]

// SolveParameters is a list of solve-parameters.
type SolveParameters = alns.SolveParameters[Solution]

var (
	solverConnect = connect.NewConnector("sdk", "NextRouteSolver")

	newSolveOperatorUnPlan func(SolveParameter) SolveOperatorUnPlan

	newSolveOperatorRestart func(SolveParameter) SolveOperatorRestart

	newSolveOperatorPlan func(SolveParameter) SolveOperatorPlan
)
