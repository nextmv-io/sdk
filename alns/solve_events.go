package alns

import (
	"github.com/nextmv-io/sdk/events"
)

// SolveEvents is a struct that contains events that are fired during the solve.
type SolveEvents[T Solution[T]] struct {
	OperatorExecuting *events.BaseEvent1[SolveInformation[T]]
	OperatorExecuted  *events.BaseEvent1[SolveInformation[T]]
	ImprovementFound  *events.BaseEvent1[SolveInformation[T]]
	Iterating         *events.BaseEvent1[SolveInformation[T]]
	Iterated          *events.BaseEvent1[SolveInformation[T]]
	ContextDone       *events.BaseEvent1[SolveInformation[T]]
	Start             *events.BaseEvent1[SolveInformation[T]]
	Reset             *events.BaseEvent2[Solution[T], SolveInformation[T]]
	Done              *events.BaseEvent1[SolveInformation[T]]
}

// NewSolveEvents creates a new instance of SolveEvents.
func NewSolveEvents[T Solution[T]]() SolveEvents[T] {
	return SolveEvents[T]{
		OperatorExecuting: &events.BaseEvent1[SolveInformation[T]]{},
		OperatorExecuted:  &events.BaseEvent1[SolveInformation[T]]{},
		ImprovementFound:  &events.BaseEvent1[SolveInformation[T]]{},
		Iterating:         &events.BaseEvent1[SolveInformation[T]]{},
		Iterated:          &events.BaseEvent1[SolveInformation[T]]{},
		ContextDone:       &events.BaseEvent1[SolveInformation[T]]{},
		Start:             &events.BaseEvent1[SolveInformation[T]]{},
		Reset:             &events.BaseEvent2[Solution[T], SolveInformation[T]]{},
		Done:              &events.BaseEvent1[SolveInformation[T]]{},
	}
}
