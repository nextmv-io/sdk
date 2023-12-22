package alns

import (
	sdkAlns "github.com/nextmv-io/sdk/alns"
	"github.com/nextmv-io/sdk/nextroute"
)

// SolveParameter is an alias for sdkAlns.SolveParameter[nextroute.Solution].
type SolveParameter = sdkAlns.SolveParameter[nextroute.Solution]

// SolveInformation is an alias for sdkAlns.SolveInformation[nextroute.Solution].
type SolveInformation = sdkAlns.SolveInformation[nextroute.Solution]

// SolveOperator is an alias for sdkAlns.SolveOperator[nextroute.Solution].
type SolveOperator = sdkAlns.SolveOperator[nextroute.Solution]
