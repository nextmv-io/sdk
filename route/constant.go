package route

// ConstantByPoint measure always estimates the same cost.
func ConstantByPoint(c float64) ByPoint {
	connect()
	return constantByPointFunc(c)
}

// Constant measure always estimates the same cost.
func Constant(c float64) ByIndex {
	connect()
	return constantFunc(c)
}

var (
	constantByPointFunc func(float64) ByPoint
	constantFunc        func(float64) ByIndex
)
