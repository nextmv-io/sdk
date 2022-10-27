package dataframe

import "fmt"

type Filter interface {
	fmt.Stringer

	And(Filter) Filter

	Not() Filter

	Or(Filter) Filter
}
