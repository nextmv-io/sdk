package statistics

import (
	"encoding/json"
	"fmt"
	"math"
)

// Float64 is a float64 that can be marshaled to and from JSON.
// It supports the special values NaN, +Inf, and -Inf.
type Float64 float64

func (f Float64) String() string {
	v := float64(f)
	switch {
	case math.IsNaN(v):
		return "nan"
	case math.IsInf(v, 1):
		return "+inf"
	case math.IsInf(v, -1):
		return "-inf"
	default:
		return fmt.Sprintf("%f", v)
	}
}

// MarshalJSON marshals the Float64 to JSON.
func (f Float64) MarshalJSON() ([]byte, error) {
	v := float64(f)
	switch {
	case math.IsNaN(v):
		return []byte(`"nan"`), nil
	case math.IsInf(v, 1):
		return []byte(`"+inf"`), nil
	case math.IsInf(v, -1):
		return []byte(`"-inf"`), nil
	}
	return json.Marshal(v)
}

// UnmarshalJSON unmarshals the Float64 from JSON.
func (f *Float64) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		*f = Float64(x)
	case string:
		switch x {
		case "nan":
			*f = Float64(math.NaN())
		case "+inf", "inf":
			*f = Float64(math.Inf(1))
		case "-inf":
			*f = Float64(math.Inf(-1))
		default:
			return fmt.Errorf("invalid float64 string %q", x)
		}
	default:
		return fmt.Errorf("invalid float64 %v", v)
	}
	return nil
}
