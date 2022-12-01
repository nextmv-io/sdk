// package main holds the implementation of the cloud-routing template.
package main

import (
	"log"

	"github.com/nextmv-io/sdk/run"
)

func main() {
	err := run.Run(solver)
	if err != nil {
		log.Fatal(err)
	}
}
