// Package validate contains Validator implementations.
package validate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"

	humaSchema "github.com/danielgtaylor/huma/schema"
	"github.com/xeipuuv/gojsonschema"
)

// JSON creates a JSON validator. If nil is passed as schema, the validator will
// try to read schema.json in the current directory. If that file does not
// exist, no validation will be performed.
func JSON[Input any](schema []byte) func(_ context.Context, input any) error {
	return JSONValidator[Input]{
		schema: schema,
	}.Validate
}

// JSONValidator validates the input against a JSON schema.
type JSONValidator[Input any] struct {
	schema []byte
}

// Validate validates the input against a JSON schema.
func (j JSONValidator[Input]) Validate(_ context.Context, input any) (retErr error) {
	// no schema is given
	if len(j.schema) == 0 {
		// generate schema for input struct
		s, err := humaSchema.Generate(reflect.TypeOf(new(Input)))
		if err != nil {
			log.Fatal(err)
		}
		// serialize s to json
		schema, err := json.Marshal(s)
		if err != nil {
			log.Fatal(err)
		}
		j.schema = schema
	}

	schemaLoader := gojsonschema.NewBytesLoader(j.schema)
	// cast input to io.Reader
	reader, ok := input.(io.Reader)
	if !ok {
		return fmt.Errorf("input is not an io.Reader")
	}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return err
	}

	loader := gojsonschema.NewBytesLoader(buf.Bytes())

	result, err := gojsonschema.Validate(schemaLoader, loader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		sb := strings.Builder{}
		for _, desc := range result.Errors() {
			sb.WriteString(desc.String() + "\n")
		}
		return fmt.Errorf(sb.String())
	}
	return nil
}
