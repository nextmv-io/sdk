package run_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/nextmv-io/sdk/run"
)

//nolint:lll
type ParallelSolveOptions struct {
	Iterations           int           `json:"iterations"  usage:"maximum number of iterations, -1 assumes no limit; iterations are counted after start solutions are generated" default:"-1"`
	Duration             time.Duration `json:"duration" usage:"maximum duration of the solver" default:"5s"`
	ParallelRuns         int           `json:"parallel_runs" usage:"maximum number of parallel runs, -1 results in using all available resources" default:"-1"`
	StartSolutions       int           `json:"start_solutions" usage:"number of solutions to generate on top of those passed in; one solution generated with sweep algorithm, the rest generated randomly" default:"-1"`
	RunDeterministically bool          `json:"run_deterministically"  usage:"run the parallel solver deterministically"`
}

type SampleOption struct {
	Solve  ParallelSolveOptions `json:"solve,omitempty"`
	Custom struct {
		Solve struct {
			Plateau struct {
				Auto bool `json:"auto" usage:"whether to enable auto plateau detection"`
			} `json:"plateau"`
		} `json:"solve"`
	} `json:"custom,omitempty"`
}

func Test_FlagParser(t *testing.T) {
	// Simulate command line arguments
	origArgs := os.Args
	os.Args = []string{
		"cmd",
		"-solve.iterations=10",
		"-custom.solve.plateau.auto=true",
		"-solve.duration=5m",
	}
	defer func() { os.Args = origArgs }()

	_, option, err := run.FlagParser[SampleOption, run.CLIRunnerConfig]()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedOption := SampleOption{
		Solve: ParallelSolveOptions{
			Iterations:           10,
			Duration:             5 * time.Minute,
			ParallelRuns:         -1,
			StartSolutions:       -1,
			RunDeterministically: false,
		},
		Custom: struct {
			Solve struct {
				Plateau struct {
					Auto bool `json:"auto" usage:"whether to enable auto plateau detection"`
				} `json:"plateau"`
			} `json:"solve"`
		}{
			Solve: struct {
				Plateau struct {
					Auto bool `json:"auto" usage:"whether to enable auto plateau detection"`
				} `json:"plateau"`
			}{
				Plateau: struct {
					Auto bool `json:"auto" usage:"whether to enable auto plateau detection"`
				}{
					Auto: true,
				},
			},
		},
	}
	if !reflect.DeepEqual(option, expectedOption) {
		t.Errorf("expected option %+v, got %+v", expectedOption, option)
	}
}
