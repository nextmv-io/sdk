/*
Package types provides types.
*/
package types

import "time"

// Duration is a wrapper around time.Duration that implements json.Marshaler.
type Duration time.Duration

// MarshalJSON implements json.Marshaler. It calls time.Duration.String() to
// serialize the duration.
func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(time.Duration(d).String()), nil
}
