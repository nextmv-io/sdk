package mip_test

import (
	"fmt"
	"time"

	"github.com/nextmv-io/sdk/mip"
)

func ExampleSolveOptions_default() {
	solveOptions := mip.NewSolveOptions()

	fmt.Println(solveOptions.Verbosity())
	fmt.Println(solveOptions.MIPGapAbsolute())
	fmt.Println(solveOptions.MIPGapRelative())
	fmt.Println(solveOptions.MaximumDuration())
	// Output:
	// 0
	// 1e-11
	// 0.05
	// 10m0s
}

func ExampleSolveOptions_change() {
	solveOptions := mip.NewSolveOptions()
	solveOptions.SetVerbosity(mip.HIGH)
	err := solveOptions.SetMIPGapAbsolute(1.23)
	if err != nil {
		panic(err)
	}
	err = solveOptions.SetMIPGapRelative(0.5)
	if err != nil {
		panic(err)
	}
	err = solveOptions.SetMaximumDuration(time.Minute)
	if err != nil {
		panic(err)
	}

	fmt.Println(solveOptions.Verbosity())
	fmt.Println(solveOptions.MIPGapAbsolute())
	fmt.Println(solveOptions.MIPGapRelative())
	fmt.Println(solveOptions.MaximumDuration())
	// Output:
	// 3
	// 1.23
	// 0.5
	// 1m0s
}
