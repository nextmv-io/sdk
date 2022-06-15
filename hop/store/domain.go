package store

import (
	"github.com/nextmv-io/sdk/hop/model"
	modeltypes "github.com/nextmv-io/sdk/hop/model/types"
	"github.com/nextmv-io/sdk/hop/store/types"
)

// Domain creates a domain of integers.
func Domain(s types.Store, ranges ...modeltypes.Range) types.Domain {
	return domainProxy{domain: Var(s, model.Domain(ranges...))}
}

// Singleton creates a domain containing one integer value.
func Singleton(s types.Store, value int) types.Domain {
	return domainProxy{domain: Var(s, model.Singleton(value))}
}

// Multiple creates a domain containing multiple integer values.
func Multiple(s types.Store, values ...int) types.Domain {
	return domainProxy{domain: Var(s, model.Multiple(values...))}
}

type domainProxy struct {
	domain types.Variable[modeltypes.Domain]
}

// Implements types.Domain.

func (d domainProxy) Add(values ...int) types.Change {
	return func(s types.Store) {
		d.domain.Set(d.Domain(s).Add(values...))(s)
	}
}

func (d domainProxy) AtLeast(value int) types.Change {
	return func(s types.Store) {
		d.domain.Set(d.Domain(s).AtLeast(value))(s)
	}
}

func (d domainProxy) AtMost(value int) types.Change {
	return func(s types.Store) {
		d.domain.Set(d.Domain(s).AtMost(value))(s)
	}
}

func (d domainProxy) Cmp(s types.Store, d2 types.Domain) int {
	return d.Domain(s).Cmp(d2.Domain(s))
}

func (d domainProxy) Contains(s types.Store, value int) bool {
	return d.Domain(s).Contains(value)
}

func (d domainProxy) Domain(s types.Store) modeltypes.Domain {
	return d.domain.Get(s)
}

func (d domainProxy) Empty(s types.Store) bool {
	return d.Domain(s).Empty()
}

func (d domainProxy) Len(s types.Store) int {
	return d.Domain(s).Len()
}

func (d domainProxy) Max(s types.Store) (int, bool) {
	return d.Domain(s).Max()
}

func (d domainProxy) Min(s types.Store) (int, bool) {
	return d.Domain(s).Min()
}

func (d domainProxy) Remove(values ...int) types.Change {
	return func(s types.Store) {
		d.domain.Set(d.Domain(s).Remove(values...))(s)
	}
}

func (d domainProxy) Slice(s types.Store) []int {
	return d.Domain(s).Slice()
}

func (d domainProxy) Value(s types.Store) (int, bool) {
	return d.Domain(s).Value()
}
