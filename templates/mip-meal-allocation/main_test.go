package main

import (
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
	if err = json.Unmarshal(b, &input); err != nil {
		t.Fatal(err)
	}

	// Declare the output.
	output, err := solver(input, Option{Duration: 5 * time.Second})
	if err != nil {
		t.Fatal(err)
	}

	got := output[0]
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
	if !reflect.DeepEqual(got.Binkies, want.Binkies) {
		t.Errorf("got %+v, want %+v", got.Binkies, want.Binkies)
	}
}
