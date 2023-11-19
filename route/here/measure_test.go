package here

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/nextmv-io/sdk/route"
)

const (
	apiKey         string = "foo"
	statusEndpoint string = "/status"
)

// This isn't the best test in the world because it relies on timing, so it
// doesn't test that every individual cancellation point handles cancellations
// properly - but that is very hard to do so this is good enough for now.
func TestCancellation(t *testing.T) {
	inPoints := []point{
		{Lon: 1.0, Lat: 0.1},
		{Lon: 2.0, Lat: 0.2},
		{Lon: 3.0, Lat: 0.3},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	s := newMockServer(ctx, apiKey, []requestSpec{
		{
			Endpoint:      "matrix",
			ExpectedAsync: true,

			ExpectedRequest: matrixRequest{
				Origins:          inPoints,
				MatrixAttributes: []string{"distances"},
				RegionDefinition: regionDefinition{
					Type: "autoCircle",
				},
			},
			BuildResponse: asyncMatrixSuccess(),
		},
		{
			Endpoint:      "status",
			BuildResponse: internalServerError(),
			latency:       60,
		},
		{
			Endpoint:      "status",
			BuildResponse: internalServerError(),
			latency:       60,
		},
		{
			Endpoint:      "status",
			BuildResponse: internalServerError(),
			latency:       60,
		},
	})
	cli := NewClient(apiKey)
	cli.(*client).retries = 100
	cli.(*client).schemeHost = s.URL
	cli.(*client).maxSyncPoints = 1

	points := []route.Point{{1.0, 0.1}, {2.0, 0.2}, {3.0, 0.3}}
	_, err := cli.DistanceMatrix(ctx, points)
	if err == nil {
		t.Errorf("expected context cancellation error: %v", err)
	}
}

func TestClientRedirect(t *testing.T) {
	apiKey := "foo"
	inPoints := []point{
		{Lon: 1.0, Lat: 0.1},
		{Lon: 2.0, Lat: 0.2},
		{Lon: 3.0, Lat: 0.3},
	}
	responseMatrix := []int{
		0, 1, 2,
		3, 0, 4,
		5, 6, 0,
	}

	serverHostName := "127.0.0.1"
	testHostName := "hereapi.com"

	tests := []struct {
		description string
		requests    []requestSpec
	}{
		{
			description: "ignores redirect for relative URL",
			requests: []requestSpec{
				{
					Endpoint:      "matrix",
					ExpectedAsync: true,
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
					},
					BuildResponse: asyncMatrixSuccess(),
				},
				{
					Endpoint:      "status",
					BuildResponse: statusRedirect(responseStatusComplete),
					redirectTo:    "/some_path_that_wont_be_followed",
				},
				{
					Endpoint: "result",
					BuildResponse: resultSuccess(matrixResponse{
						Matrix: matrix{
							Distances: responseMatrix,
						},
					}),
				},
			},
		},
		{
			description: "ignores redirect for full URL",
			requests: []requestSpec{
				{
					Endpoint:      "matrix",
					ExpectedAsync: true,
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
					},
					BuildResponse: asyncMatrixSuccess(),
				},
				{
					Endpoint:      "status",
					BuildResponse: statusRedirect(responseStatusComplete),
					redirectTo:    fmt.Sprintf("https://%s/some_path", serverHostName),
				},
				{
					Endpoint: "result",
					BuildResponse: resultSuccess(matrixResponse{
						Matrix: matrix{
							Distances: responseMatrix,
						},
					}),
				},
			},
		},
		{
			description: "ignores redirects for hostnames greedily",
			requests: []requestSpec{
				{
					Endpoint:      "matrix",
					ExpectedAsync: true,
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
					},
					BuildResponse: asyncMatrixSuccess(),
				},
				{
					Endpoint:      "status",
					BuildResponse: statusRedirect(responseStatusComplete),
					redirectTo: fmt.Sprintf(
						"https://foo.bar.baz.%s/some_path", testHostName),
				},
				{
					Endpoint: "result",
					BuildResponse: resultSuccess(matrixResponse{
						Matrix: matrix{
							Distances: responseMatrix,
						},
					}),
				},
			},
		},
	}

	for _, test := range tests {
		s := newMockServer(context.Background(), apiKey, test.requests)
		u, _ := url.Parse(s.URL)
		serverHostName = u.Hostname()

		cli := NewClient(apiKey, WithDenyRedirectPolicy(
			serverHostName, testHostName))
		cli.(*client).schemeHost = s.URL
		cli.(*client).maxSyncPoints = 1

		points := []route.Point{{1.0, 0.1}, {2.0, 0.2}, {3.0, 0.3}}
		_, err := cli.DistanceMatrix(context.Background(), points)
		if err != nil {
			t.Errorf(
				"redirect test %s failed: unexpected error: %v",
				test.description,
				err,
			)
		}
	}
}

func TestHereMeasures(t *testing.T) {
	distances := func(
		cli Client,
		ps []route.Point,
		opts ...MatrixOption,
	) (route.ByIndex, error) {
		return cli.DistanceMatrix(context.Background(), ps, opts...)
	}
	durations := func(
		cli Client,
		ps []route.Point,
		opts ...MatrixOption,
	) (route.ByIndex, error) {
		return cli.DurationMatrix(context.Background(), ps, opts...)
	}
	selectDistances := func(
		cli Client,
		ps []route.Point,
		opts ...MatrixOption,
	) (route.ByIndex, error) {
		distances, _, err := cli.DistanceDurationMatrices(
			context.Background(), ps, opts...)
		return distances, err
	}
	selectDurations := func(
		cli Client,
		ps []route.Point,
		opts ...MatrixOption,
	) (route.ByIndex, error) {
		_, durations, err := cli.DistanceDurationMatrices(
			context.Background(), ps, opts...)
		return durations, err
	}
	type test struct {
		description   string
		points        []route.Point
		maxSyncPoints int
		retries       int
		expectErr     bool
		getMatrix     func(
			Client,
			[]route.Point,
			...MatrixOption,
		) (route.ByIndex, error)
		expectedMatrix [][]float64
		requests       []requestSpec
		opts           []MatrixOption
	}

	apiKey := `foo`
	responseMatrix := []int{
		0, 1, 2,
		3, 0, 4,
		5, 6, 0,
	}
	points := []route.Point{{1.0, 0.1}, {2.0, 0.2}, {3.0, 0.3}}
	pointsWithEmptyVals := []route.Point{
		{1.0, 0.1}, {2.0, 0.2}, {}, {3.0, 0.3}, {},
	}
	inPoints := []point{
		{Lon: 1.0, Lat: 0.1},
		{Lon: 2.0, Lat: 0.2},
		{Lon: 3.0, Lat: 0.3},
	}
	expectedMatrix := [][]float64{{0, 1.0, 2.0}, {3.0, 0, 4.0}, {5.0, 6.0, 0.0}}
	expectedMatrixWithEmptyVals := [][]float64{
		{0, 1, 0, 2, 0},
		{3, 0, 0, 4, 0},
		{0, 0, 0, 0, 0},
		{5, 6, 0, 0, 0},
		{0, 0, 0, 0, 0},
	}

	truckProfile := Truck{
		ShippedHazardousGoods: []HazardousGood{
			Explosive,
			Poison,
		},
		GrossWeight:    100,
		WeightPerAxle:  25,
		Height:         1000,
		Width:          20,
		Length:         10,
		TunnelCategory: TunnelCategoryC,
		AxleCount:      4,
		Type:           TruckTypeTractor,
		WeightPerAxleGroup: &WeightPerAxleGroup{
			Tandem: 6,
			Triple: 7,
		},
	}

	tests := []test{
		{
			description:    "sync calculation with empty locations (distances)",
			points:         pointsWithEmptyVals,
			getMatrix:      distances,
			expectedMatrix: expectedMatrixWithEmptyVals,
			requests: []requestSpec{
				{
					Endpoint: "matrix",
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
					},
					BuildResponse: syncMatrixSuccess(matrixResponse{
						Matrix: matrix{
							Distances: responseMatrix,
						},
					}),
				},
			},
			maxSyncPoints: 100,
		},
		{
			description:    "sync calculation (distances)",
			points:         points,
			getMatrix:      distances,
			expectedMatrix: expectedMatrix,
			requests: []requestSpec{
				{
					Endpoint: "matrix",
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
					},
					BuildResponse: syncMatrixSuccess(matrixResponse{
						Matrix: matrix{
							Distances: responseMatrix,
						},
					}),
				},
			},
			maxSyncPoints: 100,
		},
		{
			description:    "sync calculation (durations)",
			points:         points,
			getMatrix:      durations,
			expectedMatrix: expectedMatrix,
			requests: []requestSpec{
				{
					Endpoint: "matrix",
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"travelTimes"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
					},
					BuildResponse: syncMatrixSuccess(matrixResponse{
						Matrix: matrix{
							TravelTimes: responseMatrix,
						},
					}),
				},
			},
			maxSyncPoints: 100,
		},
		{
			description:    "sync calculation (distances from distances + durations)",
			points:         points,
			getMatrix:      selectDistances,
			expectedMatrix: expectedMatrix,
			requests: []requestSpec{
				{
					Endpoint: "matrix",
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances", "travelTimes"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
					},
					BuildResponse: syncMatrixSuccess(matrixResponse{
						Matrix: matrix{
							Distances:   responseMatrix,
							TravelTimes: responseMatrix,
						},
					}),
				},
			},
			maxSyncPoints: 100,
		},
		{
			description:    "sync calculation (durations from distances + durations)",
			points:         points,
			getMatrix:      selectDurations,
			expectedMatrix: expectedMatrix,
			requests: []requestSpec{
				{
					Endpoint: "matrix",
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances", "travelTimes"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
					},
					BuildResponse: syncMatrixSuccess(matrixResponse{
						Matrix: matrix{
							Distances:   responseMatrix,
							TravelTimes: responseMatrix,
						},
					}),
				},
			},
			maxSyncPoints: 100,
		},
		{
			description:    "sync with error",
			points:         points,
			getMatrix:      selectDurations,
			expectedMatrix: expectedMatrix,
			expectErr:      true,
			requests: []requestSpec{
				{
					Endpoint: "matrix",
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances", "travelTimes"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
					},
					BuildResponse: internalServerError(),
				},
			},
			maxSyncPoints: 100,
		},
		{
			description:    "async calculation with empty locations (durations)",
			points:         pointsWithEmptyVals,
			getMatrix:      durations,
			expectedMatrix: expectedMatrixWithEmptyVals,
			requests: []requestSpec{
				{
					Endpoint:      "matrix",
					ExpectedAsync: true,
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"travelTimes"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
					},
					BuildResponse: asyncMatrixSuccess(),
				},
				{
					Endpoint:      "status",
					BuildResponse: statusSuccess(),
				},
				{
					Endpoint: "result",
					BuildResponse: resultSuccess(matrixResponse{
						Matrix: matrix{
							TravelTimes: responseMatrix,
						},
					}),
				},
			},
		},
		{
			description:    "async calculation (distances)",
			points:         points,
			getMatrix:      distances,
			expectedMatrix: expectedMatrix,
			requests: []requestSpec{
				{
					Endpoint:      "matrix",
					ExpectedAsync: true,
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
					},
					BuildResponse: asyncMatrixSuccess(),
				},
				{
					Endpoint:      "status",
					BuildResponse: statusSuccess(),
				},
				{
					Endpoint: "result",
					BuildResponse: resultSuccess(matrixResponse{
						Matrix: matrix{
							Distances: responseMatrix,
						},
					}),
				},
			},
		},
		{
			description:    "async calculation (durations)",
			points:         points,
			getMatrix:      durations,
			expectedMatrix: expectedMatrix,
			requests: []requestSpec{
				{
					Endpoint:      "matrix",
					ExpectedAsync: true,
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"travelTimes"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
					},
					BuildResponse: asyncMatrixSuccess(),
				},
				{
					Endpoint:      "status",
					BuildResponse: statusSuccess(),
				},
				{
					Endpoint: "result",
					BuildResponse: resultSuccess(matrixResponse{
						Matrix: matrix{
							TravelTimes: responseMatrix,
						},
					}),
				},
			},
		},
		{
			description: "async calculation (distances from distances " +
				"+ durations)",
			points:         points,
			getMatrix:      selectDistances,
			expectedMatrix: expectedMatrix,
			requests: []requestSpec{
				{
					Endpoint:      "matrix",
					ExpectedAsync: true,
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances", "travelTimes"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
					},
					BuildResponse: asyncMatrixSuccess(),
				},
				{
					Endpoint:      "status",
					BuildResponse: statusSuccess(),
				},
				{
					Endpoint: "result",
					BuildResponse: resultSuccess(matrixResponse{
						Matrix: matrix{
							Distances:   responseMatrix,
							TravelTimes: responseMatrix,
						},
					}),
				},
			},
		},
		{
			description: "async calculation (durations from distances " +
				"+ durations)",
			points:         points,
			getMatrix:      selectDurations,
			expectedMatrix: expectedMatrix,
			requests: []requestSpec{
				{
					Endpoint:      "matrix",
					ExpectedAsync: true,
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances", "travelTimes"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
					},
					BuildResponse: asyncMatrixSuccess(),
				},
				{
					Endpoint:      "status",
					BuildResponse: statusSuccess(),
				},
				{
					Endpoint: "result",
					BuildResponse: resultSuccess(matrixResponse{
						Matrix: matrix{
							Distances:   responseMatrix,
							TravelTimes: responseMatrix,
						},
					}),
				},
			},
		},
		{
			description:    "async calculation with status retries",
			points:         points,
			getMatrix:      distances,
			expectedMatrix: expectedMatrix,
			requests: []requestSpec{
				{
					Endpoint:      "matrix",
					ExpectedAsync: true,

					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
					},
					BuildResponse: asyncMatrixSuccess(),
				},
				{
					Endpoint:      "status",
					BuildResponse: internalServerError(),
				},
				{
					Endpoint:      "status",
					BuildResponse: statusSuccess(),
				},
				{
					Endpoint: "result",
					BuildResponse: resultSuccess(matrixResponse{
						Matrix: matrix{
							Distances: responseMatrix,
						},
					}),
				},
			},
		},
		{
			description:    "async calculation with failure from status",
			points:         points,
			getMatrix:      distances,
			expectErr:      true,
			expectedMatrix: expectedMatrix,
			retries:        1,
			requests: []requestSpec{
				{
					Endpoint:      "matrix",
					ExpectedAsync: true,

					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
					},
					BuildResponse: asyncMatrixSuccess(),
				},
				{
					Endpoint:      "status",
					BuildResponse: internalServerError(),
				},
				{
					Endpoint:      "status",
					BuildResponse: internalServerError(),
				},
			},
		},
		{
			description:    "async calculation with 404 from status",
			points:         points,
			getMatrix:      distances,
			expectErr:      true,
			expectedMatrix: expectedMatrix,
			requests: []requestSpec{
				{
					Endpoint:      "matrix",
					ExpectedAsync: true,

					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
					},
					BuildResponse: asyncMatrixSuccess(),
				},
				{
					Endpoint:      "status",
					BuildResponse: notFound(),
				},
			},
		},
		{
			description:    "async calculation with result retries",
			points:         points,
			getMatrix:      distances,
			expectedMatrix: expectedMatrix,
			requests: []requestSpec{
				{
					Endpoint:      "matrix",
					ExpectedAsync: true,
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
					},
					BuildResponse: asyncMatrixSuccess(),
				},
				{
					Endpoint:      "status",
					BuildResponse: statusSuccess(),
				},
				{
					Endpoint:      "result",
					BuildResponse: internalServerError(),
				},
				{
					Endpoint: "result",
					BuildResponse: resultSuccess(matrixResponse{
						Matrix: matrix{
							Distances: responseMatrix,
						},
					}),
				},
			},
		},
		{
			description:    "async calculation with failure from result",
			points:         points,
			getMatrix:      distances,
			expectErr:      true,
			expectedMatrix: expectedMatrix,
			retries:        1,
			requests: []requestSpec{
				{
					Endpoint:      "matrix",
					ExpectedAsync: true,

					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
					},
					BuildResponse: asyncMatrixSuccess(),
				},
				{
					Endpoint:      "status",
					BuildResponse: statusSuccess(),
				},
				{
					Endpoint:      "result",
					BuildResponse: internalServerError(),
				},
				{
					Endpoint:      "result",
					BuildResponse: internalServerError(),
				},
			},
		},
		{
			description:    "async calculation with failure from result",
			points:         points,
			getMatrix:      distances,
			expectErr:      true,
			expectedMatrix: expectedMatrix,
			requests: []requestSpec{
				{
					Endpoint:      "matrix",
					ExpectedAsync: true,

					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
					},
					BuildResponse: asyncMatrixSuccess(),
				},
				{
					Endpoint:      "status",
					BuildResponse: statusSuccess(),
				},
				{
					Endpoint:      "result",
					BuildResponse: notFound(),
				},
			},
		},
		// The following tests test various options that HERE supports - I just
		// randomly alternated between sync and async patterns and different
		// matrix types between the tests to keep things interesting
		// (I think testing all modes for all options would be overkill)
		{
			description:    "departure time (any)",
			points:         points,
			getMatrix:      selectDurations,
			expectedMatrix: expectedMatrix,
			opts: []MatrixOption{
				WithDepartureTime(time.Time{}),
			},
			requests: []requestSpec{
				{
					Endpoint:      "matrix",
					ExpectedAsync: true,
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances", "travelTimes"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
						DepartureTime: "any",
					},
					BuildResponse: asyncMatrixSuccess(),
				},
				{
					Endpoint:      "status",
					BuildResponse: statusSuccess(),
				},
				{
					Endpoint: "result",
					BuildResponse: resultSuccess(matrixResponse{
						Matrix: matrix{
							Distances:   responseMatrix,
							TravelTimes: responseMatrix,
						},
					}),
				},
			},
		},
		{
			description:    "departure time (specific time)",
			points:         pointsWithEmptyVals,
			getMatrix:      distances,
			expectedMatrix: expectedMatrixWithEmptyVals,
			opts: []MatrixOption{
				WithDepartureTime(
					time.Date(2021, 12, 10, 9, 30, 0, 0, time.UTC)),
			},
			requests: []requestSpec{
				{
					Endpoint: "matrix",
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
						DepartureTime: "2021-12-10T09:30:00Z",
					},
					BuildResponse: syncMatrixSuccess(matrixResponse{
						Matrix: matrix{
							Distances: responseMatrix,
						},
					}),
				},
			},
			maxSyncPoints: 100,
		},
		{
			description:    "transport mode",
			points:         pointsWithEmptyVals,
			getMatrix:      distances,
			expectedMatrix: expectedMatrixWithEmptyVals,
			opts: []MatrixOption{
				WithTransportMode(TransportModeBicycle),
				WithDepartureTime(
					time.Date(2021, 12, 10, 9, 30, 0, 0, time.UTC)),
			},
			requests: []requestSpec{
				{
					Endpoint: "matrix",
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
						DepartureTime: "2021-12-10T09:30:00Z",
						TransportMode: TransportMode("bicycle"),
					},
					BuildResponse: syncMatrixSuccess(matrixResponse{
						Matrix: matrix{
							Distances: responseMatrix,
						},
					}),
				},
			},
			maxSyncPoints: 100,
		},
		{
			description:    "avoid features",
			points:         pointsWithEmptyVals,
			getMatrix:      distances,
			expectedMatrix: expectedMatrixWithEmptyVals,
			opts: []MatrixOption{
				WithTransportMode(TransportModeBicycle),
				WithAvoidFeatures([]Feature{Ferry, DirtRoad}),
			},
			requests: []requestSpec{
				{
					Endpoint: "matrix",
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
						TransportMode: "bicycle",
						Avoid: &avoid{
							Features: []Feature{Ferry, DirtRoad},
						},
					},
					BuildResponse: syncMatrixSuccess(matrixResponse{
						Matrix: matrix{
							Distances: responseMatrix,
						},
					}),
				},
			},
			maxSyncPoints: 100,
		},
		{
			description:    "avoid areas",
			points:         points,
			getMatrix:      durations,
			expectedMatrix: expectedMatrix,
			opts: []MatrixOption{
				WithDepartureTime(time.Time{}),
				WithAvoidAreas([]BoundingBox{
					{
						North: 0.8,
						South: 0.6,
						East:  0.7,
						West:  0.5,
					},
				}),
				WithAvoidAreas([]BoundingBox{
					{
						North: 4.8,
						South: 2.6,
						East:  3.7,
						West:  1.5,
					},
				}),
			},
			requests: []requestSpec{
				{
					Endpoint:      "matrix",
					ExpectedAsync: true,
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"travelTimes"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
						DepartureTime: "any",
						Avoid: &avoid{
							Areas: []area{
								{
									Type:  "boundingBox",
									North: 0.8,
									South: 0.6,
									East:  0.7,
									West:  0.5,
								},
								{
									Type:  "boundingBox",
									North: 4.8,
									South: 2.6,
									East:  3.7,
									West:  1.5,
								},
							},
						},
					},
					BuildResponse: asyncMatrixSuccess(),
				},
				{
					Endpoint:      "status",
					BuildResponse: statusSuccess(),
				},
				{
					Endpoint: "result",
					BuildResponse: resultSuccess(matrixResponse{
						Matrix: matrix{
							TravelTimes: responseMatrix,
						},
					}),
				},
			},
		},
		{
			description:    "truck profile",
			points:         pointsWithEmptyVals,
			getMatrix:      distances,
			expectedMatrix: expectedMatrixWithEmptyVals,
			opts: []MatrixOption{
				WithTransportMode(TransportModeTruck),
				WithTruckProfile(truckProfile),
			},
			requests: []requestSpec{
				{
					Endpoint: "matrix",
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
						TransportMode: TransportModeTruck,
						Truck:         &truckProfile,
					},
					BuildResponse: syncMatrixSuccess(matrixResponse{
						Matrix: matrix{
							Distances: responseMatrix,
						},
					}),
				},
			},
			maxSyncPoints: 100,
		},
		{
			description:    "scooter profile",
			points:         pointsWithEmptyVals,
			getMatrix:      distances,
			expectedMatrix: expectedMatrixWithEmptyVals,
			opts: []MatrixOption{
				WithTransportMode(TransportModeScooter),
				WithScooterProfile(Scooter{AllowHighway: true}),
			},
			requests: []requestSpec{
				{
					Endpoint: "matrix",
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
						TransportMode: TransportModeScooter,
						Scooter: &Scooter{
							AllowHighway: true,
						},
					},
					BuildResponse: syncMatrixSuccess(matrixResponse{
						Matrix: matrix{
							Distances: responseMatrix,
						},
					}),
				},
			},
			maxSyncPoints: 100,
		},
		{
			description:    "taxi profile",
			points:         pointsWithEmptyVals,
			getMatrix:      distances,
			expectedMatrix: expectedMatrixWithEmptyVals,
			opts: []MatrixOption{
				WithTransportMode(TransportModeTaxi),
				WithTaxiProfile(Taxi{
					AllowDriveThroughTaxiRoads: true,
				}),
			},
			requests: []requestSpec{
				{
					Endpoint: "matrix",
					ExpectedRequest: matrixRequest{
						Origins:          inPoints,
						MatrixAttributes: []string{"distances"},
						RegionDefinition: regionDefinition{
							Type: "autoCircle",
						},
						TransportMode: TransportModeTaxi,
						Taxi: &Taxi{
							AllowDriveThroughTaxiRoads: true,
						},
					},
					BuildResponse: syncMatrixSuccess(matrixResponse{
						Matrix: matrix{
							Distances: responseMatrix,
						},
					}),
				},
			},
			maxSyncPoints: 100,
		},
	}
	for i, test := range tests {
		s := newMockServer(context.Background(), apiKey, test.requests)
		cli := NewClient(apiKey)
		cli.(*client).maxSyncPoints = 1
		if test.maxSyncPoints > 0 {
			cli.(*client).maxSyncPoints = test.maxSyncPoints
		}
		if test.retries != 0 {
			cli.(*client).retries = test.retries
		}
		cli.(*client).schemeHost = s.URL

		m, err := test.getMatrix(cli, test.points, test.opts...)
		// Error from server is higher priority than error on the client-side
		if s.Error() != nil {
			t.Errorf(
				"[%d] %s: error from mock server: %v",
				i,
				test.description, s.Error(),
			)
			continue
		}
		if err != nil {
			if !test.expectErr {
				t.Errorf("[%d] %s: unexpected error: %v", i, test.description, err)
			}
			continue
		}

		if err == nil && test.expectErr {
			t.Errorf("[%d] %s: expected error but got none", i, test.description)
			continue
		}

		got := unpackMatrix(m, len(test.points))
		if diff := cmp.Diff(got, test.expectedMatrix); diff != "" {
			t.Errorf("[%d] %s: (-want, +got)\n%s", i, test.description, diff)
		}
	}
}

func unpackMatrix(m route.ByIndex, width int) [][]float64 {
	matrix := make([][]float64, width)
	for i := 0; i < width; i++ {
		matrix[i] = make([]float64, width)
		for j := 0; j < width; j++ {
			matrix[i][j] = m.Cost(i, j)
		}
	}
	return matrix
}

type mockServer struct {
	*httptest.Server
	err              error
	expectedRequests int
	requests         []requestSpec
}

func (m mockServer) Error() error {
	if m.err != nil {
		return m.err
	}
	if len(m.requests) != 0 {
		return fmt.Errorf(
			"received fewer requests than expected (%d)", m.expectedRequests)
	}
	return nil
}

func (m *mockServer) PopRequest() *requestSpec {
	if len(m.requests) == 0 {
		m.err = fmt.Errorf(
			"got more requests than expected (%d)", m.expectedRequests)
		return nil
	}
	nextRequest := m.requests[0]
	m.requests = m.requests[1:]
	return &nextRequest
}

type requestSpec struct {
	// This is a partial URL path that uniquely identifies the endpoint
	// Possible values are: matrix, status, and result
	Endpoint        string
	ExpectedAsync   bool
	ExpectedRequest any
	BuildResponse   buildResponse
	// latency will cause this request to wait some number of seconds
	latency int64
	// If set and using a redirect header, will set the `location` header
	// on the response
	redirectTo string
}

type buildResponse func(url string) (int, any)

func notFound() buildResponse {
	return func(url string) (int, any) {
		return http.StatusNotFound, nil
	}
}

func internalServerError() buildResponse {
	return func(url string) (int, any) {
		return http.StatusInternalServerError, nil
	}
}

func statusRedirect(status responseStatus) buildResponse {
	return func(url string) (int, any) {
		resp := statusResponse{
			StatusURL: url + statusEndpoint,
			Status:    status,
		}
		if status == responseStatusComplete {
			resp.ResultURL = url + "/result"
		}
		return http.StatusMovedPermanently, resp
	}
}

func syncMatrixSuccess(resp matrixResponse) buildResponse {
	return func(url string) (int, any) {
		return http.StatusOK, resp
	}
}

func asyncMatrixSuccess() buildResponse {
	return func(url string) (int, any) {
		return http.StatusAccepted, statusResponse{
			StatusURL: url + statusEndpoint,
		}
	}
}

func statusSuccess() buildResponse {
	return func(url string) (int, any) {
		resp := statusResponse{
			StatusURL: url + statusEndpoint,
			Status:    responseStatusComplete,
			ResultURL: url + "/result",
		}

		return http.StatusOK, resp
	}
}

func resultSuccess(response matrixResponse) buildResponse {
	return func(url string) (int, any) {
		return http.StatusOK, response
	}
}

func newMockServer( //nolint:gocyclo
	ctx context.Context,
	apiKey string,
	requests []requestSpec,
) *mockServer {
	m := &mockServer{
		requests:         requests,
		expectedRequests: len(requests),
	}

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextRequest := m.PopRequest()
			if nextRequest == nil {
				return
			}

			// Check if the api key is correct.
			if r.URL.Query().Get("apiKey") != apiKey {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			async := r.URL.Query().Get("async") == "true"
			if async && !nextRequest.ExpectedAsync {
				m.err = fmt.Errorf("got unexpected async request")
				return
			}
			if !async && nextRequest.ExpectedAsync {
				m.err = fmt.Errorf("expected async request")
				return
			}

			if !strings.Contains(r.URL.Path, nextRequest.Endpoint) {
				m.err = fmt.Errorf(
					"expected request type %s, got request to %s",
					nextRequest.Endpoint,
					r.URL.Path,
				)
				return
			}

			if nextRequest.ExpectedRequest != nil {
				request := make(map[string]any)
				if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
					m.err = fmt.Errorf("decoding request: %v", err)
					return
				}
				if diff := cmp.Diff(
					toMap(nextRequest.ExpectedRequest), request); diff != "" {
					m.err = fmt.Errorf("(-want, +got)\n%s", diff)
					return
				}
			}

			if nextRequest.latency != 0 {
				select {
				case <-ctx.Done():
				case <-time.After(time.Second * time.Duration(nextRequest.latency)):
				}
			}

			statusCode, response := nextRequest.BuildResponse(m.URL)

			if statusCode == http.StatusMovedPermanently {
				var redirect string
				if strings.HasPrefix(nextRequest.redirectTo, "/") {
					redirect = fmt.Sprintf("%s%s", m.URL, nextRequest.redirectTo)
				} else {
					redirect = nextRequest.redirectTo
				}
				w.Header().Set("Location", redirect)
			}
			w.WriteHeader(statusCode)
			err := json.NewEncoder(w).Encode(response)
			m.err = err
		}))

	m.Server = s

	return m
}

func toMap(x any) map[string]any {
	b, _ := json.Marshal(x)
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	return m
}
