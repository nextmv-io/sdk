package context

import (
	"encoding/json"
	"reflect"
)

type Declared[T any] interface {
	Declarable

	Get(Context) T
	Set(T) Change
}

type declaredProxy[T any] struct {
	declarable Declarable
}

func (d declaredProxy[T]) Number() uint64 {
	return d.declarable.Number()
}

func (d declaredProxy[T]) Type() reflect.Type {
	return d.declarable.Type()
}

func (d declaredProxy[T]) Get(ctx Context) T {
	if data, ok := declaredGetFunc(ctx, d.declarable.Number()); ok {
		return data.(T)
	}

	var x T // zero value of the generic type
	return x
}

func (d declaredProxy[T]) Set(data T) Change {
	return func(ctx Context) {
		declaredSetFunc(ctx, d.declarable, data)
	}
}

func (d declaredProxy[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.declarable)
}

var declaredGetFunc func(Context, uint64) (any, bool)
var declaredSetFunc func(Context, Declarable, any)
