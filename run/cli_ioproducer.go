package run

import (
	"context"
	"io"
	"log"
	"os"
)

// CliIOProducer is the default IOProducer for the CliRunner.
func CliIOProducer(_ context.Context, config any) IOData {
	cfg, ok := config.(CLIRunnerConfig)
	if !ok {
		log.Fatal("CliIOProducer is not compatible with the runner")
	}
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
