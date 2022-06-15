package model

import "github.com/nextmv-io/sdk/hop/model/types"

// Range create a new integer range.
func Range(min, max int) types.Range {
	connect()
	return rangeFunc(min, max)
}

var rangeFunc func(int, int) types.Range
