package store

import (
	"encoding/json"
	"reflect"
)

// Variable stored in a Store.
type Variable[T any] interface {
	/*
		Get the current value of the Variable in the Store.

			s := store.New()
			x := store.NewVar(s, 10)
			s = s.Format(func(s store.Store) any {
				return map[string]int{"x": x.Get(s)}
			})

	*/
	Get(Store) T

	/*
		Set a new value on the Variable.

			s := store.New()
			x := store.NewVar(s, 10)
			s = s.Apply(x.Set(15))
	*/
	Set(T) Change
}

/*
NewVar stores a new Variable in a Store.

	s := store.New()
	x := store.NewVar(s, 10) // x is stored in s.
*/
func NewVar[T any](s Store, data T) Variable[T] {
	return variable[T]{variable: newVarFunc(s, data)}
}

type variable[T any] struct {
	variable Variable[any]
}

// Implements Variable.

func (v variable[T]) Get(s Store) T {
	return v.variable.Get(s).(T)
}

func (v variable[T]) Set(data T) Change {
	return v.variable.Set(data)
}

// Implements fmt.Stringer.

func (v variable[T]) String() string {
	var x T
	return reflect.TypeOf(x).String()
}

// Implements json.Marshaler.

func (v variable[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

var newVarFunc func(Store, any) Variable[any]
