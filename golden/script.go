package golden

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/nextmv-io/sdk/flatmap"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type scriptTest struct {
	// Path is the path to the script file.
	Path string
	// Command is the command to execute the script.
	Command string
}

// BashTest calls ScriptTest with bash as the command to execute the scripts and
// the file extension .sh. This is a convenience function for bash scripts.
func BashTest(
	t *testing.T,
	goldenDir string,
	scriptConfig ScriptConfig,
) {
	scriptConfig.ScriptExtensions = map[string]string{
		".sh": "bash",
	}
	ScriptTest(t, goldenDir, scriptConfig)
}

// ScriptTest executes a golden file test for scripts. It walks over the
// goldenDir to gather all .sh scripts present in the dir. It then executes each
// of the scripts and compares expected vs. actual outputs. If displayStdout or
// displayStderr are true, the output of each script will be composed of the
// resulting stderr + stdout.
func ScriptTest(
	t *testing.T,
	goldenDir string,
	scriptConfig ScriptConfig,
) {
	// Fail immediately, if dir does not exist
	if stat, err := os.Stat(goldenDir); err != nil || !stat.IsDir() {
		t.Fatalf("dir %s does not exist", goldenDir)
	}

	// If no script extensions are provided, we fail the test.
	if len(scriptConfig.ScriptExtensions) == 0 {
		t.Fatal("no script extensions provided in script config")
	}

	// Collect scripts.
	extensions := make([]string, 0, len(scriptConfig.ScriptExtensions))
	for ext := range scriptConfig.ScriptExtensions {
		extensions = append(extensions, ext)
	}
	var scripts []scriptTest
	fn := func(path string, _ os.FileInfo, _ error) error {
		// Check if the file should be considered as a script to test.
		extension := filepath.Ext(path)
		if slices.Contains(extensions, extension) {
			scripts = append(scripts, scriptTest{
				Path:    path,
				Command: scriptConfig.ScriptExtensions[extension],
			})
		}
		return nil
	}
	if err := filepath.Walk(goldenDir, fn); err != nil {
		t.Fatal("error walking over files: ", err)
	}

	// If no scripts were found, fail the test.
	if len(scripts) == 0 {
		t.Fatal("no scripts found in directory: ", goldenDir)
	}

	// Execute a golden file test for each script. Make the script path
	// absolute to avoid issues with custom working directories.
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal("error getting current working directory: ", err)
	}
	for _, script := range scripts {
		ScriptTestFile(t, script.Command, filepath.Join(cwd, script.Path), scriptConfig)
	}

	// Post-process files containing volatile data.
	postProcessVolatileData(t, scriptConfig)
}

// BashTestFile executes a golden file test for a single bash script. The script
// is executed and the expected output is compared with the actual output.
func BashTestFile(
	t *testing.T,
	script string,
	bashConfig ScriptConfig,
) {
	ScriptTestFile(
		t,
		"bash",
		script,
		bashConfig,
	)
}

// ScriptTestFile executes a golden file test for a single script. The script is
// executed and the expected output is compared with the actual output.
func ScriptTestFile(
	t *testing.T,
	command string,
	script string,
	scriptConfig ScriptConfig,
) {
	ext := goldenExtension
	if scriptConfig.GoldenExtension != "" {
		ext = scriptConfig.GoldenExtension
	}
	goldenFilePath := script + ext
	// Function run by the test.
	f := func(t *testing.T) {
		// Make script path absolute to avoid issues with custom working
		// directories.
		var err error
		script, err = filepath.Abs(script)
		if err != nil {
			t.Fatal("error getting absolute path for script: ", script, ": ", err)
		}

		if _, err := os.Stat(script); err != nil {
			t.Fatalf("script %s does not exist", script)
		}

		// Execute the script using the provided command.
		cmd := exec.Command(command, script)

		// Pass environment and add custom environment variables
		cmd.Env = os.Environ()
		for _, e := range scriptConfig.Envs {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", e[0], e[1]))
		}

		// Set custom working directory if provided.
		if scriptConfig.WorkingDir != "" {
			cmd.Dir = scriptConfig.WorkingDir
		}

		// Run the command and gather the output bytes.
		out, err := runCmd(cmd, scriptConfig.DisplayStdout, scriptConfig.DisplayStderr)
		if err != nil {
			t.Fatal(err)
		}

		// Process the output data before comparison.
		got := processOutput(t, out, goldenFilePath, scriptConfig.OutputProcessConfig)

		// Write the output bytes to a .golden file, if the test is being
		// updated
		if *update || scriptConfig.OutputProcessConfig.AlwaysUpdate {
			if err := os.WriteFile(goldenFilePath, []byte(got), 0o644); err != nil {
				t.Fatal("error writing script output to file: ", err)
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

	// Delay the execution of the test to adhere for rate limits.
	if scriptConfig.WaitBefore > 0 {
		t.Logf("delaying test execution for %v", scriptConfig.WaitBefore)
		<-time.After(scriptConfig.WaitBefore)
	}

	// Test is executed.
	t.Run(script, f)

	// Run post-process functions.
	for _, f := range scriptConfig.PostProcessFunctions {
		err := f(goldenFilePath)
		if err != nil {
			t.Fatalf("error running post-process function: %v", err)
		}
	}
}

func processOutput(
	t *testing.T,
	out []byte,
	goldenPath string,
	config OutputProcessConfig,
) string {
	// Check whether any JSON modifications are requested.
	jsonModifications := false
	for _, field := range config.TransientFields {
		skip, err := skipFile(goldenPath, field.FileRegex, field.FileRegexFullPath)
		if err != nil {
			t.Fatalf("error checking file regex: %v", err)
		}
		if !skip {
			jsonModifications = true
			break
		}
	}
	for _, rounding := range config.RoundingConfig {
		skip, err := skipFile(goldenPath, rounding.FileRegex, rounding.FileRegexFullPath)
		if err != nil {
			t.Fatalf("error checking file regex: %v", err)
		}
		if !skip {
			jsonModifications = true
			break
		}
	}

	// Apply JSON specific processing (if any were found).
	if jsonModifications {
		// Convert the output to a map[string]any for processing.
		var err error
		output := map[string]any{}
		if err = json.Unmarshal(out, &output); err != nil {
			t.Fatalf("transient fields or rounding config provided, but output is not valid JSON: %v", err)
		}

		// Flatten the map and apply the configured replacements /
		// modifications.
		flattenedOutput := flatmap.Do(output)
		flattenedOutput, err = replaceTransient(goldenPath, flattenedOutput, config.TransientFields...)
		if err != nil {
			t.Fatal(err)
		}
		flattenedOutput, err = roundFields(goldenPath, flattenedOutput, config.RoundingConfig...)
		if err != nil {
			t.Fatal(err)
		}
		nestedOutput, err := flatmap.Undo(flattenedOutput)
		if err != nil {
			t.Fatal(err)
		}

		// Marshal the processed output back to JSON.
		out, err = json.MarshalIndent(nestedOutput, "", "  ")
		if err != nil {
			t.Fatal("error marshaling output: ", err)
		}
	}

	got := string(out)

	// Apply regex replacements for volatile data.
	if !config.KeepVolatileData {
		// Replace default volatile content with a placeholder.
		got = regexReplaceAllDefault(got)
	}
	// Apply custom regex replacements.
	for _, r := range config.VolatileRegexReplacements {
		got = regexReplaceCustom(got, r.Replacement, r.Regex)
	}

	return got
}

func postProcessVolatileData(
	t *testing.T,
	scriptConfig ScriptConfig,
) {
	// Post-process files containing volatile data.
	for _, file := range scriptConfig.OutputProcessConfig.VolatileDataFiles {
		// Read the file.
		out, err := os.ReadFile(file)
		if err != nil {
			t.Fatal("error reading file: ", file, ": ", err)
		}
		got := string(out)

		// Replace default volatile content with a placeholder.
		if !scriptConfig.OutputProcessConfig.KeepVolatileData {
			got = regexReplaceAllDefault(got)
		}
		// Apply custom volatile regex replacements.
		for _, r := range scriptConfig.OutputProcessConfig.VolatileRegexReplacements {
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
