package store

import (
	"encoding/json"
	"math"
)

// MarshalJSON Limits.
func (l Limits) MarshalJSON() ([]byte, error) {
	m := map[string]any{}
	m["duration"] = l.Duration.String()
	if l.Nodes != math.MaxInt {
		m["nodes"] = l.Nodes
	}
	if l.Solutions != math.MaxInt {
		m["solutions"] = l.Solutions
	}

	return json.Marshal(m)
}

// MarshalJSON Options.
func (o Options) MarshalJSON() ([]byte, error) {
	search := map[string]any{}
	search["buffer"] = o.Search.Buffer
	m := map[string]any{
		"diagram": o.Diagram,
		"search":  search,
	}
	if o.Limits != (Limits{}) {
		m["limits"] = o.Limits
	}
	if o.Random.Seed != 0 {
		m["random"] = o.Random
	}
	if o.Sense.String() != "" {
		m["sense"] = o.Sense.String()
	}
	if len(o.Tags) > 0 {
		m["tags"] = o.Tags
	}
	if o.Pool.Size != 0 {
		m["pool"] = o.Pool
	}

	return json.Marshal(m)
}

// MarshalJSON Diagram.
func (d Diagram) MarshalJSON() ([]byte, error) {
	m := map[string]any{"width": d.Width}
	m["expansion"] = d.Expansion

	return json.Marshal(m)
}

// MarshalJSON Time.
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"start":           t.Start,
		"elapsed":         t.Elapsed.String(),
		"elapsed_seconds": t.Elapsed.Seconds(),
	})
}
