package run

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/itzg/go-flagsfiller"
	"github.com/nextmv-io/sdk"
)

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

	version := sdk.VERSION

	// internal test expectations use <<PRESENCE>> as the version so we do not
	// have to update expectations every time the version changes
	if ver, ok := os.LookupEnv("USE_PRESENCE"); ok && ver == "1" {
		version = "<<PRESENCE>>"
	}

	fmt.Fprintf(
		out,
		"\"Nextmv Hybrid Optimization Platform\" %s\n",
		version,
	)
	fmt.Fprint(out, "Usage:\n")
	flag.PrintDefaults()
}
