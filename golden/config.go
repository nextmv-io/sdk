package golden

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const goldenExtension = ".golden"

// Config lets a user configure the golden file tests.
type Config struct {
	// VerifyFunc is used to validate output against input, if provided.
	VerifyFunc func(input, output []byte) error
	// InputSchema definition to validate input JSON against
	InputSchema []byte
	// OutputSchema definition to validate output JSON against
	OutputSchema []byte
	// Args specifies the arguments to supply to the solver.
	Args []string
	// Envs specifies the environment variables to set for execution.
	Envs [][2]string
	// CompareConfig defines how to compare output against expectation.
	CompareConfig CompareConfig
	// OutputProcessConfig defines how to process the output before comparison.
	OutputProcessConfig OutputProcessConfig
	// SkipGoldenComparison skips the comparison against the golden file.
	SkipGoldenComparison bool
	// ExitCode defines the expected exit code of the command.
	ExitCode int
	// UseStdIn indicates whether to feed the file via stdin.
	UseStdIn bool
	// UseStdOut indicates whether to write output to stdout instead of a file.
	UseStdOut bool
	// IgnoreStdOut indicates whether to ignore the output of the command. This
	// is useful when the command writes to a file and we are not interested in
	// the output.
	IgnoreStdOut bool
	// TransientFields are keys that hold values which are transient (dynamic)
	// in nature, such as the elapsed time, version, start time, etc. Transient
	// fields have a special parsing in the .golden file and they are
	// stabilized in the comparison.
	TransientFields []TransientField
	// Tresholds by data type to be used when comparing actual and expected.
	// This configuration is optional, and if not provided then the comparison
	// between values is hard equality.
	Thresholds Tresholds
	// When DedicatedComparison is defined, then the golden file test will only
	// compare the keys that are defined in the slice. The keys are defined as
	// a [JSONPath]-like key. In general, use a dot (.) to recursively enter
	// nested objects and brackets ([]) to access array elements.
	//
	// [JSONPath]: https://goessner.net/articles/JsonPath/
	DedicatedComparison []string
	// ExecutionConfig defines the configuration for a Python golden file test. If
	// it is absent, then the golden file test is not a Python test.
	ExecutionConfig *ExecutionConfig
}

// BashConfig defines the configuration for a golden bash test.
type BashConfig struct {
	// DisplayStdout indicates whether to display or suppress stdout.
	DisplayStdout bool
	// DisplayStderr indicates whether to display or suppress stderr.
	DisplayStderr bool
	// OutputProcessConfig defines how to process the output before comparison.
	OutputProcessConfig OutputProcessConfig
	// Envs specifies the environment variables to set for execution.
	Envs [][2]string
	// PostProcessFunctions defines a list of functions to be executed after the bash
	// script has been run. This can be used to make use of the output of the bash script
	// and perform additional operations on it. The functions are executed in the order
	// they are defined and are not used for comparison.
	PostProcessFunctions []func(goldenFile string) error
}

// TransientField represents a field that is transient, this is, dynamic in
// nature. Examples of such fields include durations, times, versions, etc.
// Transient fields are replaced in golden file tests to always obtain the same
// result regardless of the moment it is executed. If a dynamic field always
// has a static value, then the golden file tests can run successfully. The
// transient field is represented by a key and a replacement. Please see the
// documentation of the fields for more information.
type TransientField struct {
	// Key is a representation of a [JSONPath]-like key. Here are some examples
	// of transient fields and how to override them in the comparison:
	//  - "version" key in the root: ".version"
	//  - "elapsed" key in the stats object in the root: ".stats.elapsed"
	//  - "start" key in the stats object in the root: ".stats.start"
	//  - "time" key in the first element of the solutions array in the root: ".solutions[0].time"
	// In general, use a dot (.) to recursively enter nested objects and
	// brackets ([]) to access array elements.
	//
	// [JSONPath]: https://goessner.net/articles/JsonPath/
	Key string

	// Replacement is optional, and it is the value that is used to stabilize
	// the transient field. If a replacement is not provided for the key, the
	// stabilization happens according to the data type. For example, a
	// [time.Time] is replaced using [StableTime], a [time.Duration] is
	// replaced using [StableDuration], etc. You can use the constants provided
	// by this package to stabilize the transient fields.
	Replacement any
}

// Tresholds by data type to be used when comparing actual and expected. If the
// absolute difference between the two values is less than or equal to the
// given threshold, then we consider the two values to be equal.
type Tresholds struct {
	// Float is the threshold to be used when comparing floats.
	Float float64
	// Int is the threshold to be used when comparing ints.
	Int int
	// Time is the threshold to be used when comparing times. Two times are
	// considered the same if the absolute difference between them, which is a
	// duration, is less than or equal to the given threshold.
	Time time.Duration
	// Duration is the threshold to be used when comparing durations.
	Duration time.Duration
	// CustomThresholds defines threshold for specific keys that override the
	// generic thresholds.
	CustomThresholds CustomThresholds
}

// CustomThresholds defines threshold for specific keys that override the
// generic thresholds.
type CustomThresholds struct {
	// Float defines specific thresholds for specific keys.
	Float map[string]float64
	// Int defines specific thresholds for specific keys.
	Int map[string]int
	// Time defines specific thresholds for specific keys.
	Time map[string]time.Duration
	// Duration defines specific thresholds for specific keys.
	Duration map[string]time.Duration
}

// CompareConfig configures how to compare actual and expected.
type CompareConfig struct {
	// Pure string comparison. If true, the output is compared as a string. If
	// false, the output is parsed as JSON and compared as a JSON object.
	TxtParse         bool
	TxtCompareLength int
}

// OutputProcessConfig defines how to process the output before comparison.
type OutputProcessConfig struct {
	// AlwaysUpdate makes the comparison always update the golden file.
	AlwaysUpdate bool
	// KeepVolatileData indicates whether to keep or replace frequently
	// changing data.
	KeepVolatileData bool
	// VolatileRegexReplacements defines regex replacements to be applied to the
	// golden file before comparison.
	VolatileRegexReplacements []VolatileRegexReplacement
	// VolatileDataFiles are files that contain volatile data and should get
	// post-processed to be more stable. This is only supported in directory
	// mode ([BashTest]) of golden bash testing, i.e., this will be ignored in
	// single file mode ([BashTestFile]).
	VolatileDataFiles []string
	// RelativeDestination is the relative path to the directory where the
	// output file will be stored. If not provided, then the output file is
	// stored in the current directory.
	RelativeDestination string
}

// ExecutionConfig defines the configuration for non-SDK golden file tests.
type ExecutionConfig struct {
	// Command is the command of the entrypoint of the app to be executed. E.g.,
	// "python3".
	Command string
	// Args are the arguments to be passed to the entrypoint of the app to be
	// executed. E.g., ["main.py"]. If InputFlag and OutputFlag are not
	// specified, then GOLDEN_INPUT and GOLDEN_OUTPUT are replaced by the input
	// and output file paths, respectively.
	Args []string
	// WorkDir is the working directory where the command will be executed. When
	// specified, the input file path will be adapted.
	WorkDir string
	// InputFlag is the argument to be used to pass the input file to the app to
	// be executed. E.g., "-input".
	InputFlag string
	// OutputFlag is the argument to be used to pass the output file to the app
	// to be executed. E.g., "-output".
	OutputFlag string
}

const (
	// ArgInputReplacement is the placeholder for the input file path in the
	// command arguments.
	ArgInputReplacement = "GOLDEN_INPUT"
	// ArgOutputReplacement is the placeholder for the output file path in the
	// command arguments.
	ArgOutputReplacement = "GOLDEN_OUTPUT"
)

// entrypoint returns the command to execute the algorithm for golden file
// comparison and the name of a temporary file where the output will be stored,
// according to the language configured by the config struct.
func (config Config) entrypoint(inputPath string) (*exec.Cmd, string, error) {
	var tempFileName string
	isCustom := config.ExecutionConfig != nil
	args := config.Args

	// Adapt input path, if using custom working directory
	if isCustom && config.ExecutionConfig.WorkDir != "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, "", err
		}
		inputPath = filepath.Join(cwd, inputPath)
	}

	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return nil, "", fmt.Errorf("input file does not exist: %s", inputPath)
	}

	// Append custom arguments to the command arguments, if using custom command
	if isCustom {
		args = append(config.ExecutionConfig.Args, args...)
	}

	// Handle how the input is passed to the command. If we are using stdin,
	// then we can ignore the input path here. Otherwise, we either use the
	// specified input flag or replace the placeholder in the command arguments.
	if !config.UseStdIn {
		if isCustom && config.ExecutionConfig.InputFlag == "" {
			for i, arg := range args {
				if strings.Contains(arg, ArgInputReplacement) {
					args[i] = strings.ReplaceAll(arg, ArgInputReplacement, inputPath)
				}
			}
		} else {
			inputFlag := "-runner.input.path"
			if isCustom && config.ExecutionConfig.InputFlag != "" {
				inputFlag = config.ExecutionConfig.InputFlag
			}
			args = append(args, inputFlag, inputPath)
		}
	}

	// Handle how the output is passed to the command. If we are using stdout,
	// then we can ignore the output path here. Otherwise, we either use the
	// specified output flag or replace the placeholder in the command
	// arguments.
	if !config.UseStdOut {
		outputFile, err := os.CreateTemp("", "output")
		if err != nil {
			return nil, "", err
		}
		tempFileName = outputFile.Name()
		if isCustom && config.ExecutionConfig.OutputFlag == "" {
			for i, arg := range args {
				if strings.Contains(arg, ArgOutputReplacement) {
					args[i] = strings.ReplaceAll(arg, ArgOutputReplacement, tempFileName)
				}
			}
		} else {
			outputFlag := "-runner.output.path"
			if isCustom && config.ExecutionConfig.OutputFlag != "" {
				outputFlag = config.ExecutionConfig.OutputFlag
			}
			args = append(args, outputFlag, tempFileName)
		}
	}

	// Assemble the command (switch working directory if needed)
	command := exec.Command("./"+binaryName, args...)
	if isCustom {
		command = exec.Command(config.ExecutionConfig.Command, args...)
		if config.ExecutionConfig.WorkDir != "" {
			command.Dir = config.ExecutionConfig.WorkDir
		}
	}

	// Pass environment and add custom environment variables
	command.Env = os.Environ()
	for _, e := range config.Envs {
		command.Env = append(command.Env, fmt.Sprintf("%s=%s", e[0], e[1]))
	}

	// Pipe input file to stdin, if using stdin
	if config.UseStdIn {
		file, err := os.Open(inputPath)
		if err != nil {
			return nil, "", err
		}
		command.Stdin = file
	}

	return command, tempFileName, nil
}
