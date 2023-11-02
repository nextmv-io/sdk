package run

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/nextmv-io/go-flagsfiller"
	"github.com/nextmv-io/sdk"
	"github.com/nextmv-io/sdk/types"
)

func init() {
	flagsfiller.RegisterSimpleType[types.Duration](customDuration)
}

func customDuration(s string, _ reflect.StructTag) (types.Duration, error) {
	d, err := time.ParseDuration(s)
	if err != nil {
		return types.Duration(d), err
	}
	return types.Duration(d), nil
}

// FlagParser parses flags and env vars and returns a runner config and options.
func FlagParser[Option, RunnerCfg any]() (
	runnerConfig RunnerCfg, option Option, err error,
) {
	// create a FlagSetFiller
	filler := flagsfiller.New(
		flagsfiller.WithEnv(""),
		flagsfiller.WithFieldRenamer(
			func(name string) string {
				repl := strings.ReplaceAll(name, "-", ".")
				return strings.ToLower(repl)
			},
		),
	)

	err = filler.Fill(flag.CommandLine, &option)
	if err != nil {
		return runnerConfig, option, err
	}

	err = filler.Fill(flag.CommandLine, &runnerConfig)
	if err != nil {
		return runnerConfig, option, err
	}
	flag.Usage = usage
	flag.Parse()

	return runnerConfig, option, nil
}

func usage() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	out := fs.Output()

	fmt.Fprintf(
		out,
		"Nextmv Hybrid Optimization Platform %s\n",
		sdk.VERSION,
	)
	fmt.Fprint(out, "Usage:\n")
	flag.PrintDefaults()
}
