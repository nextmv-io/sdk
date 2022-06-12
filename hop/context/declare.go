package context

import (
	"encoding/json"
	"reflect"
)

// Declare declares new data on a context.
func Declare[T any](ctx Context, data T) Declared[T] {
	return declaredProxy[T]{declared: declareFunc(ctx, data)}
}

// Declared represents data declared on a context that can be queried and
// modified.
type Declared[T any] interface {
	Get(Context) T
	Set(T) Change
}

type declaredProxy[T any] struct {
	declared Declared[any]
}

func (d declaredProxy[T]) Get(ctx Context) T {
	return d.declared.Get(ctx).(T)
}

func (d declaredProxy[T]) Set(data T) Change {
	return d.declared.Set(data)
}

func (d declaredProxy[T]) String() string {
	var x T
	return reflect.TypeOf(x).String()
}

func (d declaredProxy[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

var declareFunc func(Context, any) Declared[any]
