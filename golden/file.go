// Package golden holds tools for testing documentation code.
package golden

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xeipuuv/gojsonschema"
)

// update is a flag that can be passed to the test binary to update the golden
// files.
var update = flag.Bool("update", false, "update goldenfiles and templates")

// FileTests performs golden file tests for all the <.json/csv> files contained
// in the given location. A golden file test uses an input file to execute a
// program. The output of the program is compared against an expected output
// which is a .golden file. If the location is a <.json/csv> file, then a
// single test is carried out for that file. If the location is a directory,
// all <.json/csv> files are gathered and a golden file test is executed for
// each file. You can read more about [GoldenFileTesting] in this blog.
//
// [GoldenFileTesting]: https://ieftimov.com/posts/testing-in-go-golden-files/
func FileTests(t *testing.T, location string, config Config) {
	if _, err := os.Stat(location); err != nil {
		t.Fatalf("file %s does not exist", location)
	}

	var inputs []string
	err := filepath.Walk(
		location,
		func(path string, info os.FileInfo, _ error) error {
			if !validForFileComparison(info) {
				return nil
			}

			inputs = append(inputs, path)
			return nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	for _, input := range inputs {
		FileTest(t, input, config)
	}
}

// FileTest performs a golden file test using input as the input file for the
// program found at the path given by main. The output is compared against the
// golden file (.golden). If the update flag is set, the .golden file is
// updated to the observed output.
func FileTest(t *testing.T, inputPath string, config Config) {
	t.Run(
		filepath.Base(inputPath),
		func(t *testing.T) {
			t.Parallel()
			command, tempFileName, err := config.entrypoint(inputPath)
			if err != nil {
				t.Fatal(err)
			}

			if !config.UseStdOut {
				defer func() {
					err := os.Remove(tempFileName)
					if err != nil {
						panic(err)
					}
				}()
			}

			// Run the actual command.
			var stdout, stderr bytes.Buffer
			command.Stderr = &stderr
			command.Stdout = &stdout
			if err := command.Run(); err != nil {
				t.Fatal(err, stdout.String(), stderr.String())
			}

			if exitCode := command.ProcessState.ExitCode(); exitCode != config.ExitCode {
				t.Fatalf("got %v; want %v", exitCode, config.ExitCode)
			}

			actualBytes := stdout.Bytes()
			if !config.UseStdOut {
				if len(actualBytes) > 0 {
					t.Fatal("expected no stdout bytes but got: ", string(actualBytes))
				}

				bytes, err := os.ReadFile(tempFileName)
				if err != nil {
					t.Fatal(err)
				}

				actualBytes = bytes
			}

			if !config.SkipGoldenComparison {
				comparison(
					t,
					actualBytes,
					inputPath,
					config,
				)
			}

			if len(config.InputSchema) > 0 {
				input, err := os.ReadFile(inputPath)
				if err != nil {
					t.Errorf("got %v; want nil", err)
				}

				err = ValidateAgainstSchema("input", input, config.InputSchema)
				if err != nil {
					t.Error(err)
				}
			}

			if len(config.OutputSchema) > 0 {
				err := ValidateAgainstSchema("output", actualBytes, config.OutputSchema)
				if err != nil {
					t.Error(err)
				}
			}

			if config.VerifyFunc != nil {
				input, err := os.ReadFile(inputPath)
				if err != nil {
					t.Errorf("got %v; want nil", err)
				}

				err = config.VerifyFunc(input, actualBytes)
				if err != nil {
					t.Errorf("verification error: %v", err)
				}
			}
		},
	)
}

// comparison verifies that the given actual bytes are matching the expectation
// given by the goldenPath.
func comparison(
	t *testing.T,
	actualBytes []byte,
	inputPath string,
	config Config,
) {
	goldenPath := inputPath + ".golden"
	if config.OutputProcessConfig.RelativeDestination != "" {
		goldenPath = filepath.Join(
			config.OutputProcessConfig.RelativeDestination,
			filepath.Base(inputPath)+".golden",
		)
	}

	outputWithTransient := map[string]any{}
	flattenedOutput := map[string]any{}
	if !config.CompareConfig.TxtParse {
		output := map[string]any{}
		if err := json.Unmarshal(actualBytes, &output); err != nil {
			t.Fatal(err)
		}

		flattenedOutput = flatten(output)
		flattenedOutput = replaceTransient(flattenedOutput, config.TransientFields...)
		nestedOutput, err := nest(flattenedOutput)
		if err != nil {
			t.Fatal(err)
		}
		outputWithTransient = nestedOutput
	}

	// Update golden file, if requested. If we are updating the golden file, we
	// don't need to compare against it.
	if *update || config.OutputProcessConfig.AlwaysUpdate {
		updateGoldenFile(t, config, actualBytes, outputWithTransient, goldenPath)
		return
	}

	expectedBytes, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatal(err)
	}

	if config.CompareConfig.TxtParse {
		actual := string(actualBytes)
		expected := string(expectedBytes)
		if config.CompareConfig.TxtCompareLength == 0 && actual != expected {
			t.Errorf("\ngot:\n\"%s\"\nexpected:\n\"%s\"", actual, expected)
			return
		}

		chopped := actual
		if config.CompareConfig.TxtCompareLength > 0 {
			chopped = actual[0:config.CompareConfig.TxtCompareLength]
		}

		if chopped != expected {
			t.Errorf("\ngot:\n\"%s\"\nexpected:\n\"%s\"", chopped, expected)
		}

		return
	}

	// Default comparison is JSON.
	expected := map[string]any{}
	if err := json.Unmarshal(expectedBytes, &expected); err != nil {
		t.Fatal(err)
	}

	flattenedExpected := flatten(expected)

	if len(config.DedicatedComparison) > 0 {
		for _, key := range config.DedicatedComparison {
			outputValue, ok := flattenedOutput[key]
			if !ok {
				t.Errorf("key \"%s\" not present in output", key)
			}

			expectedValue, ok := flattenedExpected[key]
			if !ok {
				t.Errorf("key \"%s\" not present in expected", key)
			}

			if err := valuesAreEqual(config, key, outputValue, expectedValue); err != nil {
				t.Error(err)
			}
		}
		return
	}

	for key, outputValue := range flattenedOutput {
		expectedValue, ok := flattenedExpected[key]
		if !ok {
			t.Errorf("key \"%s\" present in output but not expected", key)
		}

		if err := valuesAreEqual(config, key, outputValue, expectedValue); err != nil {
			t.Error(err)
		}
	}
}

// updateGoldenFile updates the golden file given by the goldenPath. If the
// config states that it is a raw text comparison, the actualBytes are written
// to the golden file. On the other hand, if the config states that it is a
// JSON comparison, the jsonOutput is written to the golden file.
func updateGoldenFile(
	t *testing.T,
	config Config,
	actualBytes []byte,
	jsonOutput map[string]any,
	goldenPath string,
) {
	var updated []byte
	var err error
	if config.CompareConfig.TxtParse {
		// Replace any volatile content with a placeholder
		actualString := string(actualBytes)
		actualString = regexReplaceAllDefault(actualString)
		for _, r := range config.OutputProcessConfig.VolatileRegexReplacements {
			actualString = regexReplaceCustom(actualString, r.Replacement, r.Regex)
		}

		actualBytes = []byte(actualString)
		updated = actualBytes
		if config.CompareConfig.TxtCompareLength != 0 {
			actualString := string(actualBytes)
			choppedString := actualString[0:config.CompareConfig.TxtCompareLength]
			if len(actualString) > len(choppedString) {
				choppedString += "\n"
			}

			updated = []byte(choppedString)
		}
	} else {
		updated, err = json.MarshalIndent(jsonOutput, "", "  ")
		updated = append(updated, '\n')
		if err != nil {
			t.Fatal(err)
		}
	}

	err = os.WriteFile(goldenPath, updated, 0o644)
	if err != nil {
		t.Fatal(err)
	}
}

// ValidateAgainstSchema validates the given JSON document (docBytes) against a
// JSON schema (schemaBytes). 'name' can be used for identifying the failing
// document type in the error.
func ValidateAgainstSchema(name string, docBytes, schemaBytes []byte) error {
	// Validate JSON file against given schema
	schemaLoader := gojsonschema.NewStringLoader(string(schemaBytes))
	documentLoader := gojsonschema.NewStringLoader(string(docBytes))
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}

	// Collect and return errors (or simply success in form of nil)
	errors := []string{}
	for _, desc := range result.Errors() {
		errors = append(errors, fmt.Sprintf("- %s\n", desc))
	}
	if result.Valid() {
		return nil
	}
	return fmt.Errorf(
		"schema validation error for %s: \n%s",
		name,
		strings.Join(errors, "\n"),
	)
}

// valuesAreEqual compares two values and returns an error if they are not.
// Comparison is done based on a treshold value, which is defined in the
// GoldenConfig.
func valuesAreEqual(config Config, key string, output, expected any) error {
	if outputString, isString := output.(string); isString {
		expectedString, expectedIsString := expected.(string)
		if !expectedIsString {
			return fmt.Errorf(
				"key \"%s\": output value %v is string, expected value %v is not",
				key, outputString, expected,
			)
		}

		if outputTime, err := time.Parse(time.RFC3339, outputString); err == nil {
			expectedTime, err := time.Parse(time.RFC3339, expectedString)
			if err != nil {
				return fmt.Errorf(
					"key \"%s\": output value %v is time.Time, expected value %v is not",
					key, outputTime, expectedString,
				)
			}

			threshold := config.Thresholds.Time
			if customThreshold, ok := config.Thresholds.CustomThresholds.Time[key]; ok {
				threshold = customThreshold
			}

			if !withinToleranceTime(outputTime, expectedTime, threshold) {
				return fmt.Errorf(
					"key \"%s\" (time threshold comparison): output %v, expected %v, using threshold %v",
					key, outputTime, expectedTime, threshold,
				)
			}

			return nil
		}

		if outputDuration, err := time.ParseDuration(outputString); err == nil {
			expectedDuration, err := time.ParseDuration(expectedString)
			if err != nil {
				return fmt.Errorf(
					"key \"%s\": output value %v is time.Duration, expected value %v is not",
					key, outputDuration, expectedString,
				)
			}

			threshold := config.Thresholds.Duration
			if customThreshold, ok := config.Thresholds.CustomThresholds.Duration[key]; ok {
				threshold = customThreshold
			}

			if !withinTolerance(outputDuration, expectedDuration, threshold) {
				return fmt.Errorf(
					"key \"%s\" (duration threshold comparison): output %v, expected %v, using threshold %v",
					key, outputDuration, expectedDuration, threshold,
				)
			}

			return nil
		}

		if outputString != expectedString {
			return fmt.Errorf(
				"key \"%s\" (string comparison): output %v, expected %v",
				key, outputString, expectedString,
			)
		}

		return nil
	}

	if outputFloat, isFloat := output.(float64); isFloat {
		expectedFloat, expectedIsFloat := expected.(float64)
		if !expectedIsFloat {
			return fmt.Errorf(
				"key \"%s\": output value %v is float64, expected value %v is not",
				key, outputFloat, expected,
			)
		}

		threshold := config.Thresholds.Float
		if customThreshold, ok := config.Thresholds.CustomThresholds.Float[key]; ok {
			threshold = customThreshold
		}

		if !withinTolerance(outputFloat, expectedFloat, threshold) {
			return fmt.Errorf(
				"key \"%s\" (float threshold comparison): output %v, expected %v, using threshold %v",
				key, outputFloat, expectedFloat, threshold,
			)
		}

		return nil
	}

	if outputInt, IsInt := output.(int); IsInt {
		expectedInt, expectedIsInt := expected.(int)
		if !expectedIsInt {
			return fmt.Errorf(
				"key \"%s\": output value %v is int, expected value %v is not",
				key, outputInt, expected,
			)
		}

		threshold := config.Thresholds.Int
		if customThreshold, ok := config.Thresholds.CustomThresholds.Int[key]; ok {
			threshold = customThreshold
		}

		if !withinTolerance(outputInt, expectedInt, threshold) {
			return fmt.Errorf(
				"key \"%s\" (integer threshold comparison): output %v, expected %v, using threshold %v",
				key, outputInt, expectedInt, threshold,
			)
		}

		return nil
	}

	if outputBool, isBool := output.(bool); isBool {
		expectedBool, expectedIsBool := expected.(bool)
		if !expectedIsBool {
			return fmt.Errorf(
				"key \"%s\": output value %v is bool, expected value %v is not",
				key, outputBool, expected,
			)
		}

		if outputBool != expectedBool {
			return fmt.Errorf(
				"key \"%s\" (boolean comparison): output %v, expected %v",
				key, outputBool, expectedBool,
			)
		}
	}

	return nil
}

// withinTolerance compares two values of type int, float64 or time.Duration
// and returns true if the difference between the two values is within the
// given tolerance.
func withinTolerance[T int | float64 | time.Duration](a, b, tolerance T) bool {
	if a == b {
		return true
	}

	if a > b {
		diff := a - b
		return diff <= tolerance
	}

	diff := b - a
	return diff <= tolerance
}

// withinToleranceTime compares two values of type time.Time and returns true
// if the difference between the two values is within the given tolerance.
func withinToleranceTime(a, b time.Time, tolerance time.Duration) bool {
	if a == b {
		return true
	}

	if a.After(b) {
		diff := a.Sub(b)
		return diff <= tolerance
	}

	diff := b.Sub(a)
	return diff <= tolerance
}

// validForFileComparison returns true if the given fileInfo is valid for
// golden file comparison.
func validForFileComparison(fileInfo os.FileInfo) bool {
	if fileInfo.IsDir() {
		return false
	}

	validSuffixes := []string{
		".json",
	}
	for _, suffix := range validSuffixes {
		if !strings.HasSuffix(fileInfo.Name(), suffix) {
			return false
		}
	}

	return true
}
