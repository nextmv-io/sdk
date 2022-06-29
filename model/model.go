package model

// A Domain of integers.
type Domain interface {
	// Add values to a domain.
	Add(...int) Domain
	// AtLeast updates the domain to the subdomain of at least some value.
	AtLeast(int) Domain
	// AtMost updates the domain to the subdomain of at most some value.
	AtMost(int) Domain
	// Cmp lexically compares two integer domains. It returns a negative value
	// if the receiver is less, 0 if they are equal, and a positive value if
	// the receiver domain is greater.
	Cmp(Domain) int
	// Contains returns true if a domain contains a given value.
	Contains(int) bool
	// Empty is true if a domain is empty.
	Empty() bool
	// Len of a domain, counting all values within ranges.
	Len() int
	// Max of a domain and a boolean indicating it is nonempty.
	Max() (int, bool)
	// Min of a domain and a boolean indicating it is nonempty.
	Min() (int, bool)
	// Remove values from a domain.
	Remove(...int) Domain
	// Slice representation of a domain.
	Slice() []int
	// Value returns an int and true if a domain is singleton.
	Value() (int, bool)
}

// Domains of integers.
type Domains interface {
	// Add values to a domain by index.
	Add(int, ...int) Domains
	// Assign a singleton value to a domain by index.
	Assign(int, int) Domains
	// AtLeast updates the domain to the subdomain of at least some value.
	AtLeast(int, int) Domains
	// AtMost updates the domain to the subdomain of at most some value.
	AtMost(int, int) Domains
	// Cmp lexically compares two sequences of integer domains.
	Cmp(Domains) int
	// Domain by index.
	Domain(int) Domain
	// Empty is true if all domains are empty.
	Empty() bool
	// Len returns the number of domains.
	Len() int
	// Remove values from a domain by index.
	Remove(int, ...int) Domains
	// Singleton is true if all domains are Singleton.
	Singleton() bool
	// Slices convert domains to a slice of int slices.
	Slices() [][]int
	// Values returns the values of a sequence of singleton domains/
	Values() ([]int, bool)

	/* Domain selectors */

	// First returns the first domain index with length above 1.
	First() (int, bool)
	// Largest returns the index of the largest domain with length above 1 by
	// number of elements.
	Largest() (int, bool)
	// Last returns the last domain index with length above 1.
	Last() (int, bool)
	// Maximum returns the index of the domain containing the maximum value with
	// length above 1.
	Maximum() (int, bool)
	// Minimum returns the index of the domain containing the minimum value with
	// length above 1.
	Minimum() (int, bool)
	// Smallest returns the index of the smallest domain with length above 1 by
	// number of elements.
	Smallest() (int, bool)
}

// A Range of integers.
type Range interface {
	Min() int
	Max() int
}

// NewDomain creates a domain of integers.
func NewDomain(ranges ...Range) Domain {
	connect()
	return newDomainFunc(ranges...)
}

// Singleton creates a domain containing one integer value.
func Singleton(value int) Domain {
	connect()
	return singletonFunc(value)
}

// Multiple creates a domain containing multiple integer values.
func Multiple(values ...int) Domain {
	connect()
	return multipleFunc(values...)
}

// NewDomains creates a sequence of domains.
func NewDomains(domains ...Domain) Domains {
	connect()
	return newDomainsFunc(domains...)
}

// Repeat a domain n times.
func Repeat(n int, d Domain) Domains {
	connect()
	return repeatFunc(n, d)
}

// NewRange create a new integer range.
func NewRange(min, max int) Range {
	connect()
	return newRangeFunc(min, max)
}

var (
	newDomainFunc  func(...Range) Domain
	singletonFunc  func(int) Domain
	multipleFunc   func(...int) Domain
	newDomainsFunc func(...Domain) Domains
	repeatFunc     func(int, Domain) Domains
	newRangeFunc   func(int, int) Range
)
