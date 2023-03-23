package nextroute

import (
	"time"

	"github.com/nextmv-io/sdk/connect"
)

// PerformanceObserver is an interface that is used to observe the performance
// of the model, and it's subsequent use.
type PerformanceObserver interface {
	SolutionObserver

	Duration() time.Duration

	Report() string
}

// NewPerformanceObserver returns a new performance observer.
func NewPerformanceObserver(model Model) PerformanceObserver {
	connect.Connect(con, &newPerformanceObserver)
	return newPerformanceObserver(model)
}
