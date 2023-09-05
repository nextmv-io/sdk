package main

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/nextmv-io/sdk/mip"
	"github.com/nextmv-io/sdk/run/schema"
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

	// Set default options.
	options := options{
		OverSupplyPenalty:  1000,
		UnderSupplyPenalty: 500,
		MaxHoursPerDay:     10 * time.Hour,
		MaxHoursPerWeek:    40 * time.Hour,
		MinHoursPerShift:   2 * time.Hour,
		MaxHoursPerShift:   8 * time.Hour,
		HoursBetweenShifts: 8 * time.Hour,
		Limits:             mip.Limits{Duration: 4 * time.Minute},
	}
	output, err := solver(context.Background(), input, options)
	if err != nil {
		t.Fatal(err)
	}

	expectedOutput := schema.Output{}
	b, err = os.ReadFile("testdata/output.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &expectedOutput); err != nil {
		t.Fatal(err)
	}

	got := output.Statistics.Result.Value
	want := expectedOutput.Statistics.Result.Value

	// Compare against expected.
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
