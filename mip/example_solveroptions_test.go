package mip_test

import (
	"fmt"
	"time"

	"github.com/nextmv-io/sdk/mip"
)

func ExampleSolverOptions_default() {
	solverOptions := mip.NewSolverOptions()

	fmt.Println(solverOptions.SolverVerboseLevel())
	fmt.Println(solverOptions.MIPGapAbsolute())
	fmt.Println(solverOptions.MIPGapRelative())
	fmt.Println(solverOptions.MaximumDuration())
	// Output:
	// 0
	// 1e-11
	// 0.05
	// 10m0s

}

func ExampleSolverOptions_change() {
	solverOptions := mip.NewSolverOptions()
	solverOptions.SetVerboseLevel(mip.HIGH)
	solverOptions.SetMIPGapAbsolute(1.23)
	solverOptions.SetMIPGapRelative(0.5)
	solverOptions.SetMaximumDuration(time.Minute)

	fmt.Println(solverOptions.SolverVerboseLevel())
	fmt.Println(solverOptions.MIPGapAbsolute())
	fmt.Println(solverOptions.MIPGapRelative())
	fmt.Println(solverOptions.MaximumDuration())
	// Output:
	// 3
	// 1.23
	// 0.5
	// 1m0s
}
