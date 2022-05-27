// © 2019-2022 nextmv.io inc. All rights reserved.
// nextmv.io, inc. CONFIDENTIAL
//
// This file includes unpublished proprietary source code of nextmv.io, inc.
// The copyright notice above does not evidence any actual or intended
// publication of such source code. Disclosure of this source code or any
// related proprietary information is strictly prohibited without the express
// written permission of nextmv.io, inc.

package main

import (
	"encoding/json"
	"os"

	"github.com/nextmv-io/sdk/hop/context"
)

func main() {
	global := context.NewContext()
	x := context.Declare(global, 42)

	enc := json.NewEncoder(os.Stdout)
	err := enc.Encode(global)
	if err != nil {
		panic(err)
	}

	global = global.Apply(x.Set(13))
	global = global.Format(
		func(local context.Context) any {
			return map[string]any{"x": x.Get(local)}
		},
	)
	global = global.Value(x.Get)
	global = global.Check(context.False)

	err = enc.Encode(global)
	if err != nil {
		panic(err)
	}
}
