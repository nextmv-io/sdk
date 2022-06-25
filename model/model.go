package model

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
