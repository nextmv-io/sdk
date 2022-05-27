package context

// Valuer maps a context to an integer value.
type Valuer func(Context) int
