package dataframe

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/run/decode"
)

// FromCSV returns a decoder to decode comma seperated values (CSV) files and
// turn it into a DataFrame.
func FromCSV() decode.Decoder {
	connect.Connect(con, &fromCSV)
	return fromCSV()
}

var (
	con     = connect.NewConnector("sdk", "Dataframe")
	fromCSV func() decode.Decoder
)
