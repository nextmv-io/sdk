package main

import (
	"fmt"

	"github.com/nextmv-io/sdk"
	"github.com/nextmv-io/sdk/hop/store"
)

func main() {
	s := store.New()
	version := store.Var(s, sdk.VERSION)
	fmt.Println("Hello Hop", version.Get(s))
}
