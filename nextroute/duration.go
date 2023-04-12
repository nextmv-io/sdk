package nextroute

import (
	"time"

	"github.com/nextmv-io/sdk/connect"
)

// DurationValue returns the value of a duration in the given time unit.
// Will panic if the time unit is zero.
func DurationValue(distance Distance, speed Speed, unit time.Duration) float64 {
	connect.Connect(con, &durationValue)
	return durationValue(distance, speed, unit)
}
