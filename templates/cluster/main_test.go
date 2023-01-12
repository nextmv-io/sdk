package main

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/nextmv-io/sdk/store"
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
	opt := store.DefaultOptions()
	opt.Limits.Duration = 5 * time.Second
	opt.Diagram.Expansion.Limit = 1
	opt.Limits.Solutions = 1

	// Declare the solver.
	solver, err := solver(input, opt)
	if err != nil {
		t.Fatal(err)
	}

	// Get the solution.
	last := solver.Last(context.Background())
	b, err = json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	got := output{}
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}

	// Get the expected solution.
	want := output{}
	b, err = os.ReadFile("testdata/output.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &want); err != nil {
		t.Fatal(err)
	}

	got = output{
		Feasible:          got.Feasible,
		UnassignedIndices: got.UnassignedIndices,
	}

	want = output{
		Feasible:          want.Feasible,
		UnassignedIndices: want.UnassignedIndices,
	}

	// Compare against expected.
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}