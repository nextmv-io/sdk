// Package events provides a simple event system for the ALNS.
package events

// BaseEvent is a base event type that can be used to implement events.
type BaseEvent[T any] struct {
	handlers []Handler[T]
}

// Register adds an event handler for this event.
func (e *BaseEvent[T]) Register(handler Handler[T]) {
	e.handlers = append(e.handlers, handler)
}

// Trigger sends out an event with the payload.
func (e *BaseEvent[T]) Trigger(payload T) {
	for _, handler := range e.handlers {
		handler(payload)
	}
}

// Handler is a function that handles an event.
type Handler[T any] func(payload T)
