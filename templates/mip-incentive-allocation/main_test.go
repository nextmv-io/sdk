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
	input := incentiveAllocationProblem{}
	b, err := os.ReadFile("input.json")
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(b, &input); err != nil {
		t.Fatal(err)
	}

	// Declare the output.
	options := Option{}
	options.Limits.Duration = 5 * time.Second
	output, err := solver(context.Background(), input, options)
	if err != nil {
		t.Fatal(err)
	}

	got := output.Solutions[0].(Output)
	if err = json.Unmarshal(b, &got); err != nil {
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

	// Compare against expected.
	if !reflect.DeepEqual(got.Value, want.Value) {
		t.Errorf("got %+v, want %+v", got.Value, want.Value)
	}
}
