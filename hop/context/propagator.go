package context

type Propagator func(Context) []func(Context)
