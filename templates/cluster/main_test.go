package main

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestTemplate(t *testing.T) {
	// Read the input from the file.
	input := input{}
	b, err := os.ReadFile("input.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &input); err != nil {
		t.Fatal(err)
	}

	// Declare the options.
	opt := ClusterOptions{}
	opt.Limits.Duration = 5 * time.Second

	// Declare the solver.
	output, err := solver(context.Background(), input, opt)
	if err != nil {
		t.Fatal(err)
	}

	// Get the expected solution.
	want := Output{}
	b, err = os.ReadFile("testdata/output.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &want); err != nil {
		t.Fatal(err)
	}

	// Get the solution.
	got := output.Solutions[0].(Output)
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}

	got = Output{
		Feasible:          got.Feasible,
		UnassignedIndices: got.UnassignedIndices,
	}

	want = Output{
		Feasible:          want.Feasible,
		UnassignedIndices: want.UnassignedIndices,
	}

	// Compare against expected.
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
