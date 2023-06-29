package routingkit_test

import (
	"encoding/json"
	"testing"
	"unicode"

	"github.com/nextmv-io/sdk/route"
	"github.com/nextmv-io/sdk/route/routingkit"
)

type byPointConstantMeasure float64

func (m byPointConstantMeasure) Cost(_ route.Point, _ route.Point) float64 {
	return float64(m)
}

func TestFallback(t *testing.T) {
	sources := []route.Point{
		{7.336650, 52.145020},
	}
	dests := []route.Point{
		{1.32486, 52.14280},
		{7.31893, 52.15924},
	}
	expected := [][]float64{
		{666, 2346},
	}

	m, err := routingkit.Matrix(
		"testdata/rk_test.osm.pbf",
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
		{7.336650, 52.145020},
		{7.33293, 52.13893},
		{7.33745, 52.14758},
		{7.34979, 52.15149},
	}
	dests := []route.Point{
		{7.31893, 52.15924},
		{7.35630, 52.14031},
	}
	expected := [][]float64{
		{2346, 2465},
		{2855, 2951},
		{2186, 1961},
		{2834, 2423},
	}

	m, err := routingkit.Matrix(
		"testdata/rk_test.osm.pbf",
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
		"testdata/rk_test.osm.pbf",
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
		`"osm":"testdata/rk_test.osm.pbf","profile":{"name":"car"},"radius":1000,` +
		`"sources":[[1,2]],"type":"routingkitMatrix"}`
	if v := string(b); v != w {
		t.Errorf("got %q; want %q", v, w)
	}
}

func TestByPoint(t *testing.T) {
	p1 := route.Point{7.33665, 52.14502}
	p2 := route.Point{7.33021, 52.14789}

	m, err := routingkit.ByPoint(
		"testdata/rk_test.osm.pbf",
		1000,
		1<<30,
		routingkit.Car(),
		nil,
	)
	if err != nil {
		t.Fatalf("constructing measure: %v", err)
	}

	if v := int(m.Cost(p1, p2)); v != 722 {
		t.Errorf("got %v; want 722", v)
	}
	if v := int(m.Cost(p2, p1)); v != 722 {
		t.Errorf("got %v; want 722", v)
	}

	// get the same values from the cache
	if v := int(m.Cost(p1, p2)); v != 722 {
		t.Errorf("got %v; want 722", v)
	}
	if v := int(m.Cost(p2, p1)); v != 722 {
		t.Errorf("got %v; want 722", v)
	}
}

func TestByPointMarshal(t *testing.T) {
	m, err := routingkit.ByPoint(
		"testdata/rk_test.osm.pbf",
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
	w := `{"cache_size":1073741824,"osm":"testdata/rk_test.osm.pbf",` +
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
			input: `{"cache_size":1073741824,"osm":"testdata/rk_test.osm.pbf",` +
				`"profile":{"name":"car"},"radius":1000,"type":"routingkit"}`,
			expectedErr: false,
			from:        route.Point{7.33665, 52.14502},
			to:          route.Point{7.33021, 52.14789},
			expected:    722,
		},
		{
			input: `{"cache_size":1073741824,"osm":"testdata/rk_test.osm.pbf",` +
				`"profile":{"name":"pedestrian"},"radius":1000,"type":"routingkit"}`,
			expectedErr: false,
			from:        route.Point{7.33665, 52.14502},
			to:          route.Point{7.33021, 52.14789},
			expected:    690,
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
			input: `{"destinations":[[7.31893,52.15924],[7.3563,` +
				`52.14031]],"osm":"testdata/rk_test.osm.pbf","profile":` + `
				{"name":"car"},"radius":1000,"sources":[[7.33665,52.14502],` +
				`[7.33293,52.13893],[7.33745,52.14758],[7.34979,52.15149]],` +
				`"type":"routingkitMatrix"}`,
			expectedErr: false,
			from:        2,
			to:          0,
			expected:    2186,
		},
		// with fallback measure
		{
			input: `{"destinations":[[7.32486,52.1428],[7.31893,` +
				`52.15924]],"measure":{"type":"haversine"},"osm":` +
				`"testdata/rk_test.osm.pbf","profile":{"name":"car"},` +
				`"radius":1000,"sources":[[7.33665,52.14502]],` +
				`"type":"routingkitMatrix"}`,
			expectedErr: false,
			from:        0,
			to:          0,
			expected:    1200,
		},
		// routingkitDurationMatrix
		{
			input: `{"destinations":[[7.31893,52.15924],[7.3563,` +
				`52.14031]],"osm":"testdata/rk_test.osm.pbf","profile":` +
				`{"name":"car"},"radius":1000,"sources":[[7.33665,52.14502],` +
				`[7.33293,52.13893],[7.33745,52.14758],[7.34979,52.15149]],` +
				`"type":"routingkitDurationMatrix"}`,
			expectedErr: false,
			from:        2,
			to:          0,
			expected:    221,
		},
		// routingkitDurationMatrix with fallback measure
		{
			input: `{"destinations":[[7.32486,52.1428],[7.31893,` +
				`52.15924]],"measure":{"type":"haversine"},"osm":` +
				`"testdata/rk_test.osm.pbf","profile":{"name":"car"},` +
				`"radius":1000,"sources":[[7.33665,52.14502]],"type":` +
				`"routingkitDurationMatrix"}`,
			expectedErr: false,
			from:        0,
			to:          0,
			expected:    246,
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
	p1 := route.Point{7.33665, 52.14502}
	p2 := route.Point{7.33021, 52.14789}

	m, err := routingkit.DurationByPoint(
		"testdata/rk_test.osm.pbf", 1000, 1<<30, routingkit.Car(), nil)
	if err != nil {
		t.Fatalf("constructing measure: %v", err)
	}

	if v := int(m.Cost(p1, p2)); v != 84 {
		t.Errorf("got %v; want 84", v)
	}
	if v := int(m.Cost(p2, p1)); v != 84 {
		t.Errorf("got %v; want 84", v)
	}

	// get the same values from the cache
	if v := int(m.Cost(p1, p2)); v != 84 {
		t.Errorf("got %v; want 84", v)
	}
	if v := int(m.Cost(p2, p1)); v != 84 {
		t.Errorf("got %v; want 84", v)
	}
}

func TestDurationByPointMarshal(t *testing.T) {
	m, err := routingkit.DurationByPoint(
		"testdata/rk_test.osm.pbf", 1000, 1<<30, routingkit.Car(), nil)
	if err != nil {
		t.Fatalf("constructing measure: %v", err)
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Errorf("got %+v; want nil", err)
	}
	w := `{"cache_size":1073741824,"osm":"testdata/rk_test.osm.pbf",` +
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
