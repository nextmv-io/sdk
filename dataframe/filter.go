package dataframe

import "fmt"

// Filter defines how to filter columns out of a DataFrame instance.
type Filter interface {
	fmt.Stringer

	// And creates and returns a conjunction filter of the invoking filter
	// and filter.
	And(filter Filter) Filter

	// Not creates and returns a negation filter of the invoking filter.
	Not() Filter

	// Or creates and returns a disjunction filter of the invoking filter
	// and filter.
	Or(filter Filter) Filter
}

// Filters is the slice of Filter instances.
type Filters []Filter
