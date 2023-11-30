package dataframe

import "fmt"

// Filter defines how to filter columns out of a DataFrame instance.
//
// Deprecated: This package is deprecated and will be removed in a future.
type Filter interface {
	fmt.Stringer

	// And creates and returns a conjunction filter of the invoking filter
	// and filter.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	And(filter Filter) Filter

	// Not creates and returns a negation filter of the invoking filter.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	Not() Filter

	// Or creates and returns a disjunction filter of the invoking filter
	// and filter.
	//
	// Deprecated: This package is deprecated and will be removed in a future.
	Or(filter Filter) Filter
}

// Filters is the slice of Filter instances.
//
// Deprecated: This package is deprecated and will be removed in a future.
type Filters []Filter
