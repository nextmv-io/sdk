package context

import "reflect"

// Declarable represents data that is declared on a context. The declared data
// has a unique number in the context.
type Declarable interface {
	Number() uint64
	Type() reflect.Type
}
