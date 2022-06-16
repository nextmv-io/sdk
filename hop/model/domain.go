package model

import "github.com/nextmv-io/sdk/hop/model/types"

// Domain creates a domain of integers.
func Domain(ranges ...types.Range) types.Domain {
	connect()
	return domainFunc(ranges...)
}

// Singleton creates a domain containing one integer value.
func Singleton(value int) types.Domain {
	connect()
	return singletonFunc(value)
}

// Multiple creates a domain containing multiple integer values.
func Multiple(values ...int) types.Domain {
	connect()
	return multipleFunc(values...)
}

var (
	domainFunc    func(...types.Range) types.Domain
	singletonFunc func(int) types.Domain
	multipleFunc  func(...int) types.Domain
)
