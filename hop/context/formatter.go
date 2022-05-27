package context

// Formatter maps a context to any type with a JSON representation.
type Formatter func(Context) any
