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

// TestGoldenBash executes a golden file test, where the bash file is run and
// the output is compared against the expected one.
func TestGoldenBash(t *testing.T) {
	// Execute the rest of the bash commands.
	golden.BashTest(t, "./bash", golden.BashConfig{
		DisplayStdout: true,
		DisplayStderr: true,
	})
}
