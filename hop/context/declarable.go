package context

import "reflect"

type Declarable interface {
	Number() uint64
	Type() reflect.Type
}
