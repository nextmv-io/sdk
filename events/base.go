// Package events provides a simple event system for the ALNS.
package events

// BaseEvent1 is a base event type that can be used to implement events
// with one payload.
type BaseEvent1[T any] struct {
	handlers []Handler1[T]
}

// Register adds an event handler for this event.
func (e *BaseEvent1[T]) Register(handler Handler1[T]) {
	e.handlers = append(e.handlers, handler)
}

// Trigger sends out an event with the payload.
func (e *BaseEvent1[T]) Trigger(payload T) {
	for _, handler := range e.handlers {
		handler(payload)
	}
}

// Handler1 is a function that handles an event with one payload.
type Handler1[T any] func(payload T)

// BaseEvent2 is a base event type that can be used to implement events
// with two payloads.
type BaseEvent2[S any, T any] struct {
	handlers []Handler2[S, T]
}

// Register adds an event handler for this event.
func (e *BaseEvent2[S, T]) Register(handler Handler2[S, T]) {
	e.handlers = append(e.handlers, handler)
}

// Trigger sends out an event with the payload.
func (e *BaseEvent2[S, T]) Trigger(payload1 S, payload2 T) {
	for _, handler := range e.handlers {
		handler(payload1, payload2)
	}
}

// Handler2 is a function that handles an event with two payloads.
type Handler2[S any, T any] func(payload1 S, payload2 T)
