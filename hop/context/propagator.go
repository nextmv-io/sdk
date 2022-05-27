package context

// Propagator propagates a context into multiple ones.
type Propagator func(Context) []func(Context)
