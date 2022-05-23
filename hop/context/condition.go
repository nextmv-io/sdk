package context

type Condition func(Context) bool

func And(c1 Condition, c2 Condition, conditions ...Condition) Condition {
	return andFunc(c1, c2, conditions...)
}

func False(ctx Context) bool {
	return falseFunc(ctx)
}

func Not(c Condition) Condition {
	return notFunc(c)
}

func Or(c1 Condition, c2 Condition, conditions ...Condition) Condition {
	return orFunc(c1, c2, conditions...)
}

func True(ctx Context) bool {
	return trueFunc(ctx)
}

func Xor(c1, c2 Condition) Condition {
	return xorFunc(c1, c2)
}

var (
	andFunc   func(Condition, Condition, ...Condition) Condition
	falseFunc func(Context) bool
	notFunc   func(Condition) Condition
	orFunc    func(Condition, Condition, ...Condition) Condition
	trueFunc  func(Context) bool
	xorFunc   func(Condition, Condition) Condition
)
