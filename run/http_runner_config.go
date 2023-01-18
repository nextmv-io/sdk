package run

import (
	"log"
	"time"
)

// HTTPRunnerConfig defines the configuration of the HTTPRunner.
type HTTPRunnerConfig struct {
	Runner struct {
		Log    *log.Logger
		Output struct {
			Solutions string `default:"all" usage:"Return all or last solution"`
			Quiet     bool   `default:"false" usage:"Do not return statistics"`
		}
		HTTP struct {
			Address           string        `default:":9000" usage:"The host address"`
			Certificate       string        `usage:"The certificate file path"`
			Key               string        `usage:"The key file path"`
			ReadHeaderTimeout time.Duration `default:"60s" usage:"The maximum duration for reading the request headers"`
			MaxParallel       int           `default:"1" usage:"The max number of requests"`
		}
	}
}

// Quiet returns the quiet flag.
func (c HTTPRunnerConfig) Quiet() bool {
	return c.Runner.Output.Quiet
}

// Solutions returns the configured solutions.
func (c HTTPRunnerConfig) Solutions() (Solutions, error) {
	return ParseSolutions(c.Runner.Output.Solutions)
}
