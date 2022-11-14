package dataframe

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/run/decode"
)

// FromCSV returns a decoder to decode comma separated values (CSV) files and
// turns it into a DataFrame.
func FromCSV() decode.Decoder {
	connect.Connect(con, &fromCSV)
	return fromCSV()
}

// FromFeather returns a decoder to decode Apache Arrow Feather files in IPC
// format and turns it into a DataFrame. This decoder is not a streaming decoder
// and will load the entire file into memory. It does not support the Arrow
// Dictionary type.
func FromFeather() decode.Decoder {
	connect.Connect(con, &fromFeather)
	return fromFeather()
}

var (
	con         = connect.NewConnector("sdk", "DataFrame")
	fromCSV     func() decode.Decoder
	fromFeather func() decode.Decoder
)
