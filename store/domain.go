package store

import (
	"github.com/nextmv-io/sdk/model"
)

/*
NewDomain creates a Domain of integers and stores it in a Store.

	s := store.New()
	d1 := store.NewDomain(s, model.Range(1, 10)) // 1 through 10
	d2 := store.NewDomain( // 1 through 10 and 20 through 29
		s,
		model.Range(1, 10),
		model.Range(20, 29),
	)
*/
func NewDomain(s Store, ranges ...model.Range) Domain {
	return domainProxy{domain: NewVar(s, model.NewDomain(ranges...))}
}

/*
NewSingleton creates a Domain containing one integer value and stores it in a
Store.

    s := store.New()
    fortyTwo := store.NewSingleton(s, 42)
*/
func NewSingleton(s Store, value int) Domain {
	return domainProxy{domain: NewVar(s, model.NewSingleton(value))}
}

/*
NewMultiple creates a Domain containing multiple integer values and stores it
in a Store.

	s := store.New()
	even := store.NewMultiple(2, 4, 6, 8)
*/
func NewMultiple(s Store, values ...int) Domain {
	return domainProxy{domain: NewVar(s, model.NewMultiple(values...))}
}

type domainProxy struct {
	domain Variable[model.Domain]
}

// Implements store.Domain.

func (d domainProxy) Add(values ...int) Change {
	return func(s Store) {
		d.domain.Set(d.Domain(s).Add(values...))(s)
	}
}

func (d domainProxy) AtLeast(value int) Change {
	return func(s Store) {
		d.domain.Set(d.Domain(s).AtLeast(value))(s)
	}
}

func (d domainProxy) AtMost(value int) Change {
	return func(s Store) {
		d.domain.Set(d.Domain(s).AtMost(value))(s)
	}
}

func (d domainProxy) Cmp(s Store, d2 Domain) int {
	return d.Domain(s).Cmp(d2.Domain(s))
}

func (d domainProxy) Contains(s Store, value int) bool {
	return d.Domain(s).Contains(value)
}

func (d domainProxy) Domain(s Store) model.Domain {
	return d.domain.Get(s)
}

func (d domainProxy) Empty(s Store) bool {
	return d.Domain(s).Empty()
}

func (d domainProxy) Len(s Store) int {
	return d.Domain(s).Len()
}

func (d domainProxy) Max(s Store) (int, bool) {
	return d.Domain(s).Max()
}

func (d domainProxy) Min(s Store) (int, bool) {
	return d.Domain(s).Min()
}

func (d domainProxy) Remove(values ...int) Change {
	return func(s Store) {
		d.domain.Set(d.Domain(s).Remove(values...))(s)
	}
}

func (d domainProxy) Slice(s Store) []int {
	return d.Domain(s).Slice()
}

func (d domainProxy) Value(s Store) (int, bool) {
	return d.Domain(s).Value()
}
