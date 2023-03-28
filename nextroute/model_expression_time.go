package nextroute

import (
	"time"

	"github.com/nextmv-io/sdk/connect"
)

// NewTimeExpression creates a new time expression.
func NewTimeExpression(
	expression ModelExpression,
	epoch time.Time,
) TimeExpression {
	connect.Connect(con, &newTimeExpression)
	return newTimeExpression(expression, epoch)
}

// NewStopTimeExpression creates a new duration expression.
func NewStopTimeExpression(
	name string,
	time time.Time,
) StopTimeExpression {
	connect.Connect(con, &newStopTimeExpression)
	return newStopTimeExpression(name, time)
}

// TimeExpression is a ModelExpression that returns a time.
type TimeExpression interface {
	ModelExpression
	// Time returns the time for the given parameters.
	Time(ModelVehicleType, ModelStop, ModelStop) time.Time
}

// StopTimeExpression is a ModelExpression that returns a time per stop and
// allows to set the time per stop.
type StopTimeExpression interface {
	ModelExpression
	// Time returns the time for the given stop.
	Time(stop ModelStop) time.Time
	// SetTime sets the time for the given stop.
	SetTime(ModelStop, time.Time)
}
