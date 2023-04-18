package nextroute

import (
	"time"

	"github.com/nextmv-io/sdk/connect"
)

// NewTimeDependentDurationExpression creates a
// [TimeDependentDurationExpression], which takes a default expression as input.
// This default expression is used when the time is not within any of the
// intervals set by the SetExpression method.
//
// The [TimeDependentDurationExpression.SetExpression] method can be used to
// define expressions for intervals of time, which are inclusive of the start
// time and exclusive of the end time. For example, if the interval is
// [9:30, 10:30), the expression is valid for all times from 9:30 up to (but not
// including) 10:30. Note that the interval must be on the minute boundary, and
// the start time must be before the end time.
//
// The overall interval defined by the minimum of all start times and the
// maximum of all end times cannot exceed more than one week.
// The default expression is used for all times that go outside the defined
// intervals and may not contain any negative values.
func NewTimeDependentDurationExpression(
	model Model,
	defaultExpression DurationExpression,
) (TimeDependentDurationExpression, error) {
	connect.Connect(con, &newTimeDependentDurationExpression)
	return newTimeDependentDurationExpression(model, defaultExpression)
}

// NewTimeIndependentDurationExpression creates a
// [TimeDependentDurationExpression] which is not dependent on time.
// This expression has the same interface as the time dependent expression
// but the time is not used in any of the calculations. All values originate
// from the base expression.
func NewTimeIndependentDurationExpression(
	expression DurationExpression,
) TimeDependentDurationExpression {
	connect.Connect(con, &newTimeIndependentDurationExpression)
	return newTimeIndependentDurationExpression(expression)
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
	// minute boundary. Expression is not allowed to contain negative values.
	SetExpression(start, end time.Time, expression DurationExpression) error

	// ExpressionAtTime returns the expression for the given time.
	ExpressionAtTime(time.Time) DurationExpression
	// ExpressionAtValue returns the expression for the given value.
	ExpressionAtValue(float64) DurationExpression

	// ValueAtTime returns the value for the given time.
	ValueAtTime(
		time time.Time,
		vehicleType ModelVehicleType,
		from, to ModelStop,
	) float64
	// ValueAtValue returns the value for the given value.
	ValueAtValue(
		value float64,
		vehicleType ModelVehicleType,
		from, to ModelStop,
	) float64
}
