package golden

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// BashTest executes a golden file test for a bash command. It walks over
// the goldenDir to gather all .sh scripts present in the dir. It then executes
// each of the scripts and compares expected vs. actual outputs. If
// displayStdout or displayStderr are true, the output of each script will be
// composed of the resulting stderr + stdout.
func BashTest(
	t *testing.T,
	goldenDir string,
	bashConfig BashConfig,
) {
	// Fail immediately, if dir does not exist
	if stat, err := os.Stat(goldenDir); err != nil || !stat.IsDir() {
		t.Fatalf("dir %s does not exist", goldenDir)
	}

	// Collect bash scripts.
	var scripts []string
	fn := func(path string, _ os.FileInfo, _ error) error {
		// Only consider .sh files
		if strings.HasSuffix(path, ".sh") {
			scripts = append(scripts, path)
		}

		return nil
	}
	if err := filepath.Walk(goldenDir, fn); err != nil {
		t.Fatal("error walking over files: ", err)
	}

	// Execute a golden file test for each script.
	for _, script := range scripts {
		BashTestFile(t, script, bashConfig)
	}

	// Post-process files containing volatile data.
	postProcessVolatileData(t, bashConfig)
}

// BashTestFile executes a golden file test for a single bash script. The
// script is executed and the expected output is compared with the actual
// output.
func BashTestFile(
	t *testing.T,
	script string,
	bashConfig BashConfig,
) {
	goldenFilePath := script + goldenExtension
	// Function run by the test.
	f := func(t *testing.T) {
		// Execute a bash command which consists of executing a .sh file.
		cmd := exec.Command("bash", script)

		// Add additional command flags.
		cmd.Args = append(cmd.Args, bashConfig.AdditionalCommandFlags...)

		// Pass environment and add custom environment variables
		cmd.Env = os.Environ()
		for _, e := range bashConfig.Envs {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", e[0], e[1]))
		}

		// Run the command and gather the output bytes.
		out, err := runCmd(cmd, bashConfig.DisplayStdout, bashConfig.DisplayStderr)
		if err != nil {
			t.Fatal(err)
		}

		got := string(out)
		if !bashConfig.OutputProcessConfig.KeepVolatileData {
			// Replace default volatile content with a placeholder
			got = regexReplaceAllDefault(got)
		}
		// Apply custom volatile regex replacements
		for _, r := range bashConfig.OutputProcessConfig.VolatileRegexReplacements {
			got = regexReplaceCustom(got, r.Replacement, r.Regex)
		}
		// Write the output bytes to a .golden file, if the test is being
		// updated
		if *update || bashConfig.OutputProcessConfig.AlwaysUpdate {
			if err := os.WriteFile(goldenFilePath, []byte(got), 0o644); err != nil {
				t.Fatal("error writing bash output to file: ", err)
			}
		}

		// Read the .golden file.
		outGolden, err := os.ReadFile(goldenFilePath)
		if err != nil {
			t.Fatal("error reading file: ", goldenFilePath, ": ", err)
		}

		// Perform the golden file comparison.
		expected := string(outGolden)
		if got != expected {
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(got, expected, true)
			t.Errorf(
				"\ngot:\n%s\nexpected:\n%s\ndiffs (look for the colors):\n%s",
				got,
				expected,
				dmp.DiffPrettyText(diffs),
			)
		}
	}

	// Test is executed.
	t.Run(script, f)

	// Run post-process functions.
	for _, f := range bashConfig.PostProcessFunctions {
		err := f(goldenFilePath)
		if err != nil {
			t.Fatalf("error running post-process function: %v", err)
		}
	}
}

func postProcessVolatileData(
	t *testing.T,
	bashConfig BashConfig,
) {
	// Post-process files containing volatile data.
	for _, file := range bashConfig.OutputProcessConfig.VolatileDataFiles {
		// Read the file.
		out, err := os.ReadFile(file)
		if err != nil {
			t.Fatal("error reading file: ", file, ": ", err)
		}
		got := string(out)

		// Replace default volatile content with a placeholder.
		if !bashConfig.OutputProcessConfig.KeepVolatileData {
			got = regexReplaceAllDefault(got)
		}
		// Apply custom volatile regex replacements.
		for _, r := range bashConfig.OutputProcessConfig.VolatileRegexReplacements {
			got = regexReplaceCustom(got, r.Replacement, r.Regex)
		}

		// Write the output back to the file.
		if err := os.WriteFile(file, []byte(got), 0o644); err != nil {
			t.Fatal("error writing stabilized output to file: ", err)
		}
	}
}

func runCmd(cmd *exec.Cmd, displayStdout, displayStderr bool) (output []byte, err error) {
	var out []byte
	if displayStderr && displayStdout {
		// If we are interested in displaying both stderr and stdout, we
		// need to get the combined output in order to get the correct
		// order of the output.
		if out, err = cmd.CombinedOutput(); err != nil {
			return nil, fmt.Errorf("error running command: %v: %s", err, out)
		}
		return out, nil
	}

	// Initialize variables that will hold the output of stdout and stderr.
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	// Assign the pointers of stderr and stdout to populate when the
	// command is run.
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	// Run the actual command.
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error running command: %v: %s %s", err, stderr.String(), stdout.String())
	}
	// Gather the output bytes. Writing stderr and stdout is optional.
	if displayStderr {
		out = stderr.Bytes()
	}
	if displayStdout {
		out = append(out, stdout.Bytes()...)
	}
	return out, nil
}
