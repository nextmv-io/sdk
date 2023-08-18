package osrm_test

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/nextmv-io/sdk/measure"
	"github.com/nextmv-io/sdk/route"
	"github.com/nextmv-io/sdk/route/osrm"
)

var expectedDistances = [3][3]float64{
	{0, 10283.9, 8160.5},
	{10323.8, 0, 1931.5},
	{8441.5, 1378, 0},
}

var expectedDurations = [3][3]float64{
	{0, 781.1, 767.1},
	{804.9, 0, 187.2},
	{793.5, 119.5, 0},
}

var expectedDurationsUae = [3][3]float64{
	{0, 1111.7, 1154.6},
	{1295.4, 0, 761.1},
	{1167.7, 505.1, 0},
}

var expectedDurationsScaledBy2 = [3][3]float64{
	{0, 1562.2, 1534.2},
	{1609.8, 0, 374.4},
	{1587.0, 239, 0},
}

var p = []route.Point{
	{-105.050583, 39.762631},
	{-104.983978, 39.711413},
	{-104.983978, 39.721413},
}

var oceanPoints = []route.Point{
	{-42.79808862899699, 28.670649170472345},
	{-41.847832598419565, 17.13854940278748},
}

const (
	hostEnv          = "ENGINES_TEST_OSRM_HOST"
	hostEnvNotSetMsg = "skipping because " + hostEnv + " is not set"
)

func newMockOSRM(
	distances [][]float64,
	durations [][]float64,
) *httptest.Server {
	return httptest.NewServer(
		http.AllowQuerySemicolons(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_ = json.NewEncoder(w).Encode(map[string]any{
					"code":      "Ok",
					"distances": distances,
					"durations": durations,
					"message":   "Everything worked",
				})
			}),
		),
	)
}

func TestDistanceMatrixWithZeroPoints(t *testing.T) {
	type matrixFunc func(c osrm.Client,
		points []route.Point,
		parallelQueries int,
	) (route.ByIndex, error)

	selectDistance := func(c osrm.Client,
		points []route.Point,
		parallelQueries int,
	) (route.ByIndex, error) {
		distance, _, err := osrm.DistanceDurationMatrices(c, points, parallelQueries)
		return distance, err
	}
	selectDuration := func(c osrm.Client,
		points []route.Point,
		parallelQueries int,
	) (route.ByIndex, error) {
		_, duration, err := osrm.DistanceDurationMatrices(c, points, parallelQueries)
		return duration, err
	}
	type testCase struct {
		matrixFunc        matrixFunc
		distancesResponse [][]float64
		durationsResponse [][]float64
		points            []route.Point
		expected          [][]float64
	}
	matrixWithInfs := [][]float64{
		{1.0, 2.0, 999.0, 3.0, 999.0},
		{4.0, 5.0, 999.0, 6.0, 999.0},
		{999.0, 999.0, 999.0, 999.0, 999.0},
		{7.0, 8.0, 999.0, 9.0, 999.0},
	}
	matrixWithZeroes := [][]float64{
		{1.0, 2.0, 0.0, 3.0, 0.0},
		{4.0, 5.0, 0.0, 6.0, 0.0},
		{0.0, 0.0, 0.0, 0.0, 0.0},
		{7.0, 8.0, 0.0, 9.0, 0.0},
		{0.0, 0.0, 0.0, 0.0, 0.0},
	}
	matrixWithNoInfs := [][]float64{
		{1.0, 2.0, 3.0},
		{4.0, 5.0, 6.0},
		{7.0, 8.0, 9.0},
	}
	cases := []testCase{
		{
			points: []route.Point{
				{1.0, 1.0},
				{2.0, 2.0},
				{},
				{3.0, 3.0},
				{},
			},
			distancesResponse: matrixWithInfs,
			expected:          matrixWithZeroes,
			matrixFunc:        osrm.DistanceMatrix,
		},
		{
			points: []route.Point{
				{1.0, 1.0},
				{2.0, 2.0},
				{},
				{3.0, 3.0},
				{},
			},
			distancesResponse: matrixWithInfs,
			expected:          matrixWithZeroes,
			matrixFunc:        osrm.DistanceMatrix,
		},
		{
			points: []route.Point{
				{1.0, 1.0},
				{2.0, 2.0},
				{},
				{3.0, 3.0},
				{},
			},
			durationsResponse: matrixWithInfs,
			expected:          matrixWithZeroes,
			matrixFunc:        osrm.DurationMatrix,
		},
		{
			points: []route.Point{
				{1.0, 1.0},
				{2.0, 2.0},
				{},
				{3.0, 3.0},
				{},
			},
			distancesResponse: matrixWithInfs,
			expected:          matrixWithZeroes,
			matrixFunc:        selectDistance,
		},
		{
			points: []route.Point{
				{1.0, 1.0},
				{2.0, 2.0},
				{3.0, 3.0},
			},
			distancesResponse: matrixWithNoInfs,
			expected:          matrixWithNoInfs,
			matrixFunc:        osrm.DistanceMatrix,
		},
		{
			points: []route.Point{
				{1.0, 1.0},
				{2.0, 2.0},
				{3.0, 3.0},
			},
			durationsResponse: matrixWithNoInfs,
			expected:          matrixWithNoInfs,
			matrixFunc:        osrm.DurationMatrix,
		},
		{
			points: []route.Point{
				{1.0, 1.0},
				{2.0, 2.0},
				{3.0, 3.0},
			},
			durationsResponse: matrixWithNoInfs,
			expected:          matrixWithNoInfs,
			matrixFunc:        selectDuration,
		},
	}
	for i, test := range cases {
		s := newMockOSRM(test.distancesResponse, test.durationsResponse)
		c := osrm.DefaultClient(s.URL, false)
		m, err := test.matrixFunc(c, test.points, 0)
		if err != nil {
			t.Errorf("[%d] unexpected error: %v", i, err)
		}
		got := unpackMeasure(m, test.points)
		if !reflect.DeepEqual(test.expected, got) {
			t.Errorf("[%d] expected %v, got %v", i, test.expected, got)
		}
		s.Close()
	}
}

func unpackMeasure(m route.ByIndex, points []route.Point) [][]float64 {
	matrix := make([][]float64, len(points))
	for i := range matrix {
		matrix[i] = make([]float64, len(points))
		for j := range matrix[i] {
			matrix[i][j] = m.Cost(i, j)
		}
	}
	return matrix
}

func TestDistanceMatrix(t *testing.T) {
	osrmHost := os.Getenv(hostEnv)
	if osrmHost == "" {
		t.Skip(hostEnvNotSetMsg)
	}

	c := osrm.DefaultClient(osrmHost, true)
	disMeasurer, err := osrm.DistanceMatrix(c, p, 0)
	if err != nil {
		t.Fatalf("error requesting matrix: %v", err)
	}

	for i, row := range expectedDistances {
		for j, col := range row {
			if c := disMeasurer.Cost(i, j); c != col {
				t.Errorf("want: %f; got: %f\n", col, c)
			}
		}
	}
}

func TestSnappingFailed(t *testing.T) {
	osrmHost := os.Getenv(hostEnv)
	if osrmHost == "" {
		t.Skip(hostEnvNotSetMsg)
	}

	c := osrm.DefaultClient(osrmHost, true)
	_, err := osrm.DistanceMatrix(c, oceanPoints, 0)
	if err == nil {
		t.Fatalf("snapping should have failed")
	}
	expectedErrorMsg := "fetching matrix: expected \"Ok\" response code; got" +
		" \"NoSegment\" (\"Could not find a matching segment for coordinate 0\")"
	actualErrorMsg := err.Error()
	if actualErrorMsg != expectedErrorMsg {
		t.Errorf("want: %v; got: %v\n", expectedErrorMsg, actualErrorMsg)
	}
}

func TestDeflate(t *testing.T) {
	osrmHost := os.Getenv(hostEnv)
	if osrmHost == "" {
		t.Skip(hostEnvNotSetMsg)
	}

	// Extend base matrix to include a 0,0 point
	points := make([]measure.Point, len(p)+1)
	copy(points, p)
	points[len(p)] = measure.Point{0, 0}
	expDistances := make([][]float64, len(points))
	for i := range expectedDistances {
		expDistances[i] = make([]float64, len(points))
		copy(expDistances[i], expectedDistances[i][:])
	}
	expDistances[len(points)-1] = make([]float64, len(points))
	expDurations := make([][]float64, len(points))
	for i := range expectedDurations {
		expDurations[i] = make([]float64, len(points))
		copy(expDurations[i], expectedDurations[i][:])
	}
	expDurations[len(points)-1] = make([]float64, len(points))

	c := osrm.DefaultClient(osrmHost, true)
	c.IgnoreEmpty(true)
	distances, durations, err := c.Table(points, osrm.WithDistance(), osrm.WithDuration())
	if err != nil {
		t.Fatalf("error requesting matrices: %v", err)
	}

	if !reflect.DeepEqual(expDistances, distances) {
		t.Errorf("want: %v; got: %v\n", expDistances, distances)
	}
	if !reflect.DeepEqual(expDurations, durations) {
		t.Errorf("want: %v; got: %v\n", expDurations, durations)
	}
}

func TestDistanceMatrixWithParallelQueries(t *testing.T) {
	osrmHost := os.Getenv(hostEnv)
	if osrmHost == "" {
		t.Skip(hostEnvNotSetMsg)
	}

	c := osrm.DefaultClient(osrmHost, true)
	disMeasurer, err := osrm.DistanceMatrix(c, p, 2)
	if err != nil {
		t.Fatalf("error requesting matrix: %v", err)
	}

	for i, row := range expectedDistances {
		for j, col := range row {
			if c := disMeasurer.Cost(i, j); c != col {
				t.Errorf("want: %f; got: %f\n", col, c)
			}
		}
	}
}

func TestDurationMatrix(t *testing.T) {
	osrmHost := os.Getenv(hostEnv)
	if osrmHost == "" {
		t.Skip(hostEnvNotSetMsg)
	}

	c := osrm.DefaultClient(osrmHost, true)
	durMeasurer, err := osrm.DurationMatrix(c, p, 0)
	if err != nil {
		t.Fatalf("error requesting matrix: %v", err)
	}

	for i, row := range expectedDurations {
		for j, col := range row {
			if c := durMeasurer.Cost(i, j); c != col {
				t.Errorf("want: %f; got: %f\n", col, c)
			}
		}
	}
}

// You will need http://download.geofabrik.de/asia/gcc-states.html for this
// test.
func TestDurationMatrixLarge(t *testing.T) {
	osrmHost := os.Getenv(hostEnv)
	if osrmHost == "" {
		t.Skip(hostEnvNotSetMsg)
	}

	testpoints := make([]route.Point, 1000)
	for i := 0; i < len(testpoints); i++ {
		latMin := 24.17041401832874
		latMax := 24.253760
		lat := latMin + rand.Float64()*(latMax-latMin)

		lonMin := 55.656787
		lonMax := 55.81164689921658
		lon := lonMin + rand.Float64()*(lonMax-lonMin)
		testpoints[i] = route.Point{lon, lat}
	}

	c := osrm.DefaultClient(osrmHost, true)

	durMeasurer, err := osrm.DurationMatrix(c, testpoints, 0)
	if err != nil {
		t.Fatalf("error requesting matrix: %v", err)
	}

	// Check some expected durations
	for i, row := range expectedDurationsUae {
		for j, col := range row {
			if c := durMeasurer.Cost(i, j); c != col {
				t.Errorf("want: %f; got: %f\n", col, c)
			}
		}
	}

	// Check that all durations are present and non-negative
	for i := range testpoints {
		for j := range testpoints {
			if c := durMeasurer.Cost(i, j); c < 0 {
				t.Errorf("received negative duration: %f\n", c)
			}
		}
	}
}

func TestDurationMatrixWithScaleFactor(t *testing.T) {
	osrmHost := os.Getenv(hostEnv)
	if osrmHost == "" {
		t.Skip(hostEnvNotSetMsg)
	}

	c := osrm.DefaultClient(osrmHost, true)
	err := c.ScaleFactor(2.0)
	if err != nil {
		t.Fatalf("error requesting matrix with scale factor: %v", err)
	}
	durMeasurer, err := osrm.DurationMatrix(c, p, 0)
	if err != nil {
		t.Fatalf("error requesting matrix: %v", err)
	}

	for i, row := range expectedDurationsScaledBy2 {
		for j, col := range row {
			if c := durMeasurer.Cost(i, j); c != col {
				t.Errorf("want: %f; got: %f\n", col, c)
			}
		}
	}
}

func TestDurationMatrixWithParallelQueries(t *testing.T) {
	osrmHost := os.Getenv(hostEnv)
	if osrmHost == "" {
		t.Skip(hostEnvNotSetMsg)
	}

	c := osrm.DefaultClient(osrmHost, true)
	durMeasurer, err := osrm.DurationMatrix(c, p, 2)
	if err != nil {
		t.Fatalf("error requesting matrix: %v", err)
	}

	for i, row := range expectedDurations {
		for j, col := range row {
			if c := durMeasurer.Cost(i, j); c != col {
				t.Errorf("want: %f; got: %f\n", col, c)
			}
		}
	}
}

func TestDistanceDurationMatrices(t *testing.T) {
	osrmHost := os.Getenv(hostEnv)
	if osrmHost == "" {
		t.Skip(hostEnvNotSetMsg)
	}

	c := osrm.DefaultClient(osrmHost, true)
	disMeasurer, durMeasurer, err := osrm.DistanceDurationMatrices(c, p, 0)
	if err != nil {
		t.Fatalf("error requesting matrices: %v", err)
	}

	for i, row := range expectedDistances {
		for j, col := range row {
			if c := disMeasurer.Cost(i, j); c != col {
				t.Errorf("want: %f; got: %f\n", col, c)
			}
		}
	}

	for i, row := range expectedDurations {
		for j, col := range row {
			if c := durMeasurer.Cost(i, j); c != col {
				t.Errorf("want: %f; got: %f\n", col, c)
			}
		}
	}
}

func TestDistanceDurationMatricesWithParallelQueries(t *testing.T) {
	osrmHost := os.Getenv(hostEnv)
	if osrmHost == "" {
		t.Skip(hostEnvNotSetMsg)
	}

	c := osrm.DefaultClient(osrmHost, true)
	disMeasurer, durMeasurer, err := osrm.DistanceDurationMatrices(c, p, 2)
	if err != nil {
		t.Fatalf("error requesting matrices: %v", err)
	}

	for i, row := range expectedDistances {
		for j, col := range row {
			if c := disMeasurer.Cost(i, j); c != col {
				t.Errorf("want: %f; got: %f\n", col, c)
			}
		}
	}

	for i, row := range expectedDurations {
		for j, col := range row {
			if c := durMeasurer.Cost(i, j); c != col {
				t.Errorf("want: %f; got: %f\n", col, c)
			}
		}
	}
}
