# Hopping About / Working with Stores / Creating Stores

[home](../README.md)

To start, let's import the `hop/store` subpackage of Nextmv's SDK and create a
`main` stub.

```go
package main

import (
	"fmt"

	"github.com/nextmv-io/sdk/hop/store"
)

func main() {
    // code goes here
}
```

A store is lexical scope, similar to the lexical scope in most programming 
languages. It can hold variable declarations, variable assignments, and logic.

```go
s := store.New()
```

You can add any variables to a store. The store will manage them for you.

```go
x := store.Var(s, 42)                    // x is a stored int
y := store.Var(s, []float64{3.14, 2.72}) // y is a stored []float64
```

You can retrieve typed variable values from the store with their `Get` methods. 
The store knows which type they are so you don't have to think about it.

```go
fmt.Println(
	x.Get(s)*10, // x.Get(s) returns an int
	y.Get(s)[0], // y.Get(s) returns a []float64
)
```

## Source

Let's put this together and try it. Save the [full source][source] to a file 
called `main.go` inside its own directory under `hoppity/`. For example, 
`working-with-stores/create-stores/main.go` works nicely. Run it using
`go run -trimpath`. You should see output like this.

```bash
hoppity$ go run -trimpath storing-data/create-a-store/main.go
420 3.14
```

In future examples, we'll leave out those steps. You can save the files wherever
you like under the `hoppity/` folder, so long as there is only one `func main` 
in any subfolder.

## Exercises

* Add more variables of different types to the store. Print their values.
* Create a seconds store and add variables to it.
* What happens when you retrieve a value from the wrong store?

[source]: creating-stores/main.go
