package store

import (
	"github.com/nextmv-io/sdk/model"
)

// Domains of integers.
type Domains interface {
	/*
		Add values to a Domain by index.

			s1 := store.New()
			d := store.Repeat(s1, 1, model.Singleton(42)) // [42, 42, 42]
			s2 := s1.Apply(d.Add(1, 41, 43))              // [42, [41,43], 42]
	*/
	Add(int, ...int) Change

	/*
		Assign a Singleton value to a Domain by index.

			s1 := store.New()
			d := store.Repeat(s1, 3, model.Singleton(42)) // [42, 42, 42]
			s2 := s1.Apply(d.Assign(0, 10))               // [10, 42, 42]
	*/
	Assign(int, int) Change

	/*
		AtLeast updates the Domain to the sub-Domain of at least some value.

			s1 := store.New()
			d := store.Repeat( // [[1, 100], [1, 100]]
				s1,
				2,
				model.NewDomain(model.NewRange(1, 100)),
			)
			s2 := s1.Apply(d.AtLeast(1, 50)) // [[1, 100], [50, 100]]
	*/
	AtLeast(int, int) Change

	/*
		AtMost updates the Domain to the sub-Domain of at most some value.

			s1 := store.New()
			d := store.Repeat( // [[1, 100], [1, 100]]
				s1,
				2,
				model.NewDomain(model.NewRange(1, 100)),
			)
			s2 := s1.Apply(d.AtMost(1, 50)) // [[1, 100], [1, 50]]
	*/
	AtMost(int, int) Change

	/*
		Cmp lexically compares two sequences of integer Domains.

			s := store.New()
			d1 := store.Repeat(s, 3, model.Singleton(42)) // [42, 42, 42]
			d2 := store.Repeat(s, 2, model.Singleton(43)) // [43, 43]]
			d1.Cmp(s, d2) // < 0
	*/
	Cmp(Store, Domains) int

	/*
		Domain by index.

			s := store.New()
			d := store.NewDomains(
				s,
				model.NewDomain(),
				model.Singleton(42),
			)
			d.Domain(s, 0) // {}
			d.Domain(s, 1) // 42
	*/
	Domain(Store, int) model.Domain

	/*
		Domains in the sequence.

			s := store.New()
			d := store.NewDomains(
				s,
				model.NewDomain(),
				model.Singleton(42),
			)
			d.Domains(s) // [{}, 42}
	*/
	Domains(Store) model.Domains

	/*
		Empty is true if all Domains are empty.

			s := store.New()
			d := store.NewDomains(s, model.NewDomain())
			d.Empty(s) // true
	*/
	Empty(Store) bool

	/*
		Len returns the number of Domains.

			s := store.New()
			d := store.Repeat(s, 5, model.NewDomain())
			d.Len(s) // 5
	*/
	Len(Store) int

	/*
		Remove values from a Domain by index.

			s1 := store.New()
			d := store.NewDomains(s1, model.Multiple(42, 13)) // {13, 42}
			s2 := s1.Apply(d.Remove(13)) // {42}
	*/
	Remove(int, ...int) Change

	/*
		Singleton is true if all Domains are Singleton.

			s := store.New()
			d := store.Repeat(s, 5, model.Singleton(42))
			d.Singleton(s) // true
	*/
	Singleton(Store) bool

	/*
		Slices converts Domains to a slice of int slices.

			s := store.New()
			d := store.NewDomains(s, model.NewDomain(), model.Multiple(1, 3))
			d.Slices(s) // [[], [1, 2, 3]]
	*/
	Slices(Store) [][]int

	/*
		Values returns the values of a sequence of Singleton Domains.

			s1 := store.New()
			d := store.Repeat(s1, 3, model.Singleton(42))
			s2 := store.Apply(d.Add(0, 41))
			d.Values(s1) // ([42, 42, 42], true)
			d.Values(s2) // (_, false)
	*/
	Values(Store) ([]int, bool)

	/* Domain selectors */

	// First returns the first Domain index with length above 1.
	First(Store) (int, bool)
	// Largest returns the index of the largest Domain with length above 1 by
	// number of elements.
	Largest(Store) (int, bool)
	// Last returns the last Domain index with length above 1.
	Last(Store) (int, bool)
	// Maximum returns the index of the Domain containing the maximum value
	// with length above 1.
	Maximum(Store) (int, bool)
	// Minimum returns the index of the Domain containing the minimum value
	// with length above 1.
	Minimum(Store) (int, bool)
	// Smallest returns the index of the smallest Domain with length above 1 by
	// number of elements.
	Smallest(Store) (int, bool)
}

/*
NewDomains creates a sequence of Domains and stores the sequence in a Store.

	s := store.New()
	domains := store.NewDomains( // [1 to 10, 42, odds, evens]
		s,
		model.Domain(model.NewRange(1, 10)),
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
