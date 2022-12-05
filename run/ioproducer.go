package run

import (
	"context"
	"io"
	"log"
	"os"
)

// CliIOProducer is a test IOProducer.
func CliIOProducer(_ context.Context, config any) IOData {
	cfg, ok := config.(CliRunnerConfig)
	if !ok {
		log.Fatal("DefaultIOProducer is not compatible with the runner")
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
