package routingkit_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"unicode"

	"github.com/nextmv-io/sdk/measure"
	"github.com/nextmv-io/sdk/measure/routingkit"
	"github.com/twpayne/go-polyline"
)

const cacheSizeInput string = `{"cache_size":1073741824,"osm":"testdata/rk_test.osm.pbf",`

type byPointConstantMeasure float64

func (m byPointConstantMeasure) Cost(_ measure.Point, _ measure.Point) float64 {
	return float64(m)
}

func TestFallback(t *testing.T) {
	sources := []measure.Point{
		{7.336650, 52.145020},
	}
	dests := []measure.Point{
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
	sources := []measure.Point{
		{7.336650, 52.145020},
		{7.33293, 52.13893},
		{7.33745, 52.14758},
		{7.34979, 52.15149},
	}
	dests := []measure.Point{
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
		[]measure.Point{{1.0, 2.0}},
		[]measure.Point{{3.0, 4.0}},
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
	p1 := measure.Point{7.33665, 52.14502}
	p2 := measure.Point{7.33021, 52.14789}

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
	w := cacheSizeInput +
		`"profile":{"name":"pedestrian"},` +
		`"radius":1000,"type":"routingkit"}`
	if v := string(b); v != w {
		t.Errorf("got %q; want %q", v, w)
	}
}

func TestByPointLoader(t *testing.T) {
	tests := []struct {
		input       string
		from        measure.Point
		to          measure.Point
		expectedErr bool
		expected    int
	}{
		{
			input: cacheSizeInput +
				`"profile":{"name":"car"},"radius":1000,"type":"routingkit"}`,
			expectedErr: false,
			from:        measure.Point{7.33665, 52.14502},
			to:          measure.Point{7.33021, 52.14789},
			expected:    722,
		},
		{
			input: cacheSizeInput +
				`"profile":{"name":"pedestrian"},"radius":1000,"type":"routingkit"}`,
			expectedErr: false,
			from:        measure.Point{7.33665, 52.14502},
			to:          measure.Point{7.33021, 52.14789},
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
	p1 := measure.Point{7.33665, 52.14502}
	p2 := measure.Point{7.33021, 52.14789}

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
	w := cacheSizeInput +
		`"profile":{"name":"car"},` +
		`"radius":1000,"type":"routingkitDuration"}`
	if v := string(b); v != w {
		t.Errorf("got %q; want %q", v, w)
	}
}

func TestDistanceClient(t *testing.T) {
	testRoute := []measure.Point{
		{7.336189, 52.146548},
		{7.335031, 52.146057},
		{7.335073697312962, 52.145657185840214},
	}
	expectedCosts := [][]float64{
		{0, 165, 217},
		{165, 0, 52},
		{217, 52, 0},
	}
	expectedWholePoly := [][]float64{
		{7.3362, 52.14631},
		{7.3358799999999995, 52.14691},
		{7.3353399999999995, 52.146499999999996},
		{7.3350599999999995, 52.146209999999996},
		{7.3350599999999995, 52.146209999999996},
		{7.334689999999999, 52.14579},
	}
	expectedLegPolys := [][][]float64{
		{
			{7.3362, 52.14631},
			{7.3358799999999995, 52.14691},
			{7.3353399999999995, 52.146499999999996},
			{7.3350599999999995, 52.146209999999996},
		},
		{
			{7.33506, 52.14621},
			{7.33469, 52.145790000000005},
		},
	}

	c, err := routingkit.NewDistanceClient(
		"testdata/rk_test.osm.pbf", routingkit.Car())
	if err != nil {
		t.Fatalf("constructing measure: %v", err)
	}
	m, err := c.Measure(1000, 1<<30, nil)
	if err != nil {
		t.Fatalf("constructing measure: %v", err)
	}
	mat, err := c.Matrix(testRoute, testRoute)
	if err != nil {
		t.Fatalf("constructing matrix: %v", err)
	}

	checkClient(
		t,
		m,
		mat,
		func(points []measure.Point) (string, []string, error) {
			return c.Polyline(points)
		},
		testRoute,
		expectedCosts,
		expectedWholePoly,
		expectedLegPolys,
	)
}

func TestDurationClient(t *testing.T) {
	testRoute := []measure.Point{
		{7.336189, 52.146548},
		{7.335031, 52.146057},
		{7.335073697312962, 52.145657185840214},
	}
	expectedCosts := [][]float64{
		{0, 11.88, 15.624},
		{11.88, 0, 3.744},
		{15.624, 3.744, 0},
	}
	expectedWholePoly := [][]float64{
		{7.3362, 52.14631},
		{7.3358799999999995, 52.14691},
		{7.3353399999999995, 52.146499999999996},
		{7.3350599999999995, 52.146209999999996},
		{7.3350599999999995, 52.146209999999996},
		{7.334689999999999, 52.14579},
	}
	expectedLegPolys := [][][]float64{
		{
			{7.3362, 52.14631},
			{7.3358799999999995, 52.14691},
			{7.3353399999999995, 52.146499999999996},
			{7.3350599999999995, 52.146209999999996},
		},
		{
			{7.33506, 52.14621},
			{7.33469, 52.145790000000005},
		},
	}

	c, err := routingkit.NewDurationClient(
		"testdata/rk_test.osm.pbf", routingkit.Car())
	if err != nil {
		t.Fatalf("constructing measure: %v", err)
	}
	m, err := c.Measure(1000, 1<<30, nil)
	if err != nil {
		t.Fatalf("constructing measure: %v", err)
	}
	mat, err := c.Matrix(testRoute, testRoute)
	if err != nil {
		t.Fatalf("constructing matrix: %v", err)
	}

	checkClient(
		t,
		m,
		mat,
		func(points []measure.Point) (string, []string, error) {
			return c.Polyline(points)
		},
		testRoute,
		expectedCosts,
		expectedWholePoly,
		expectedLegPolys,
	)
}

// checkClient implements re-usable checks for testing polyline
// generation functionality. It takes a measure, a function that generates a
// polyline, a test route, and expected values for the cost matrix, the whole
// polyline (start of route to end), and the leg polylines (legs for each pair
// of points in the route).
func checkClient(
	t *testing.T,
	m measure.ByPoint,
	mat measure.ByIndex,
	polyliner func(points []measure.Point) (string, []string, error),
	testRoute []measure.Point,
	expectedCosts [][]float64,
	expectedWholePoly [][]float64,
	expectedLegPolys [][][]float64,
) {
	// Test by point measure
	for i, p := range testRoute {
		for j, q := range testRoute {
			if v := m.Cost(p, q); v != expectedCosts[i][j] {
				t.Errorf("got %v; want %v", v, expectedCosts[i][j])
			}
		}
	}

	// Test by index measure / matrix
	for i := range testRoute {
		for j := range testRoute {
			if v := mat.Cost(i, j); v != expectedCosts[i][j] {
				t.Errorf("got %v; want %v", v, expectedCosts[i][j])
			}
		}
	}

	// Test polyline generation
	poly, legs, err := polyliner(testRoute)
	if err != nil {
		t.Fatalf("error getting polyline: %v", err)
	}
	wholePoly, _, err := polyline.DecodeCoords([]byte(poly))
	legPolys := make([][][]float64, len(legs))
	for i, leg := range legs {
		legPolys[i], _, err = polyline.DecodeCoords([]byte(leg))
		if err != nil {
			t.Fatalf("error decoding leg polyline: %v", err)
		}
	}
	fmt.Println(wholePoly)
	fmt.Println(legPolys)
	if err != nil {
		t.Fatalf("error decoding polyline: %v", err)
	}
	if len(wholePoly) != len(expectedWholePoly) {
		t.Fatalf("got %d points; want %d", len(wholePoly), len(expectedWholePoly))
	}
	for i, p := range wholePoly {
		if !equalPoint(p, expectedWholePoly[i]) {
			t.Errorf("got %v polygon point at index %d; want %v", p, i, expectedWholePoly[i])
		}
	}
	if len(legPolys) != len(expectedLegPolys) {
		t.Fatalf("got %d legs; want %d", len(legPolys), len(expectedLegPolys))
	}
	for i, leg := range legPolys {
		if len(leg) != len(expectedLegPolys[i]) {
			t.Errorf("got %d points in leg %d; want %d", len(leg), i, len(expectedLegPolys[i]))
		}
		for j, p := range leg {
			if !equalPoint(p, expectedLegPolys[i][j]) {
				t.Errorf("got %v point at index %d in leg at index %d; want %v", p, j, i, expectedLegPolys[i][j])
			}
		}
	}
}

func equalPoint(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i, ai := range a {
		if ai != b[i] {
			return false
		}
	}
	return true
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
