package dataframe

import (
	"github.com/nextmv-io/sdk/connect"
	"github.com/nextmv-io/sdk/run/decode"
)

func FromCSV() decode.Decoder {
	connect.Connect(con, &fromCSV)
	return fromCSV()
}

var (
	con     = connect.NewConnector("sdk", "Dataframe")
	fromCSV func() decode.Decoder
)
