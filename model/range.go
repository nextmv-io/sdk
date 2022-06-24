package model

// NewRange create a new integer range.
func NewRange(min, max int) Range {
	connect()
	return newRangeFunc(min, max)
}

var newRangeFunc func(int, int) Range
