package alns

import (
	"github.com/nextmv-io/sdk/connect"
)

// NewSolveOperatorIndex returns the next unique solve operator index.
func NewSolveOperatorIndex() int {
	connect.Connect(con, &newSolveOperatorIndex)
	return newSolveOperatorIndex()
}
