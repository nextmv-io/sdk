package store

import (
	"github.com/nextmv-io/sdk/model"
)

// A Domain of integers.
//
// Deprecated: This package is deprecated and will be removed in the future.
type Domain interface {
	/*
		Add values to a Domain.

			s1 := store.New()
			d := store.Multiple(s1, 1, 3, 5)
			s2 := s1.Apply(d.Add(2, 4))

			d.Domain(s1) // {1, 3, 5}}
			d.Domain(s2) // [1, 5]]

		Deprecated: This package is deprecated and will be removed in the future.
	*/
	Add(...int) Change

	/*
		AtLeast updates the Domain to the sub-Domain of at least some value.

			s1 := store.New()
			d := store.NewDomain(s1, model.NewRange(1, 10), model.NewRange(101, 110))
			s2 := s1.Apply(d.AtLeast(50))

			d.Domain(s1) // {[1, 10], [101, 110]}
			d.Domain(s2) // [101, 110]

		Deprecated: This package is deprecated and will be removed in the future.
	*/
	AtLeast(int) Change

	/*
		AtMost updates the Domain to the sub-Domain of at most some value.

			s1 := store.New()
			d := store.NewDomain(s1, model.NewRange(1, 10), model.NewRange(101, 110))
			s2 := s1.Apply(d.AtMost(50))

			d.Domain(s1) // {[1, 10], [101, 110]}
			d.Domain(s2) // [1, 10]

		Deprecated: This package is deprecated and will be removed in the future.
	*/
	AtMost(int) Change

	/*
		Cmp lexically compares two integer Domains. It returns a negative value
		if the receiver is less, 0 if they are equal, and a positive value if
		the receiver Domain is greater.

			s := store.New()
			d1 := store.NewDomain(s, model.NewRange(1, 5), model.NewRange(8, 10))
			d2 := store.Multiple(s, -1, 1)
			d1.Cmp(s, d2) // > 0

		Deprecated: This package is deprecated and will be removed in the future.
	*/
	Cmp(Store, Domain) int

	/*
		Contains returns true if a Domain contains a given value.

			s := store.New()
			d := store.NewDomain(s, model.NewRange(1, 10))
			d.Contains(s, 5)  // true
			d.Contains(s, 15) // false

		Deprecated: This package is deprecated and will be removed in the future.
	*/
	Contains(Store, int) bool

	/*
		Domain returns a Domain unattached to a Store.

			s := store.New()
			d := store.NewDomain(s, model.NewRange(1, 10))
			d.Domain(s) // model.NewDomain(model.NewRange(1, 10))

		Deprecated: This package is deprecated and will be removed in the future.
	*/
	Domain(Store) model.Domain

	/*
		Empty is true if a Domain is empty for a Store.

			s := store.New()
			d1 := store.NewDomain(s)
			d2 := store.Singleton(s, 42)
			d1.Empty(s) // true
			d2.Empty(s) // false

		Deprecated: This package is deprecated and will be removed in the future.
	*/
	Empty(Store) bool

	/*
		Len of a Domain, counting all values within ranges.

			s := store.New()
			d := store.NewDomain(s, model.NewRange(1, 10), model.NewRange(-5, -1))
			d.Len(s) // 15

		Deprecated: This package is deprecated and will be removed in the future.
	*/
	Len(Store) int

	/*
		Max of a Domain and a boolean indicating it is non-empty.

			s := store.New()
			d1 := store.NewDomain(s)
			d2 := store.NewDomain(s, model.NewRange(1, 10), model.NewRange(-5, -1))
			d1.Max(s) // returns (_, false)
			d2.Max(s) // returns (10, true)

		Deprecated: This package is deprecated and will be removed in the future.
	*/
	Max(Store) (int, bool)

	/*
		Min of a Domain and a boolean indicating it is non-empty.

			s := store.New()
			d1 := store.NewDomain(s)
			d2 := store.NewDomain(s, model.NewRange(1, 10), model.NewRange(-5, -1))
			d1.Min(s) // returns (_, false)
			d2.Min(s) // returns (-5, true)

		Deprecated: This package is deprecated and will be removed in the future.
	*/
	Min(Store) (int, bool)

	/*
		Remove values from a Domain.

			s1 := store.New()
			d := store.NewDomain(s1, model.NewRange(1, 5))
			s2 := s1.Apply(d.Remove([]int{2, 4}))

			d.Domain(s1) // [1, 5]
			d.Domain(s2) // {1, 3, 5}

		Deprecated: This package is deprecated and will be removed in the future.
	*/
	Remove([]int) Change

	/*
		Slice representation of a Domain.

			s := store.New()
			d := store.NewDomain(s, model.NewRange(1, 5))
			d.Slice(s) // [1, 2, 3, 4, 5]

		Deprecated: This package is deprecated and will be removed in the future.
	*/
	Slice(Store) []int

	/*
		Value returns an int and true if a Domain is Singleton.

			s := store.New()
			d1 := store.NewDomain(s)
			d2 := store.Singleton(s, 42)
			d3 := store.Multiple(s, 1, 3, 5)
			d1.Value(s) // returns (0, false)
			d2.Value(s) // returns (42, true)
			d3.Value(s) // returns (0, false)

		Deprecated: This package is deprecated and will be removed in the future.
	*/
	Value(Store) (int, bool)
}

/*
NewDomain creates a Domain of integers and stores it in a Store.

	s := store.New()
	d1 := store.NewDomain(s, model.NewRange(1, 10)) // 1 through 10
	d2 := store.NewDomain( // 1 through 10 and 20 through 29
		s,
		model.NewRange(1, 10),
		model.NewRange(20, 29),
	)

Deprecated: This package is deprecated and will be removed in the future.
*/
func NewDomain(s Store, ranges ...model.Range) Domain {
	return domainProxy{domain: NewVar(s, model.NewDomain(ranges...))}
}

/*
Singleton creates a Domain containing one integer value and stores it in a
Store.

	s := store.New()
	fortyTwo := store.Singleton(s, 42)

Deprecated: This package is deprecated and will be removed in the future.
*/
func Singleton(s Store, value int) Domain {
	return domainProxy{domain: NewVar(s, model.Singleton(value))}
}

/*
Multiple creates a Domain containing multiple integer values and stores it
in a Store.

	s := store.New()
	even := store.Multiple(s, 2, 4, 6, 8)

Deprecated: This package is deprecated and will be removed in the future.
*/
func Multiple(s Store, values ...int) Domain {
	return domainProxy{domain: NewVar(s, model.Multiple(values...))}
}

type domainProxy struct {
	domain Var[model.Domain]
}

// Implements store.Domain.
//
// Deprecated: This package is deprecated and will be removed in the future.

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

func (d domainProxy) Remove(values []int) Change {
	return func(s Store) {
		d.domain.Set(d.Domain(s).Remove(values))(s)
	}
}

func (d domainProxy) Slice(s Store) []int {
	return d.Domain(s).Slice()
}

func (d domainProxy) Value(s Store) (int, bool) {
	return d.Domain(s).Value()
}
