package routingkit

import (
	"encoding/json"
	"fmt"

	"github.com/nextmv-io/go-routingkit/routingkit"
	"github.com/nextmv-io/sdk/route"
)

// ByPointLoader can be embedded in schema structs and unmarshals a ByPoint JSON
// object into the appropriate implementation, including a routingkit.ByPoint.
//
// Deprecated: This package is deprecated and will be removed in the next major release.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/routingkit].
type ByPointLoader struct {
	byPoint route.ByPoint
}

type byPointJSON struct {
	ByPoint       *ByPointLoader `json:"measure"` //nolint:tagliatelle
	Type          string         `json:"type"`
	OSMFile       string         `json:"osm"` //nolint:tagliatelle
	Radius        float64        `json:"radius"`
	CacheSize     int64          `json:"cache_size"`
	ProfileLoader *ProfileLoader `json:"profile"` //nolint:tagliatelle
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
	case "routingkit":
		byPoint, err := ByPoint(
			j.OSMFile,
			j.Radius,
			j.CacheSize,
			j.ProfileLoader.To(),
			j.ByPoint.To(),
		)
		if err != nil {
			return fmt.Errorf(`constructing measure: %v`, err)
		}
		l.byPoint = byPoint
	case "routingkitDuration":
		byPoint, err := DurationByPoint(
			j.OSMFile,
			j.Radius,
			j.CacheSize,
			j.ProfileLoader.To(),
			j.ByPoint.To(),
		)
		if err != nil {
			return fmt.Errorf(`constructing measure: %v`, err)
		}
		l.byPoint = byPoint
	default:
		var byPointLoader route.ByPointLoader
		if err := byPointLoader.UnmarshalJSON(b); err != nil {
			return err
		}
		l.byPoint = byPointLoader.To()
	}
	return nil
}

// To returns the underlying ByPoint.
func (l *ByPointLoader) To() route.ByPoint {
	if l == nil {
		return nil
	}
	return l.byPoint
}

// ByIndexLoader can be embedded in schema structs and unmarshals a ByIndex JSON
// object into the appropriate implementation, including a routingkit.ByIndex.
//
// Deprecated: This package is deprecated and will be removed in the next major release.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/routingkit].
type ByIndexLoader struct {
	byIndex route.ByIndex
}

type byIndexJSON struct {
	Measure       *ByPointLoader `json:"measure"`
	OSMFile       string         `json:"osm"` //nolint:tagliatelle
	Type          string         `json:"type"`
	Sources       []route.Point  `json:"sources"`
	Destinations  []route.Point  `json:"destinations"`
	Radius        float64        `json:"radius"`
	ProfileLoader *ProfileLoader `json:"profile,omitempty"` //nolint:tagliatelle
}

// MarshalJSON returns the JSON representation for the underlying ByIndex.
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

	var m route.ByPoint
	if j.Measure != nil {
		m = j.Measure.To()
	}

	switch j.Type {
	case "routingkitMatrix":
		byIndex, err := Matrix(
			j.OSMFile,
			j.Radius,
			j.Sources,
			j.Destinations,
			j.ProfileLoader.To(),
			m,
		)
		if err != nil {
			return fmt.Errorf(`constructing measure: %v`, err)
		}
		l.byIndex = byIndex
	case "routingkitDurationMatrix":
		byIndex, err := DurationMatrix(
			j.OSMFile,
			j.Radius,
			j.Sources,
			j.Destinations,
			j.ProfileLoader.To(),
			m,
		)
		if err != nil {
			return fmt.Errorf(`constructing measure: %v`, err)
		}
		l.byIndex = byIndex
	default:
		var byIndexLoader route.ByIndexLoader
		if err := byIndexLoader.UnmarshalJSON(b); err != nil {
			return err
		}
		l.byIndex = byIndexLoader.To()
	}
	return nil
}

// To returns the underlying ByIndex.
func (l *ByIndexLoader) To() route.ByIndex {
	return l.byIndex
}

// ProfileLoader can be embedded in schema structs and unmarshals a
// routingkit.Profile JSON object into the appropriate implementation.
//
// Deprecated: This package is deprecated and will be removed in the next major release.
// It is used with the router engine which was replaced by
// [github.com/nextmv-io/sdk/measure/routingkit].
type ProfileLoader struct {
	profile *routingkit.Profile
}

type profileJSON struct {
	Name string `json:"name"`
}

// MarshalJSON returns the JSON representation for the underlying Profile.
func (l ProfileLoader) MarshalJSON() ([]byte, error) {
	if l.profile == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(map[string]any{
		"name": l.profile.Name,
	})
}

// UnmarshalJSON converts the bytes into the appropriate implementation of
// Profile.
func (l *ProfileLoader) UnmarshalJSON(b []byte) error {
	var p profileJSON
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}
	var profile routingkit.Profile
	switch p.Name {
	case "car":
		profile = routingkit.Car()
	case "bike":
		profile = routingkit.Bike()
	case "pedestrian":
		profile = routingkit.Pedestrian()
	default:
		return fmt.Errorf(
			"%s is not an unmarshallable profile type: only car, bike, "+
				"and pedestrian can be unmarshalled",
			p.Name,
		)
	}
	l.profile = &profile

	return nil
}

// To returns the underlying Profile.
func (l *ProfileLoader) To() routingkit.Profile {
	if l == nil || l.profile == nil {
		return routingkit.Car()
	}
	return *l.profile
}
