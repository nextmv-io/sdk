package nextroute

import (
	"github.com/nextmv-io/sdk/connect"
	"time"
)

// PerformanceObserver is an interface that is used to observe the performance
// of the model, and it's subsequent use.
type PerformanceObserver interface {
	SolutionObserver

	Duration() time.Duration

	Marshal() ([]byte, error)

	Report() string
}

func NewPerformanceObserver(model Model) PerformanceObserver {
	connect.Connect(con, &newPerformanceObserver)
	return newPerformanceObserver(model)
}
