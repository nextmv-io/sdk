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
	"github.com/nextmv-io/sdk/types"
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
				MinDuration  types.Duration "json:\"min_duration\" default:\"2h\" usage:\"minimum working time per shift\""
				MaxDuration  types.Duration "json:\"max_duration\" default:\"8h\" usage:\"maximum working time per shift\""
				RecoveryTime types.Duration "json:\"recovery_time\" default:\"8h\" usage:\"minimum time between shifts\""
			}{
				MinDuration:  types.Duration(2 * time.Hour),
				MaxDuration:  types.Duration(8 * time.Hour),
				RecoveryTime: types.Duration(8 * time.Hour),
			},
			Week: struct {
				MaxDuration types.Duration "json:\"max_duration\" default:\"40h\" usage:\"maximum working time per week\""
			}{
				MaxDuration: types.Duration(40 * time.Hour),
			},
			Day: struct {
				MaxDuration types.Duration "json:\"max_duration\" default:\"10h\" usage:\"maximum working time per day\""
			}{
				MaxDuration: types.Duration(10 * time.Hour),
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
