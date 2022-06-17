package main

import (
	"fmt"

	"github.com/nextmv-io/sdk/hop/store"
)

type bunny struct {
	name       string
	fluffiness float64
	activities []string
}

func (b bunny) String() string {
	fluffy := "fluffy"
	if b.fluffiness < 0.5 {
		fluffy = "not fluffy"
	}
	return fmt.Sprintf("%s is %s and likes %v", b.name, fluffy, b.activities)
}

func main() {
	s := store.New()

	peter := store.Var(s, bunny{
		name:       "Peter Rabbit",
		fluffiness: 0.52,
		activities: []string{
			"stealing and eating vegetables",
			"losing his jacket and shoes",
		},
	})

	fmt.Println(peter.Get(s))
}
