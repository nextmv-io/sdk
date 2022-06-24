package model

// NewDomain creates a domain of integers.
func NewDomain(ranges ...Range) Domain {
	connect()
	return newDomainFunc(ranges...)
}

// NewSingleton creates a domain containing one integer value.
func NewSingleton(value int) Domain {
	connect()
	return newSingletonFunc(value)
}

// NewMultiple creates a domain containing multiple integer values.
func NewMultiple(values ...int) Domain {
	connect()
	return newMultipleFunc(values...)
}

var (
	newDomainFunc    func(...Range) Domain
	newSingletonFunc func(int) Domain
	newMultipleFunc  func(...int) Domain
)
