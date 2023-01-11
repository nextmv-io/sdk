package measure

import (
	"encoding/json"
	"sync/atomic"
)

// Override measure uses a default measure for all arcs that are not true for a
// condition. It uses an override measure for all arcs that are true for the
// condition.
func Override(
	defaultByIndex ByIndex,
	overrideByIndex ByIndex,
	condition func(from, to int) bool,
) ByIndex {
	return &override{
		defaultByIndex:  defaultByIndex,
		overrideByIndex: overrideByIndex,
		condition:       condition,
	}
}

// DebugOverride returns an Override that when marshalled will include debugging
// information describing the number of queries for default and override
// elements.
func DebugOverride(
	defaultByIndex ByIndex,
	overrideByIndex ByIndex,
	condition func(from, to int) bool,
) ByIndex {
	return &override{
		defaultByIndex:  defaultByIndex,
		overrideByIndex: overrideByIndex,
		condition:       condition,
		debugMode:       true,
	}
}

type override struct {
	defaultByIndex  ByIndex
	overrideByIndex ByIndex
	condition       func(from, to int) bool

	debugMode bool

	defaultCount  uint64
	overrideCount uint64
}

func (o *override) Cost(from, to int) float64 {
	if o.condition(from, to) {
		if o.debugMode {
			atomic.AddUint64(&o.overrideCount, 1)
		}
		return o.overrideByIndex.Cost(from, to)
	}

	if o.debugMode {
		atomic.AddUint64(&o.defaultCount, 1)
	}
	return o.defaultByIndex.Cost(from, to)
}

func (o *override) MarshalJSON() ([]byte, error) {
	output := map[string]any{
		"default":  o.defaultByIndex,
		"override": o.overrideByIndex,
		"type":     "override",
	}

	if o.debugMode {
		output["counts"] = map[string]uint64{
			"default":  o.defaultCount,
			"override": o.overrideCount,
		}
	}

	return json.Marshal(output)
}
