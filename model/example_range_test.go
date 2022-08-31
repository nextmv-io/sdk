package model_test

import (
	"fmt"

	"github.com/nextmv-io/sdk/model"
)

func ExampleNewRange() {
	r1 := model.NewRange(0, 1)
	fmt.Println(r1)
	// Output:
	// {0 1}
}
