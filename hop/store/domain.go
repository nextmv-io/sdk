package store

import (
	"github.com/nextmv-io/sdk/hop/model"
	modeltypes "github.com/nextmv-io/sdk/hop/model/types"
	"github.com/nextmv-io/sdk/hop/store/types"
)

/*
Domain creates a domain of integers associated with a store.

	s := store.New()
	d1 := store.Domain(s, model.Range(1, 10)) // 1 through 10
	d2 := store.Domain( // 1 through 10 and 20 through 29
		s,
		model.Range(1, 10),
		model.Range(20, 29),
	)
*/
func Domain(s types.Store, ranges ...modeltypes.Range) types.Domain {
	return domainProxy{domain: Var(s, model.Domain(ranges...))}
}

/*
Singleton creates a domain containing one integer value in a store.

	s := store.New()
	fortyTwo := store.Singleton(s, 42)
*/
func Singleton(s types.Store, value int) types.Domain {
	return domainProxy{domain: Var(s, model.Singleton(value))}
}

/*
Multiple creates a domain containing multiple integer values.

	s := store.New()
	even := store.Multiple(2, 4, 6, 8)
*/
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
