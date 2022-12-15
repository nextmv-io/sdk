package run

import (
	"context"
	"io"
	"log"
	"os"
)

// CliIOProducer is the IOProducer for the CliRunner. The input and output paths
// are used to configure the input and output readers and writers. If the paths
// are empty, os.Stdin and os.Stdout are used.
func CliIOProducer(_ context.Context, cfg CLIRunnerConfig) IOData {
	reader := os.Stdin
	if cfg.Runner.Input.Path != "" {
		r, err := os.Open(cfg.Runner.Input.Path)
		if err != nil {
			log.Fatal(err)
		}
		reader = r
	}
	var writer io.Writer = os.Stdout
	if cfg.Runner.Output.Path != "" {
		w, err := os.Create(cfg.Runner.Output.Path)
		if err != nil {
			log.Fatal(err)
		}
		writer = w
	}
	return NewIOData(
		reader,
		nil,
		writer,
	)
}
