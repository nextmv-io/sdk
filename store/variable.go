package store

import (
	"encoding/json"
	"reflect"

	"github.com/nextmv-io/sdk/connect"
)

// Var is a variable stored in a Store.
//
// Deprecated: This package is deprecated and will be removed in the future.
type Var[T any] interface {
	/*
		Get the current value of the variable in the Store.

			s := store.New()
			x := store.NewVar(s, 10)
			s = s.Format(func(s store.Store) any {
				return map[string]int{"x": x.Get(s)}
			})

		Deprecated: This package is deprecated and will be removed in the future.
	*/
	Get(Store) T

	/*
		Set a new value on the variable.

			s := store.New()
			x := store.NewVar(s, 10)
			s = s.Apply(x.Set(15))

		Deprecated: This package is deprecated and will be removed in the future.
	*/
	Set(T) Change
}

/*
NewVar stores a new variable in a Store.

	s := store.New()
	x := store.NewVar(s, 10) // x is stored in s.

Deprecated: This package is deprecated and will be removed in the future.
*/
func NewVar[T any](s Store, data T) Var[T] {
	connect.Connect(con, &newVarFunc)
	return variable[T]{variable: newVarFunc(s, data)}
}

type variable[T any] struct {
	variable Var[any]
}

// Implements Var.
//
// Deprecated: This package is deprecated and will be removed in the future.

func (v variable[T]) Get(s Store) T {
	if value := v.variable.Get(s); value != nil {
		return value.(T)
	}

	// zero-value of T
	var value T
	return value
}

func (v variable[T]) Set(data T) Change {
	return v.variable.Set(data)
}

// Implements fmt.Stringer.
//
// Deprecated: This package is deprecated and will be removed in the future.

func (v variable[T]) String() string {
	var x T
	return reflect.TypeOf(x).String()
}

// Implements json.Marshaler.
//
// Deprecated: This package is deprecated and will be removed in the future.

func (v variable[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

var newVarFunc func(Store, any) Var[any]
