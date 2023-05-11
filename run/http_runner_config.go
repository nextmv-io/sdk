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
			Solutions string `default:"last" usage:"Return all or last solution"`
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

// Solutions returns the configured solutions.
func (c HTTPRunnerConfig) Solutions() (Solutions, error) {
	return ParseSolutions(c.Runner.Output.Solutions)
}
