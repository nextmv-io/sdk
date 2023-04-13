package nextroute

import (
	"time"

	"github.com/nextmv-io/sdk/connect"
)

// NewTimeDependentDurationExpression creates a new time dependent expression.
// The expression passed in is the default expression for the time dependent
// expression. The default expression is used when the time is not within any
// of the intervals set by SetExpression. Expressions can be defined for
// intervals of time. The intervals are inclusive of the start time and
// exclusive of the end time. For example, if the interval is [9:30, 10:30),
// the expression is valid for 9:30, 9:31, ..., 10:29:59:999. The interval must
// be on the minute boundary. The start time must be before the end time.
// The interval defined my the minimum of all start times and the maximum
// of all end times can not exceed more than one week.
func NewTimeDependentDurationExpression(
	model Model,
	defaultExpression DurationExpression,
) (TimeDependentDurationExpression, error) {
	connect.Connect(con, &newTimeDependentDurationExpression)
	return newTimeDependentDurationExpression(model, defaultExpression)
}

// NewTimeInDependentDurationExpression creates a new time in-dependent
// duration expression.
// This expression has the same interface as the time dependent expression
// but the time is not used in any of the calculations. All values originate
// from the base expression.
func NewTimeInDependentDurationExpression(
	expression DurationExpression,
) TimeDependentDurationExpression {
	connect.Connect(con, &newTimeInDependentDurationExpression)
	return newTimeInDependentDurationExpression(expression)
}

// TimeDependentDurationExpression is a DurationExpression that returns a value
// based on time on top of a base expression.
type TimeDependentDurationExpression interface {
	DurationExpression

	// DefaultExpression returns the default expression.
	DefaultExpression() DurationExpression

	// IsDependentOnTime returns true if the expression is dependent on time.
	// The expression is dependent on time if the expression is not the same
	// for all time intervals.
	IsDependentOnTime() bool

	// SetExpression sets the expression for the given time interval
	// [start, end). If the interval overlaps with an existing interval,
	// the existing interval is replaced. Both start and end must be on the
	// minute boundary.
	SetExpression(start, end time.Time, expression DurationExpression) error

	// ExpressionAtTime returns the expression for the given time.
	ExpressionAtTime(time.Time) DurationExpression
	// ExpressionAtValue returns the expression for the given value.
	ExpressionAtValue(float64) DurationExpression

	// ValueAtTime returns the value for the given time.
	ValueAtTime(time time.Time, vehicleType ModelVehicleType, from, to ModelStop) float64
	// ValueAtValue returns the value for the given value.
	ValueAtValue(value float64, vehicleType ModelVehicleType, from, to ModelStop) float64
}
