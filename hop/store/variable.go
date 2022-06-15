package store

import (
	"encoding/json"
	"reflect"

	"github.com/nextmv-io/sdk/hop/store/types"
)

// Var stores a new Variable in a Store.
func Var[T any](s types.Store, data T) types.Variable[T] {
	return variable[T]{variable: varFunc(s, data)}
}

type variable[T any] struct {
	variable types.Variable[any]
}

// Implements types.Variable.

func (v variable[T]) Get(s types.Store) T {
	return v.variable.Get(s).(T)
}

func (v variable[T]) Set(data T) types.Change {
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

var varFunc func(types.Store, any) types.Variable[any]
