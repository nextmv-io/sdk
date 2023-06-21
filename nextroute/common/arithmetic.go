package common

// Numbers is a type constraint for all numeric types.
type Numbers interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

// Sum returns the sum of the given numbers.
func Sum[N Numbers](data []N) N {
	var sum N
	for _, d := range data {
		sum += d
	}
	return sum
}

// SumDefined returns the sum of the given numbers by applying the given function
// to each item in data.
func SumDefined[T any, N Numbers](data []T, f func(T) N) N {
	var sum N
	for _, d := range data {
		sum += f(d)
	}
	return sum
}
