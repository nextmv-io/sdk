package model

import (
	"math/bits"

	"github.com/nextmv-io/sdk/connect"
)

// Constants for integer bounds.
const (
	// MaxInt is the maximum value for an integer.
	MaxInt int = (1<<bits.UintSize)/2 - 1
	// MinInt is the minimum value for an integer.
	MinInt = (1 << bits.UintSize) / -2
)

// A Domain of integers.
//
// Deprecated: This package is deprecated and will be removed in the future.
type Domain interface {
	// Add values to a domain.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Add(v ...int) Domain
	// AtLeast updates the domain to the subdomain of at least some value.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	AtLeast(int) Domain
	// AtMost updates the domain to the subdomain of at most some value.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	AtMost(int) Domain
	// Cmp lexically compares two integer domains. It returns a negative value
	// if the receiver is less, 0 if they are equal, and a positive value if
	// the receiver domain is greater.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Cmp(Domain) int
	// Contains returns true if a domain contains a given value.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Contains(int) bool
	// Empty is true if a domain is empty.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Empty() bool
	// Iterator over a domain.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Iterator() Iterator
	// Len of a domain, counting all values within ranges.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Len() int
	// Max of a domain and a boolean indicating it is nonempty.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Max() (int, bool)
	// Min of a domain and a boolean indicating it is nonempty.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Min() (int, bool)
	// Overlaps returns true if a domain overlaps another domain.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Overlaps(other Domain) bool
	// Remove values from a domain.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Remove([]int) Domain
	// Slice representation of a domain.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Slice() []int
	// Value returns an int and true if a domain is singleton.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Value() (int, bool)
}

// Domains of integers.
//
// Deprecated: This package is deprecated and will be removed in the future.
type Domains interface {
	// Add values to a domain by index.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Add(int, ...int) Domains
	// Assign a singleton value to a domain by index.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Assign(int, int) Domains
	// AtLeast updates the domain to the subdomain of at least some value.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	AtLeast(int, int) Domains
	// AtMost updates the domain to the subdomain of at most some value.
	AtMost(int, int) Domains
	// Cmp lexically compares two sequences of integer domains.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Cmp(Domains) int
	// Domain by index.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Domain(int) Domain
	// Empty is true if all domains are empty.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Empty() bool
	// Len returns the number of domains.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Len() int
	// Remove values from a domain by index.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Remove(int, []int) Domains
	// Singleton is true if all domains are Singleton.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Singleton() bool
	// Slices convert domains to a slice of int slices.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Slices() [][]int
	// Values returns the values of a sequence of singleton domains/
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Values() ([]int, bool)

	/* Domain selectors */

	// First returns the first domain index with length above 1.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	First() (int, bool)
	// Largest returns the index of the largest domain with length above 1 by
	// number of elements.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Largest() (int, bool)
	// Last returns the last domain index with length above 1.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Last() (int, bool)
	// Maximum returns the index of the domain containing the maximum value with
	// length above 1.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Maximum() (int, bool)
	// Minimum returns the index of the domain containing the minimum value with
	// length above 1.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Minimum() (int, bool)
	// Smallest returns the index of the smallest domain with length above 1 by
	// number of elements.
	//
	// Deprecated: This package is deprecated and will be removed in the future.
	Smallest() (int, bool)
}

// A Range of integers.
//
// Deprecated: This package is deprecated and will be removed in the future.
type Range interface {
	Min() int
	Max() int
}

// An Iterator allows one to iterate over a range or a domain.
//
//	it := model.Domain(model.Range(1, 10)).Iterator()
//	for it.Next() {
//	    fmt.Println(it.Value()) // 1, ..., 10
//	}
//
// Deprecated: This package is deprecated and will be removed in the future.
type Iterator interface {
	Next() bool
	Value() int
}

// NewDomain creates a domain of integers.
//
// Deprecated: This package is deprecated and will be removed in the future.
func NewDomain(ranges ...Range) Domain {
	connect.Connect(con, &newDomainFunc)
	return newDomainFunc(ranges...)
}

// Singleton creates a domain containing one integer value.
//
// Deprecated: This package is deprecated and will be removed in the future.
func Singleton(value int) Domain {
	connect.Connect(con, &singletonFunc)
	return singletonFunc(value)
}

// Multiple creates a domain containing multiple integer values.
//
// Deprecated: This package is deprecated and will be removed in the future.
func Multiple(values ...int) Domain {
	connect.Connect(con, &multipleFunc)
	return multipleFunc(values...)
}

// NewDomains creates a sequence of domains.
//
// Deprecated: This package is deprecated and will be removed in the future.
func NewDomains(domains ...Domain) Domains {
	connect.Connect(con, &newDomainsFunc)
	return newDomainsFunc(domains...)
}

// Repeat a domain n times.
//
// Deprecated: This package is deprecated and will be removed in the future.
func Repeat(n int, d Domain) Domains {
	connect.Connect(con, &repeatFunc)
	return repeatFunc(n, d)
}

// NewRange create a new integer range.
//
// Deprecated: This package is deprecated and will be removed in the future.
func NewRange(min, max int) Range {
	connect.Connect(con, &newRangeFunc)
	return newRangeFunc(min, max)
}

var (
	con            = connect.NewConnector("sdk", "Model")
	newDomainFunc  func(...Range) Domain
	singletonFunc  func(int) Domain
	multipleFunc   func(...int) Domain
	newDomainsFunc func(...Domain) Domains
	repeatFunc     func(int, Domain) Domains
	newRangeFunc   func(int, int) Range
)
