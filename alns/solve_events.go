package alns

import (
	"github.com/nextmv-io/sdk/events"
)

// SolveEvents is a struct that contains events that are fired during the solve.
type SolveEvents[T Solution[T]] struct {
	OperatorExecuting *events.BaseEvent[SolveInformation[T]]
	OperatorExecuted  *events.BaseEvent[SolveInformation[T]]
	ImprovementFound  *events.BaseEvent[SolveInformation[T]]
	Iterating         *events.BaseEvent[SolveInformation[T]]
	Iterated          *events.BaseEvent[SolveInformation[T]]
	ContextDone       *events.BaseEvent[SolveInformation[T]]
	Start             *events.BaseEvent[SolveInformation[T]]
	Reset             *events.BaseEvent[SolveInformation[T]]
	Done              *events.BaseEvent[SolveInformation[T]]
}

// NewSolveEvents creates a new instance of SolveEvents.
func NewSolveEvents[T Solution[T]]() SolveEvents[T] {
	return SolveEvents[T]{
		OperatorExecuting: &events.BaseEvent[SolveInformation[T]]{},
		OperatorExecuted:  &events.BaseEvent[SolveInformation[T]]{},
		ImprovementFound:  &events.BaseEvent[SolveInformation[T]]{},
		Iterating:         &events.BaseEvent[SolveInformation[T]]{},
		Iterated:          &events.BaseEvent[SolveInformation[T]]{},
		ContextDone:       &events.BaseEvent[SolveInformation[T]]{},
		Start:             &events.BaseEvent[SolveInformation[T]]{},
		Reset:             &events.BaseEvent[SolveInformation[T]]{},
		Done:              &events.BaseEvent[SolveInformation[T]]{},
	}
}
