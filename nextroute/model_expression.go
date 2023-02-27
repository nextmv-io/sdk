package nextroute

import (
	"fmt"
	"sync/atomic"
)

// ModelExpression is an expression that can be used in a model to define
// values for constraints and objectives. The expression is evaluated for
// each stop in the solution by invoking the Value() method. The value of
// the expression is then used in the constraints and objective.
type ModelExpression interface {
	// Index returns the unique index of the expression.
	Index() int

	// Name returns the name of the expression.
	Name() string

	// Value returns the value of the expression for the given vehicle type,
	// from stop and to stop.
	Value(ModelVehicleType, ModelStop, ModelStop) float64
}

// ModelExpressions is a slice of ModelExpression.
type ModelExpressions []ModelExpression

// DeltaValue is a helper function which returns the difference in expression
// value if candidate would be positioned after stop. Will panic if stop is
// not planned.
//
// The difference is calculated as follows:
//
//	v1 = expression.Value(vehicle, stop.ModelStop(), candidateStop)
//	v2 = expression.Value(vehicle, candidateStop, stop.Next().ModelStop())
//	v3 = expression.Value(vehicle, stop.ModelStop(), stop.Next().ModelStop())
//	return v1 + v2 - v3
func DeltaValue(
	stop SolutionStop,
	candidate SolutionStop,
	expression ModelExpression,
) float64 {
	if !stop.IsPlanned() {
		panic(fmt.Sprintf("stop %s is not planned", stop))
	}

	vehicle := stop.Vehicle().VehicleType()
	fromStop := stop.ModelStop()
	candidateStop := candidate.ModelStop()
	toStop := stop.Next().ModelStop()

	currentValue := expression.Value(
		vehicle,
		fromStop,
		toStop,
	)
	newValue1 := vehicle.TravelDurationExpression().Value(
		vehicle,
		fromStop,
		candidateStop,
	)
	newValue2 := vehicle.TravelDurationExpression().Value(
		vehicle,
		candidateStop,
		toStop,
	)
	return newValue1 + newValue2 - currentValue
}

var expressionIndex uint32

// NextExpressionIndex returns the next unique expression index.
func NextExpressionIndex() int {
	return int(atomic.AddUint32(&expressionIndex, 1) - 1)
}
