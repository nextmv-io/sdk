package context

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Declare declares new data on a context.
func Declare[T any](ctx Context, data T) Declared[T] {
	connect()
	return declaredProxy[T]{declared: declareFunc(ctx, data)}
}

// Declared represents data declared on a context that can be queried and
// modified.
type Declared[T any] interface {
	Number() uint64
	Get(Context) T
	Set(T) Change
}

type declaredProxy[T any] struct {
	declared Declared[any]
}

func (d declaredProxy[T]) Number() uint64 {
	return d.declared.Number()
}

func (d declaredProxy[T]) Get(ctx Context) T {
	return d.declared.Get(ctx).(T)
}

func (d declaredProxy[T]) Set(data T) Change {
	return d.declared.Set(data)
}


func (d declaredProxy[T]) String() string {
	return fmt.Sprintf("<%d>%v", d.Number(), d.typeOf())
}

func (d declaredProxy[T]) MarshalJSON() ([]byte, error) {
	// A declared variable only has value in a context.
	return json.Marshal(
		map[string]any{
			"number": d.Number(),
			"type":   d.typeOf().String(),
		},
	)
}

func (d declaredProxy[T]) typeOf() reflect.Type {
	var v T
	return reflect.TypeOf(v)
}

var declareFunc func(Context, any) Declared[any]
