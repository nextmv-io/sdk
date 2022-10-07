// © 2019-2022 nextmv.io inc. All rights reserved.
// nextmv.io, inc. CONFIDENTIAL
//
// This file includes unpublished proprietary source code of nextmv.io, inc.
// The copyright notice above does not evidence any actual or intended
// publication of such source code. Disclosure of this source code or any
// related proprietary information is strictly prohibited without the express
// written permission of nextmv.io, inc.

package route

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ByPointLoader can be embedded in schema structs and unmarshals a ByPoint JSON
// object into the appropriate implementation.
type ByPointLoader struct {
	byPoint ByPoint
}

type pointType string

const (
	typeScale     pointType = "scale"
	typeEuclidean pointType = "euclidean"
	typeHaversine pointType = "haversine"
	typeTaxicab   pointType = "taxicab"
	typeConstant  pointType = "constant"
)

type byPointJSON struct {
	ByPoint  *ByPointLoader `json:"measure"`
	Type     pointType      `json:"type"`
	Scale    float64        `json:"scale"`
	Constant float64        `json:"constant"`
}

// MarshalJSON returns the JSON representation for the underlying ByPoint.
func (l ByPointLoader) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.byPoint)
}

// UnmarshalJSON converts the bytes into the appropriate implementation of
// ByPoint.
func (l *ByPointLoader) UnmarshalJSON(b []byte) error {
	var j byPointJSON
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	switch j.Type {
	case "":
		return errors.New(`no "type" field in json input`)
	case typeEuclidean:
		l.byPoint = EuclideanByPoint()
	case typeHaversine:
		l.byPoint = HaversineByPoint()
	case typeTaxicab:
		l.byPoint = TaxicabByPoint()
	case typeScale:
		l.byPoint = ScaleByPoint(j.ByPoint.To(), j.Scale)
	case typeConstant:
		l.byPoint = ConstantByPoint(j.Constant)
	default:
		return fmt.Errorf(`invalid type "%s"`, j.Type)
	}
	return nil
}

// To returns the underlying ByPoint.
func (l *ByPointLoader) To() ByPoint {
	return l.byPoint
}

// ByIndexLoader can be embedded in schema structs and unmarshals a ByIndex JSON
// object into the appropriate implementation.
type ByIndexLoader struct {
	byIndex ByIndex
}

// byIndexJSON includes the union of all fields that may appear on a ByIndex
// JSON object (like a C oneof). We unmarshal onto this data structure instead
// of onto a map[string]any for type safety and because this will allow
// recursive measures to be automatically unmarshalled.
type byIndexJSON struct {
	ByIndex  *ByIndexLoader          `json:"measure"`
	Arcs     map[int]map[int]float64 `json:"arcs"`
	Type     string                  `json:"type"`
	Measures []ByIndexLoader         `json:"measures"`
	Costs    []float64               `json:"costs"`
	Matrix   [][]float64             `json:"matrix"`
	Constant float64                 `json:"constant"`
	Scale    float64                 `json:"scale"`
	Exponent float64                 `json:"exponent"`
	Lower    float64                 `json:"lower"`
	Upper    float64                 `json:"upper"`
}

// MarshalJSON returns the JSON representation for the underlying Byindex.
func (l ByIndexLoader) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.byIndex)
}

// UnmarshalJSON converts the bytes into the appropriate implementation of
// ByIndex.
func (l *ByIndexLoader) UnmarshalJSON(b []byte) error {
	var j byIndexJSON
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	requiresByIndex := j.Type == "location" ||
		j.Type == "power" ||
		j.Type == "scale" ||
		j.Type == "sparse" ||
		j.Type == "truncate"
	if requiresByIndex && j.ByIndex == nil {
		return errors.New(`location measure must include a "by_index" field`)
	}

	switch j.Type {
	case "":
		return errors.New(`no "type" field in json input`)
	case "constant":
		l.byIndex = Constant(j.Constant)
	case "sum":
		measures := make([]ByIndex, len(j.Measures))
		for i, l := range j.Measures {
			measures[i] = l.To()
		}
		l.byIndex = Sum(measures...)
	case "location":
		l.byIndex, _ = Location(j.ByIndex.To(), j.Costs, nil)
	case "matrix":
		l.byIndex = Matrix(j.Matrix)
	case "power":
		l.byIndex = Power(j.ByIndex.To(), j.Exponent)
	case "scale":
		l.byIndex = Scale(j.ByIndex.To(), j.Scale)
	case "sparse":
		l.byIndex = Sparse(j.ByIndex.To(), j.Arcs)
	case "truncate":
		l.byIndex = Truncate(j.ByIndex.To(), j.Lower, j.Upper)
	default:
		return fmt.Errorf(`invalid type "%s"`, j.Type)
	}

	return nil
}

// To returns the underlying ByIndex.
func (l *ByIndexLoader) To() ByIndex {
	return l.byIndex
}
