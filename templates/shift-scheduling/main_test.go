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
		Penalty: penalty{
			OverSupply:  1000,
			UnderSupply: 500,
		},
		Solve: mip.SolveOptions{
			Duration: 30 * time.Second,
		},
		Limits: limits{
			Shift: struct {
				MinDuration  time.Duration "json:\"min_duration\" default:\"2h\" usage:\"minimum working time per shift\""
				MaxDuration  time.Duration "json:\"max_duration\" default:\"8h\" usage:\"maximum working time per shift\""
				RecoveryTime time.Duration "json:\"recovery_time\" default:\"8h\" usage:\"minimum time between shifts\""
			}{
				MinDuration:  2 * time.Hour,
				MaxDuration:  8 * time.Hour,
				RecoveryTime: 8 * time.Hour,
			},
			Week: struct {
				MaxDuration time.Duration "json:\"max_duration\" default:\"40h\" usage:\"maximum working time per week\""
			}{
				MaxDuration: 40 * time.Hour,
			},
			Day: struct {
				MaxDuration time.Duration "json:\"max_duration\" default:\"10h\" usage:\"maximum working time per day\""
			}{
				MaxDuration: 10 * time.Hour,
			},
		},
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
