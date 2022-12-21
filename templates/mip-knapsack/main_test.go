package main

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"
)

type output struct {
	Status  string  `json:"status,omitempty"`
	Runtime string  `json:"runtime,omitempty"`
	Value   float64 `json:"value,omitempty"`
	Items   []item  `json:"items,omitempty"`
}

func TestTemplate(t *testing.T) {
	// Read the input from the file.
	input := input{}
	b, err := os.ReadFile("input.json")
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(b, &input); err != nil {
		t.Fatal(err)
	}

	// Declare the solution.
	solution, err := solver(input, Option{Duration: 5 * time.Second})
	if err != nil {
		t.Fatal(err)
	}

	b, err = json.MarshalIndent(solution, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	got := output{}
	if err = json.Unmarshal(b, &got); err != nil {
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

	// Compare against expected.
	if !reflect.DeepEqual(got.Value, want.Value) {
		t.Errorf("got %+v, want %+v", got.Value, want.Value)
	}
}
