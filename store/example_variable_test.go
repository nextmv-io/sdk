package store_test

import (
	"encoding/json"
	"fmt"

	"github.com/nextmv-io/sdk/store"
)

func ExampleVar_get() {
	s := store.New()
	x := store.NewVar(s, 10)
	fmt.Println(x.Get(s))
	// Output:
	// 10
}

func ExampleVar_set() {
	s := store.New()
	x := store.NewVar(s, 10)
	fmt.Println(x.Get(s))
	s1 := s.Apply(x.Set(15))
	fmt.Println(x.Get(s1))
	// Output:
	// 10
	// 15
}

// Declaring a new variable adds the variable to the store.
func ExampleNewVar() {
	s := store.New()
	x := store.NewVar(s, 10)
	s = s.Apply(x.Set(15))
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// [15]
}
