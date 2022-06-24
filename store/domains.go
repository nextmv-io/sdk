package store

import (
	"github.com/nextmv-io/sdk/model"
)

/*
NewDomains creates a sequence of Domains and stores the sequence in a Store.

	s := store.New()
	domains := store.NewDomains( // [1 to 10, 42, odds, evens]
		s,
		model.Domain(model.Range(1, 10)),
		model.Singleton(42),
		model.Multiple(1, 3, 5, 7),
		model.Multiple(2, 4, 6, 8),
	)
*/
func NewDomains(s Store, domains ...model.Domain) Domains {
	return domainsProxy{domains: NewVar(s, model.NewDomains(domains...))}
}

/*
Repeat a Domain n times and store the sequence in a Store.

	s := store.New()
	domains := store.Repeat(s, 3, model.Domain(model.Range(1, 10)))
*/
func Repeat(s Store, n int, domain model.Domain) Domains {
	return domainsProxy{domains: NewVar(s, model.Repeat(n, domain))}
}

type domainsProxy struct {
	domains Variable[model.Domains]
}

// Implements store.Domains.

func (d domainsProxy) Add(i int, v ...int) Change {
	return func(s Store) {
		d.domains.Set(d.Domains(s).Add(i, v...))(s)
	}
}

func (d domainsProxy) Assign(i int, v int) Change {
	return func(s Store) {
		d.domains.Set(d.Domains(s).Assign(i, v))(s)
	}
}

func (d domainsProxy) AtLeast(i int, v int) Change {
	return func(s Store) {
		d.domains.Set(d.Domains(s).AtLeast(i, v))(s)
	}
}

func (d domainsProxy) AtMost(i int, v int) Change {
	return func(s Store) {
		d.domains.Set(d.Domains(s).AtMost(i, v))(s)
	}
}

func (d domainsProxy) Cmp(s Store, d2 Domains) int {
	return d.Domains(s).Cmp(d2.Domains(s))
}

func (d domainsProxy) Domain(s Store, i int) model.Domain {
	return d.Domains(s).Domain(i)
}

func (d domainsProxy) Domains(s Store) model.Domains {
	return d.domains.Get(s)
}

func (d domainsProxy) Empty(s Store) bool {
	return d.Domains(s).Empty()
}

func (d domainsProxy) Len(s Store) int {
	return d.Domains(s).Len()
}

func (d domainsProxy) Remove(i int, v ...int) Change {
	return func(s Store) {
		d.domains.Set(d.domains.Get(s).Remove(i, v...))(s)
	}
}

func (d domainsProxy) Singleton(s Store) bool {
	return d.Domains(s).Singleton()
}

func (d domainsProxy) Slices(s Store) [][]int {
	return d.Domains(s).Slices()
}

func (d domainsProxy) Values(s Store) ([]int, bool) {
	return d.Domains(s).Values()
}

/* Domain selectors */

func (d domainsProxy) First(s Store) (int, bool) {
	return d.Domains(s).First()
}

func (d domainsProxy) Largest(s Store) (int, bool) {
	return d.Domains(s).Largest()
}

func (d domainsProxy) Last(s Store) (int, bool) {
	return d.Domains(s).Last()
}

func (d domainsProxy) Maximum(s Store) (int, bool) {
	return d.Domains(s).Maximum()
}

func (d domainsProxy) Minimum(s Store) (int, bool) {
	return d.Domains(s).Minimum()
}

func (d domainsProxy) Smallest(s Store) (int, bool) {
	return d.Domains(s).Smallest()
}
