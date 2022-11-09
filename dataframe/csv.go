package dataframe

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/run/decode"
)

// FromCSV returns a decoder to decode comma separated values (CSV) files and
// turn it into a DataFrame.
func FromCSV() decode.Decoder {
	connect.Connect(con, &fromCSV)
	return fromCSV()
}

var (
	con     = connect.NewConnector("sdk", "DataFrame")
	fromCSV func() decode.Decoder
)
