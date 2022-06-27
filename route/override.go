package route

// Override measure uses a default measure for all arcs that are not true for a
// condition. It uses an override measure for all arcs that are true for the
// condition.
func Override(
	defaultByIndex ByIndex,
	overrideByIndex ByIndex,
	condition func(from, to int) bool,
) ByIndex {
	return overrideFunc(defaultByIndex, overrideByIndex, condition)
}

var overrideFunc func(ByIndex,
	ByIndex, func(int, int) bool,
) ByIndex
