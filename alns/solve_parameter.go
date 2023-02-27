package alns

// SolveParameter is an interface for a parameter that can change
// during the solving. The parameter can be used to control the
// behavior of the solver and it's operators.
type SolveParameter[T baseSolution[T]] interface {
	// Update updates the parameter based on the given solve information.
	// Update is invoked after each iteration of the solver.
	Update(SolveInformation[T])

	// Value returns the current value of the parameter.
	Value() int
}

// SolveParameters is a slice of solve parameters.
type SolveParameters[T baseSolution[T]] []SolveParameter[T]

// NewConstSolveParameter creates a new constant solve parameter.
func NewConstSolveParameter[T baseSolution[T]](value int) SolveParameter[T] {
	return &constParameterImpl[T]{value: value}
}

// NewSolveParameter creates a new solve parameter.
//   - startValue is the initial value of the parameter.
//   - deltaAfterIterations is the number of iterations without an improvement
//     before the value is changed.
//   - delta is the initial change in value after deltaAfterIterations.
//   - minValue is the minimum value of the parameter.
//   - maxValue is the maximum value of the parameter.
//   - snapBackAfterImprovement is a flag that indicates if the value should
//     snap back to the start value after an improvement.
//   - zigzag is a flag that indicates if the value should zigzag between
//     the min and max value. If the value is at the min value and delta is
//     negative, the delta is changed to positive. If the value is at the
//     max value and delta is positive, the delta is changed to negative.
func NewSolveParameter[T baseSolution[T]](
	startValue int,
	deltaAfterIterations int,
	delta int,
	minValue int,
	maxValue int,
	snapBackAfterImprovement bool,
	zigzag bool,
) SolveParameter[T] {
	if startValue == maxValue && delta < 0 {
		delta = -delta
	}
	if startValue == minValue && delta > 0 {
		delta = -delta
	}
	return &intParameterImpl[T]{
		startValue:               startValue,
		startDelta:               delta,
		deltaAfterIterations:     deltaAfterIterations,
		delta:                    delta,
		minValue:                 minValue,
		maxValue:                 maxValue,
		value:                    startValue,
		snapBackAfterImprovement: snapBackAfterImprovement,
		zigzag:                   zigzag,
	}
}

type intParameterImpl[T baseSolution[T]] struct {
	startValue               int
	startDelta               int
	deltaAfterIterations     int
	delta                    int
	maxValue                 int
	minValue                 int
	value                    int
	snapBackAfterImprovement bool
	zigzag                   bool
	iterations               int
}

func (i *intParameterImpl[T]) Value() int {
	return i.value
}

func (i *intParameterImpl[T]) Update(solveInformation SolveInformation[T]) {
	if solveInformation.DeltaScore() < 0.0 {
		i.iterations = 0
		if i.snapBackAfterImprovement && i.value != i.startValue {
			i.delta = i.startDelta
			i.value = i.startValue
		}
		return
	}
	i.iterations++
	if i.iterations > i.deltaAfterIterations {
		if i.value == i.maxValue || i.value == i.minValue {
			i.delta = -i.delta
		}

		i.iterations = 0
		i.value += i.delta
		if i.value > i.maxValue {
			i.value = i.maxValue
		}
		if i.value < i.minValue {
			i.value = i.minValue
		}
	}
}

type constParameterImpl[T baseSolution[T]] struct {
	value int
}

func (c *constParameterImpl[T]) Update(_ SolveInformation[T]) {
}

func (c *constParameterImpl[T]) Value() int {
	return c.value
}
