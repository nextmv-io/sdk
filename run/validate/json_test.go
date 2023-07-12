package validate_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/nextmv-io/sdk/nextroute/schema"
	"github.com/nextmv-io/sdk/run/validate"
)

// TestNextrouteTemplateInput validates the input template against the schema.
func TestNextrouteTemplateInput(t *testing.T) {
	// read input file
	file, err := os.Open("../../templates/nextroute/input.json")
	if err != nil {
		t.Fatal(err)
	}
	err = validate.JSON[schema.Input](nil)(context.Background(), file)
	if err != nil {
		t.Fatal(err)
	}
}

// TestNextrouteAdditionalAttributeInput validates the input template against the schema.
// The input holds an additional attribute that is not defined in the schema and
// therefore schema validation should fail.
func TestNextrouteAdditionalAttributeInput(t *testing.T) {
	// read input file
	file, err := os.Open("testdata/nextroute_input_additional_attribute.json")
	if err != nil {
		t.Fatal(err)
	}
	err = validate.JSON[schema.Input](nil)(context.Background(), file)
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "Additional property some_attribute is not allowed") {
		t.Fatal("expected validation error to contain 'Additional property'")
	}
}

// TestNextrouteMissingFieldInput validates the input template against the
// schema.
// The input is missing a required field and therefore schema validation should
// fail.
func TestNextrouteMissingFieldInput(t *testing.T) {
	// read input file
	file, err := os.Open("testdata/nextroute_input_missing_field.json")
	if err != nil {
		t.Fatal(err)
	}
	err = validate.JSON[schema.Input](nil)(context.Background(), file)
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "defaults.vehicles.start_location: lon is required") {
		t.Fatal("expected validation error to contain 'defaults.vehicles.start_location: lon is required'")
	}
}
