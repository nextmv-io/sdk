package context

import "github.com/nextmv-io/sdk/hop/model"

// Bounder maps a context to monotonically tightening bounds.
type Bounder func(Context) model.Bounds
