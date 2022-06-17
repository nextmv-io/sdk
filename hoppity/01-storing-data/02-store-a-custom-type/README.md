# Hopping About: Store a Custom Type

A Hop store can manage any concrete type, including custom structs. Let's define
a `bunny` type with a few fields and a `String` method.

```go
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

```

```go
	s := store.New()

	peter := store.Var(s, bunny{
		name:       "Peter Rabbit",
		fluffiness: 0.52,
		activities: []string{
			"stealing and eating vegetables",
			"losing his jacket and shoes",
		},
	})
```

If we retrieve the value of `peter` from our store and print it, we should get a
the results of the `bunny.String` method.

```go
	fmt.Println(peter.Get(s))
```

## Source

```go
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
```

## Exercises

* Add another bunny to the store and call its `String` method directly.
* Add another method to the `bunny` type. Retrieve `peter` from the store and
  call this new method.

