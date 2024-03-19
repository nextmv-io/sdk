package main

import (
	"os"
	"testing"

	"github.com/nextmv-io/sdk/golden"
)

func TestMain(m *testing.M) {
	golden.Setup()
	code := m.Run()
	golden.Teardown()
	os.Exit(code)
}

// TestGolden executes a golden file test, where the .json input is fed and an
// output is expected.
func TestGolden(t *testing.T) {
	golden.FileTests(
		t,
		"input.json",
		golden.Config{
			Args: []string{
				"-duration=1s",
			},
			TransientFields: []golden.TransientField{
				{Key: ".version.sdk", Replacement: golden.StableVersion},
				{Key: ".solutions[0].statistics.time.elapsed", Replacement: golden.StableDuration},
				{Key: ".solutions[0].statistics.time.elapsed_seconds", Replacement: golden.StableFloat},
				{Key: ".solutions[0].statistics.time.start", Replacement: golden.StableTime},
			},
		},
	)
}

// TestGoldenBash executes a golden file test, where the bash file is run and
// the output is compared against the expected one.
func TestGoldenBash(t *testing.T) {
	// Execute the rest of the bash commands.
	golden.BashTest(t, "./bash", golden.BashConfig{
		DisplayStdout: true,
		DisplayStderr: true,
	})
}
