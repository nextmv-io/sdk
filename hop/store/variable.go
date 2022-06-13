package store

import (
	"encoding/json"
	"reflect"
)

// Variable that can be stored in a Store.
type Variable[T any] interface {
	Get(Store) T
	Set(T) Change
}

// Var stores a new Variable in a Store.
func Var[T any](s Store, data T) Variable[T] {
	return variable[T]{variable: varFunc(s, data)}
}

type variable[T any] struct {
	variable Variable[any]
}

// Implements Variable

func (v variable[T]) Get(s Store) T {
	return v.variable.Get(s).(T)
}

func (v variable[T]) Set(data T) Change {
	return v.variable.Set(data)
}

// Implements fmt.Stringer

func (v variable[T]) String() string {
	var x T
	return reflect.TypeOf(x).String()
}

// Implements json.Marshaler

func (v variable[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

var varFunc func(Store, any) Variable[any]
