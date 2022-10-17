package routingkit_test

import (
	"encoding/json"
	"testing"
	"unicode"

	"github.com/nextmv-io/sdk/route"
	"github.com/nextmv-io/sdk/route/routingkit"
)

type byPointConstantMeasure float64

func (m byPointConstantMeasure) Cost(a route.Point, b route.Point) float64 {
	return float64(m)
}

func TestFallback(t *testing.T) {
	sources := []route.Point{
		{-76.587490, 39.299710},
	}
	dests := []route.Point{
		{-76.60548, 39.30772},
		{-76.582855, 39.309095},
	}
	expected := [][]float64{
		{666, 1496},
	}

	m, err := routingkit.Matrix(
		"testdata/maryland.osm.pbf",
		1000,
		sources,
		dests,
		routingkit.Car(),
		byPointConstantMeasure(666),
	)
	if err != nil {
		t.Fatalf("constructing measure: %v", err)
	}
	for i := range sources {
		for j := range dests {
			v := m.Cost(i, j)
			if v != expected[i][j] {
				t.Errorf("[%d,%d] expected %f, got %f", i, j, expected[i][j], v)
			}
		}
	}
}

func TestMatrix(t *testing.T) {
	sources := []route.Point{
		{-76.587490, 39.299710},
		{-76.594045, 39.300524},
		{-76.586664, 39.290938},
		{-76.598423, 39.289484},
	}
	dests := []route.Point{
		{-76.582855, 39.309095},
		{-76.599388, 39.302014},
	}
	expected := [][]float64{
		{1496, 1259},
		{1831, 575},
		{2372, 2224},
		{3399, 1548},
	}

	m, err := routingkit.Matrix(
		"testdata/maryland.osm.pbf",
		1000,
		sources,
		dests,
		routingkit.Car(),
		nil,
	)
	if err != nil {
		t.Fatalf("constructing measure: %v", err)
	}
	for i := range expected {
		for j, expectedV := range expected[i] {
			v := m.Cost(i, j)
			if v != expectedV {
				t.Errorf("[%d,%d] expected %f, got %f", i, j, expectedV, v)
			}
		}
	}
}

func TestMatrixMarshal(t *testing.T) {
	m, err := routingkit.Matrix(
		"testdata/maryland.osm.pbf",
		1000,
		[]route.Point{{1.0, 2.0}},
		[]route.Point{{3.0, 4.0}},
		routingkit.Car(),
		nil,
	)
	if err != nil {
		t.Fatalf("constructing measure: %v", err)
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Errorf("got %+v; want nil", err)
	}
	w := `{"destinations":[[3,4]],` +
		`"osm":"testdata/maryland.osm.pbf","profile":{"name":"car"},"radius":1000,` +
		`"sources":[[1,2]],"type":"routingkitMatrix"}`
	if v := string(b); v != w {
		t.Errorf("got %q; want %q", v, w)
	}
}

func TestByPoint(t *testing.T) {
	p1 := route.Point{-76.58749, 39.29971}
	p2 := route.Point{-76.59735, 39.30587}

	m, err := routingkit.ByPoint(
		"testdata/maryland.osm.pbf",
		1000,
		1<<30,
		routingkit.Car(),
		nil,
	)
	if err != nil {
		t.Fatalf("constructing measure: %v", err)
	}

	if v := int(m.Cost(p1, p2)); v != 1567 {
		t.Errorf("got %v; want 1567", v)
	}
	if v := int(m.Cost(p2, p1)); v != 1706 {
		t.Errorf("got %v; want 1706", v)
	}

	// get the same values from the cache
	if v := int(m.Cost(p1, p2)); v != 1567 {
		t.Errorf("got %v; want 1567", v)
	}
	if v := int(m.Cost(p2, p1)); v != 1706 {
		t.Errorf("got %v; want 1706", v)
	}
}

func TestByPointMarshal(t *testing.T) {
	m, err := routingkit.ByPoint(
		"testdata/maryland.osm.pbf",
		1000,
		1<<30,
		routingkit.Pedestrian(),
		nil)
	if err != nil {
		t.Fatalf("constructing measure: %v", err)
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Errorf("got %+v; want nil", err)
	}
	w := `{"cache_size":1073741824,"osm":"testdata/maryland.osm.pbf",` +
		`"profile":{"name":"pedestrian"},` +
		`"radius":1000,"type":"routingkit"}`
	if v := string(b); v != w {
		t.Errorf("got %q; want %q", v, w)
	}
}

func TestByPointLoader(t *testing.T) {
	tests := []struct {
		input       string
		from        route.Point
		to          route.Point
		expectedErr bool
		expected    int
	}{
		{
			input: `{"cache_size":1073741824,"osm":"testdata/maryland.osm.pbf",` +
				`"profile":{"name":"car"},"radius":1000,"type":"routingkit"}`,
			expectedErr: false,
			from:        route.Point{-76.58749, 39.29971},
			to:          route.Point{-76.59735, 39.30587},
			expected:    1567,
		},
		{
			input: `{"cache_size":1073741824,"osm":"testdata/maryland.osm.pbf",` +
				`"profile":{"name":"pedestrian"},"radius":1000,"type":"routingkit"}`,
			expectedErr: false,
			from:        route.Point{-76.58749, 39.29971},
			to:          route.Point{-76.59735, 39.30587},
			expected:    1555,
		},
	}
	for i, test := range tests {
		var loader routingkit.ByPointLoader
		if err := json.Unmarshal([]byte(test.input), &loader); err != nil {
			if !test.expectedErr {
				t.Errorf("[%d] unexpected error: %v", i, err)
			}
			continue
		}
		if test.expectedErr {
			t.Errorf("[%d] expected error but got none", i)
			continue
		}
		res := loader.To().Cost(test.from, test.to)
		if int(res) != test.expected {
			t.Errorf("[%d] expected %d, got %d", i, test.expected, int(res))
		}
		marshalled, err := json.Marshal(loader)
		if err != nil {
			t.Errorf("error marshalling loader: %v", err)
		}
		got := string(marshalled)
		want := removeSpace(test.input)
		if got != want {
			t.Errorf("[%d] got %s, want %s", i, got, want)
		}
	}
}

func TestByIndexLoader(t *testing.T) {
	tests := []struct {
		input       string
		from        int
		to          int
		expectedErr bool
		expected    int
	}{
		{
			input: `{"destinations":[[-76.582855,39.309095],[-76.599388,39.302014]],` +
				`"osm":"testdata/maryland.osm.pbf","profile":{"name":"car"},` +
				`"radius":1000,` +
				`"sources":[[-76.58749,39.29971],[-76.594045,39.300524],` +
				`[-76.586664,39.290938],[-76.598423,39.289484]],` +
				`"type":"routingkitMatrix"}`,
			expectedErr: false,
			from:        2,
			to:          0,
			expected:    2372,
		},
		// with fallback measure
		{
			input: `{"destinations":[[-76.60548,39.30772],[-76.582855,39.309095]],` +
				`"measure":{"type":"haversine"},` +
				`"osm":"testdata/maryland.osm.pbf","profile":{"name":"car"},` +
				`"radius":1000,"sources":[[-76.58749,39.29971]],` +
				`"type":"routingkitMatrix"}`,
			expectedErr: false,
			from:        0,
			to:          0,
			expected:    1785,
		},
		// routingkitDurationMatrix
		{
			input: `{"destinations":[[-76.582855,39.309095],[-76.599388,39.302014]],` +
				`"osm":"testdata/maryland.osm.pbf",` +
				`"profile":{"name":"car"},"radius":1000,` +
				`"sources":[[-76.58749,39.29971],[-76.594045,39.300524],` +
				`[-76.586664,39.290938],[-76.598423,39.289484]],` +
				`"type":"routingkitDurationMatrix"}`,
			expectedErr: false,
			from:        2,
			to:          0,
			expected:    205,
		},
		// routingkitDurationMatrix with fallback measure
		{
			input: `{"destinations":[[-76.60548,39.30772],[-76.582855,39.309095]],` +
				`"measure":{"type":"haversine"},` +
				`"osm":"testdata/maryland.osm.pbf","profile":{"name":"car"},` +
				`"radius":1000,"sources":[[-76.58749,39.29971]],` +
				`"type":"routingkitDurationMatrix"}`,
			expectedErr: false,
			from:        0,
			to:          0,
			expected:    1785,
		},
	}
	for i, test := range tests {
		var loader routingkit.ByIndexLoader
		if err := json.Unmarshal([]byte(test.input), &loader); err != nil {
			if !test.expectedErr {
				t.Errorf("[%d] unexpected error: %v", i, err)
			}
			continue
		}
		if test.expectedErr {
			t.Errorf("[%d] expected error but got none", i)
			continue
		}
		res := loader.To().Cost(test.from, test.to)
		if int(res) != test.expected {
			t.Errorf("[%d] expected %d, got %d", i, test.expected, int(res))
		}
		marshalled, err := json.Marshal(loader)
		if err != nil {
			t.Errorf("error marshalling loader: %v", err)
		}
		got := string(marshalled)
		want := removeSpace(test.input)
		if got != want {
			t.Errorf("[%d] got %s, want %s", i, got, want)
		}
	}
}

func TestDurationByPoint(t *testing.T) {
	p1 := route.Point{-76.58749, 39.29971}
	p2 := route.Point{-76.59735, 39.30587}

	m, err := routingkit.DurationByPoint(
		"testdata/maryland.osm.pbf", 1000, 1<<30, routingkit.Car(), nil)
	if err != nil {
		t.Fatalf("constructing measure: %v", err)
	}

	if v := int(m.Cost(p1, p2)); v != 140 {
		t.Errorf("got %v; want 140", v)
	}
	if v := int(m.Cost(p2, p1)); v != 156 {
		t.Errorf("got %v; want 156", v)
	}

	// get the same values from the cache
	if v := int(m.Cost(p1, p2)); v != 140 {
		t.Errorf("got %v; want 140", v)
	}
	if v := int(m.Cost(p2, p1)); v != 156 {
		t.Errorf("got %v; want 156", v)
	}
}

func TestDurationByPointMarshal(t *testing.T) {
	m, err := routingkit.DurationByPoint(
		"testdata/maryland.osm.pbf", 1000, 1<<30, routingkit.Car(), nil)
	if err != nil {
		t.Fatalf("constructing measure: %v", err)
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Errorf("got %+v; want nil", err)
	}
	w := `{"cache_size":1073741824,"osm":"testdata/maryland.osm.pbf",` +
		`"profile":{"name":"car"},` +
		`"radius":1000,"type":"routingkitDuration"}`
	if v := string(b); v != w {
		t.Errorf("got %q; want %q", v, w)
	}
}

func removeSpace(s string) string {
	rr := make([]rune, 0, len(s))
	for _, r := range s {
		if !unicode.IsSpace(r) {
			rr = append(rr, r)
		}
	}
	return string(rr)
}
