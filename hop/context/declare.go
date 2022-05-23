package context

func Declare[T any](ctx Context, data T) Declared[T] {
	return declaredProxy[T]{declarable: declareFunc(ctx, data)}
}

var declareFunc func(Context, any) Declarable
