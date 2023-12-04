package decode_test

import (
	"os"
	"testing"

	"github.com/nextmv-io/sdk/nextroute/decode"
	"github.com/nextmv-io/sdk/nextroute/schema"
)

func TestDecodeTspLibParser(t *testing.T) {
	// read file in instances/120.2.vrp
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
	err = decode.Parse(reader, &input)
	if err != nil {
		t.Error(err)
	}
	if len(input.Stops) != 120 {
		t.Errorf("expected %d, got %d", 120, len(input.Stops))
	}
}
