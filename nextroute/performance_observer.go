package nextroute

import (
	"time"

	"github.com/nextmv-io/sdk/connect"
)

// PerformanceObserver is an interface that is used to observe the performance
// of the model, and it's subsequent use.
type PerformanceObserver interface {
	SolutionObserver

	// Duration returns the duration since the creation of the
	// PerformanceObserver.
	Duration() time.Duration

	// Report creates a report of the performance of the model, and it's
	// subsequent use.
	Report() string
}

// NewPerformanceObserver creates a new PerformanceObserver. The performance
// observer is used to observe the performance of the model, and it's subsequent
// use. The created PerformanceObserver is not connected to the model, and
// therefore will not observe anything. To connect the PerformanceObserver to
// the model, use the AddSolutionObserver() method on the Model.
func NewPerformanceObserver(model Model) PerformanceObserver {
	connect.Connect(con, &newPerformanceObserver)
	return newPerformanceObserver(model)
}
