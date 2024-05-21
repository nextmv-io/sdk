package decode_test

import (
	"os"
	"testing"

	"github.com/nextmv-io/nextroute/schema"
	"github.com/nextmv-io/sdk/nextroute/decode"
)

func TestDecodeTspLibParser(t *testing.T) {
	reader, err := os.Open("instances/120.2.vrp")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err = reader.Close()
		if err != nil {
			t.Error(err)
		}
	}()
	input := schema.Input{}
	decoder := decode.TSPLIBDecoder{}
	err = decoder.Decode(reader, &input)
	if err != nil {
		t.Error(err)
	}
	if len(input.Stops) != 119 {
		t.Errorf("expected %d, got %d", 119, len(input.Stops))
	}
}
