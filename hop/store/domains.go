package store

import (
	"github.com/nextmv-io/sdk/hop/model"
	modeltypes "github.com/nextmv-io/sdk/hop/model/types"
	"github.com/nextmv-io/sdk/hop/store/types"
)

// Domains creates a sequence of domains.
func Domains(s types.Store, domains ...modeltypes.Domain) types.Domains {
	return domainsProxy{domains: Var(s, model.Domains(domains...))}
}

// Repeat repeats a domain n times.
func Repeat(s types.Store, n int, domain modeltypes.Domain) types.Domains {
	return domainsProxy{domains: Var(s, model.Repeat(n, domain))}
}

type domainsProxy struct {
	domains types.Variable[modeltypes.Domains]
}

// Implements types.Domains.

func (d domainsProxy) Add(i int, v ...int) types.Change {
	return func(s types.Store) {
		d.domains.Set(d.Domains(s).Add(i, v...))(s)
	}
}

func (d domainsProxy) Assign(i int, v int) types.Change {
	return func(s types.Store) {
		d.domains.Set(d.Domains(s).Assign(i, v))(s)
	}
}

func (d domainsProxy) AtLeast(i int, v int) types.Change {
	return func(s types.Store) {
		d.domains.Set(d.Domains(s).AtLeast(i, v))(s)
	}
}

func (d domainsProxy) AtMost(i int, v int) types.Change {
	return func(s types.Store) {
		d.domains.Set(d.Domains(s).AtMost(i, v))(s)
	}
}

func (d domainsProxy) Cmp(s types.Store, d2 types.Domains) int {
	return d.Domains(s).Cmp(d2.Domains(s))
}

func (d domainsProxy) Domain(s types.Store, i int) modeltypes.Domain {
	return d.Domains(s).Domain(i)
}

func (d domainsProxy) Domains(s types.Store) modeltypes.Domains {
	return d.domains.Get(s)
}

func (d domainsProxy) Empty(s types.Store) bool {
	return d.Domains(s).Empty()
}

func (d domainsProxy) Len(s types.Store) int {
	return d.Domains(s).Len()
}

func (d domainsProxy) Remove(i int, v ...int) types.Change {
	return func(s types.Store) {
		d.domains.Set(d.domains.Get(s).Remove(i, v...))(s)
	}
}

func (d domainsProxy) Singleton(s types.Store) bool {
	return d.Domains(s).Singleton()
}

func (d domainsProxy) Slices(s types.Store) [][]int {
	return d.Domains(s).Slices()
}

func (d domainsProxy) Values(s types.Store) ([]int, bool) {
	return d.Domains(s).Values()
}

/* Domain selectors */

func (d domainsProxy) First(s types.Store) (int, bool) {
	return d.Domains(s).First()
}

func (d domainsProxy) Largest(s types.Store) (int, bool) {
	return d.Domains(s).Largest()
}

func (d domainsProxy) Last(s types.Store) (int, bool) {
	return d.Domains(s).Last()
}

func (d domainsProxy) Maximum(s types.Store) (int, bool) {
	return d.Domains(s).Maximum()
}

func (d domainsProxy) Minimum(s types.Store) (int, bool) {
	return d.Domains(s).Minimum()
}

func (d domainsProxy) Smallest(s types.Store) (int, bool) {
	return d.Domains(s).Smallest()
}
