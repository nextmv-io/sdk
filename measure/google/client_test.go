package google

import (
	"encoding/base64"
	"fmt"
	"os"
	"reflect"
	"testing"

	"googlemaps.github.io/maps"
)

// TestSplit tests the split function to assert that the complete list of
// origins and destinations is broken into proper requests that fulfill the
// usage limits.
func TestSplit(t *testing.T) {
	type expected struct {
		originsStart      int
		numOrigins        int
		destinationsStart int
		numDestinations   int
	}
	type testCase struct {
		r        *maps.DistanceMatrixRequest
		expected []expected
	}
	tests := []testCase{
		// Case: 30 origins and 8 destinations produce 3 requests to accommodate
		// 240 elements: 100, 100 and 40.
		{
			r: &maps.DistanceMatrixRequest{
				Origins: []string{
					"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
					"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
					"21", "22", "23", "24", "25", "26", "27", "28", "29", "30",
				},
				Destinations: []string{
					"1", "2", "3", "4", "5", "6", "7", "8",
				},
			},
			expected: []expected{
				{
					originsStart:      0,
					numOrigins:        25,
					destinationsStart: 0,
					numDestinations:   4,
				},
				{
					originsStart:      0,
					numOrigins:        25,
					destinationsStart: 4,
					numDestinations:   4,
				},
				{
					originsStart:      25,
					numOrigins:        5,
					destinationsStart: 0,
					numDestinations:   8,
				},
			},
		},
		// Case: 8 origins and 30 destinations produce 3 requests to accommodate
		// 240 elements: 96, 96 and 48.
		{
			r: &maps.DistanceMatrixRequest{
				Origins: []string{
					"1", "2", "3", "4", "5", "6", "7", "8",
				},
				Destinations: []string{
					"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
					"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
					"21", "22", "23", "24", "25", "26", "27", "28", "29", "30",
				},
			},
			expected: []expected{
				{
					originsStart:      0,
					numOrigins:        8,
					destinationsStart: 0,
					numDestinations:   12,
				},
				{
					originsStart:      0,
					numOrigins:        8,
					destinationsStart: 12,
					numDestinations:   12,
				},
				{
					originsStart:      0,
					numOrigins:        8,
					destinationsStart: 24,
					numDestinations:   6,
				},
			},
		},
		// Case: 30 origins and 30 destinations produce 10 requests to
		// accommodate 900 elements: 100, 100, 100, 100, 100, 100, 100, 50, 100
		// and 50.
		{
			r: &maps.DistanceMatrixRequest{
				Origins: []string{
					"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
					"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
					"21", "22", "23", "24", "25", "26", "27", "28", "29", "30",
				},
				Destinations: []string{
					"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
					"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
					"21", "22", "23", "24", "25", "26", "27", "28", "29", "30",
				},
			},
			expected: []expected{
				{
					originsStart:      0,
					numOrigins:        25,
					destinationsStart: 0,
					numDestinations:   4,
				},
				{
					originsStart:      0,
					numOrigins:        25,
					destinationsStart: 4,
					numDestinations:   4,
				},
				{
					originsStart:      0,
					numOrigins:        25,
					destinationsStart: 8,
					numDestinations:   4,
				},
				{
					originsStart:      0,
					numOrigins:        25,
					destinationsStart: 12,
					numDestinations:   4,
				},
				{
					originsStart:      0,
					numOrigins:        25,
					destinationsStart: 16,
					numDestinations:   4,
				},
				{
					originsStart:      0,
					numOrigins:        25,
					destinationsStart: 20,
					numDestinations:   4,
				},
				{
					originsStart:      0,
					numOrigins:        25,
					destinationsStart: 24,
					numDestinations:   4,
				},
				{
					originsStart:      0,
					numOrigins:        25,
					destinationsStart: 28,
					numDestinations:   2,
				},
				{
					originsStart:      25,
					numOrigins:        5,
					destinationsStart: 0,
					numDestinations:   20,
				},
				{
					originsStart:      25,
					numOrigins:        5,
					destinationsStart: 20,
					numDestinations:   10,
				},
			},
		},
	}

	// Convenient function for comparing fields in the expected struct.
	compare := func(i, j, k, got, want int) {
		if got != want {
			t.Errorf(
				"test %d, request %d, element %d: got %d, want %d",
				i, j, k, got, want,
			)
		}
	}

	// Compare all the expected fields by test.
	for i, test := range tests {
		requests := split(test.r)
		for j, request := range requests {
			comparisons := []struct {
				got  int
				want int
			}{
				// Origins start index.
				{
					got:  request.origins.start,
					want: test.expected[j].originsStart,
				},
				// Origins length.
				{
					got:  request.origins.num,
					want: test.expected[j].numOrigins,
				},
				// Destinations start index.
				{
					got:  request.destinations.start,
					want: test.expected[j].destinationsStart,
				},
				// Destinations length.
				{
					got:  request.destinations.num,
					want: test.expected[j].numDestinations,
				},
			}
			// Execute comparisons of the fields in the expected struct.
			for k, comparison := range comparisons {
				compare(i, j, k, comparison.got, comparison.want)
			}
		}
	}
}

// TestGroup tests the group function to assert that several requests are
// grouped together respecting the maxElementsPerSecond limit.
func TestGroup(t *testing.T) {
	type testCase struct {
		request  matrixRequest
		num      int
		expected int
	}

	tests := []testCase{
		// Case: 22 100-element-long requests result in 3 groups of requests.
		{
			request: matrixRequest{
				origins:      reference{num: 25},
				destinations: reference{num: 4},
			},
			num:      22,
			expected: 3,
		},
		// Case: 20 50-element-long requests result in 1 group of requests.
		{
			request: matrixRequest{
				origins:      reference{num: 25},
				destinations: reference{num: 2},
			},
			num:      20,
			expected: 1,
		},
	}

	for i, test := range tests {
		// Create large number of requests.
		requests := make([]matrixRequest, test.num)
		for r := range requests {
			requests[r] = test.request
		}

		// Assert number of groups is correct.
		got := len(group(requests))
		want := test.expected
		if got != want {
			t.Errorf("test %d: got %d, want %d", i, got, want)
		}
	}
}

func TestPolylines(t *testing.T) {
	type testCase struct {
		points           [][2]float64
		expectedPolyFull string
		expectedPolyLegs []string
	}

	tests := []testCase{
		{
			points: [][2]float64{
				{-74.028297, 4.875835},
				{-74.046965, 4.872842},
				{-74.041763, 4.885648},
			},
			// these lines are base64 encoded
			expectedPolyFull: "c2h3XHxzeWJNQFJIRlxQSExCVkdYU1RJRl1KYUBEV0pbVE9U" +
				"S0dxQFVRWnNAcEFHRE1SfUBiQkdSQ1RhQHxAYUFsQlNaXWpAaUBuQV9BdkJjQWhCS" +
				"0R5QFNnQEdBP2VGaUJjRXVBY0BNaUFmRGlBcERjQHZBQ0JBQkFGRExMRG5DakJWVFJ" +
				"gQGZAZkFMYEBAVEFgQFVsQFlaa0BWeUBSU0JbQ1VHcUBbbUFvQG9BbUBZRW1AP2F" +
				"EYUB1Q11xQUdhR1l5Q1V5Q1lvQllxQFFpRXNBfUd3QmVCd0BPP2FAS3NAWVNASUpC" +
				"TEhKTEZiQmBASkpASmJBWnRDdkBwR3BCfkJwQGZFZEBqRFp2Rlp+QkZUQGpCVG5B" +
				"XHBCcEB4RHJBbEBUaUBiQlFeXWJBWX5ASUpHQnNASUlGSVJ9QGpDU1xTUlFWV3BA" +
				"cUJsR3FAekI=",
			expectedPolyLegs: []string{
				"c2h3XHxzeWJNP0hARD9CQD9GRkxITkZEREJGQko/SkFKRUxHSktIRURDQEtG" +
					"UUJRQk9ATURJRE9IS0pPVD8/S0dxQFU/P1Fac0BwQT8/R0RFREdMW2pAYUB2" +
					"QEVKQUY/TENGQUBDRkVGVWpAY0B6QF1wQEdMS0xNUE9YV2pAUWJAcUB8QU1Y" +
					"fUB+QUVIQ0BBQEFAQz9BP0U/R0FpQFFbR0NBQT9FQEE/Pz9jRmlCPz9BP0M/" +
					"QUFNRW9EbUE/P2NATWlBZkRzQGBDVW5AX0B0QT8/Q0BDQj9AQT8/QD9AP0B" +
					"BPz9AP0A/QEA/P0A/QD9AQEA/QEA/P0BAPz9AQD8/QEA/QD9AQEA/QD9sQGJ" +
					"AYEBYXFRgQFZOSkZIUF5AQFRmQFBeTGBAQFRBSD9WR1JNWENEVVRTSldKW0" +
					"hNQkVASUJFP01CQT9ZQ1VHcUBbRUNnQWtAb0FtQFlFUUFbQD8/UUVZQ2tAS" +
					"T8/aUFNbUJXZ0BFY0BFbUBBa0JJdUNPdUBHY0JNcUFLZ0FNc0BJe0BPW0dVS" +
					"XtEbUFNRUlDZ0VxQXNAV1dJTUdZTWtAV1FJQ0E/P0VAQT9DP0k/V0ttQFdF" +
					"QUNBQz9FQEVARUJBQj9AQUA/QD9EQEJAQEJGREJEQkZCUkZ8QFJMQkJAQkJ" +
					"AP0JCQEJARD9EPz9SRm5AUl5MbkBMVkZsQFJuQ3hAYEN2QHpAWGJBVlxGXk" +
					"RiQUpkQUpqRFp2RlpmQUR2QEBIQEo/cEBGeEBMZEBKaEBQUkZ8QWhAbkFiQ" +
					"GZBYEBMRFJGXk5MRFhIbEBSXExmQFJIQHBCZEBoQEhgQEhiQEZIQmBARFBC" +
					"TkBeRExAaEVeXEJ6REg/P0tgQEFAR1RBSEdaWXZAT2RAQUJRZEBDSkNGQ0Z" +
					"FSkFARUhJTEVMQUBLUmNAaEFHUlNqQElWVWpAY0BiQVtwQHtAYkJjQGRBa0" +
					"BuQU1USU5FSElOT2ZAQURXdEBFREVCQ0JFQlNGX0BKRUBFP0E/QT9DP0E/Q" +
					"T9BP3FBa0BVSWFBWUlBQ0FBP0E/Qz9DQENAQUBBQEFAQUBnQHpBSVZZZkFx" +
					"QHRCZUFyRD8/QEBAP0BAQD9AP0BAQD9AQEA/QEBAP0A/QEBAQEA/QEBAP0A" +
					"/QEBAP0BAQD9AQEA/QD9AQEA/P0BAP0BAQD9AQEA/QD9AQEA/QEBAP0BAQD" +
					"9APz9AQD9AQEA/QEBAP0A/QEBAP0BAQD9AQEA/QD8/QEA/QEBAP0BAQD9AP" +
					"0BAQD9AQEA/QEBAP0A/P0BAP0BAQD9AQEA/QD9AQEA/QD9AQEA/QEBAP0A/" +
					"P0BAP0A/QEBAP0A/QEBAP0BAQD9wQFZIQnZBbkBkQWZAWk5kQWZAekBgQG" +
					"ZAWHhAZkBGRGBBcEBOSlxWfkBwQERCRkROSk5KeEFiQXBBfkBMSGJCekFU" +
					"VGRAXGZAWlBKVk5kQFZKRnpAXkhCXk5CQHpAXEA/QD8/QEA/QEBAP0A/P0" +
					"BAP0A/P0BAP0A/P0BAP0A/P0BAP0A/QEBAP0BAQD9AQEA/QD8/QEA/QD8/" +
					"QEA/QD8/QEA/QD8/QEA/QD9AQEA/QEBAP0BAQD9APz9AQD9oQE5SSERAWk" +
					"pMRD8/XmlBXmlBWHVAVnVAZEBnQVJnQFBfQGRAaUFMW25Ad0FQZUBOZUBCR" +
					"T9DdkB3Qj8/SFk=",
				"dXd2XHxnfWJNSVg/P3dAdkI/QkNET2RAUWRAb0B2QU1aZUBoQVFeU2ZAZU" +
					"BmQVd0QFl0QF9AaEFfQGhBPz9NRVtLRUFTSWlAT0E/P0FBP0E/QUFBP0FB" +
					"QT9BQUE/QT8/QUE/QT8/QUE/QT8/QUE/QT8/QUE/QT9BQUE/QUFBP0FBQT" +
					"9BPz9BQT9BPz9BQT9BPz9BQT9BPz9BQT9BP0FBQT8/QUE/QT97QF1DQV9" +
					"AT0lDe0BfQEtHZUBXV09RS2dAW2VAXVVVY0J7QU1JcUFfQXlBY0FPS09L" +
					"R0VFQ19BcUBdV09LYUFxQEdFeUBnQGdAWXtAYUBlQWdAW09lQWdAd0FvQ" +
					"ElDcUBXQT9BQUE/QUFBP0E/QUFBP0E/P0FBP0E/QUFBP0FBQT9BP0FBQT" +
					"9BP0FBQT9BQUE/P0FBP0E/QUFBP0FBQT9BQUE/QT9BQUE/QUFBPz9BQT9" +
					"BP0FBQT9BQUE/QUFBP0E/QUFBP0FBQT8/QUE/QT9BQUE/QUFBP0FBQT9B" +
					"P0FBQT9BQUE/P0FBP0FBQT9BP0FBQT9BQUE/QUFBP0E/QUFBP0FBQUFBP" +
					"0E/QUFBP0FBQT9BQUE/QT9BQUE/QUFBP0FBQT9BQUE/QUFBP0E/QUFBP0" +
					"FBQUFBP0E/QUFBP0FBQT9BQUE/QT9BQUE/QUFBP0FBQT9BQUE/QUFDP0F" +
					"BQT9BQUE/QUFBP0FBQz9BQUFBQT9BQUE/QUFDP0FBQT9BQUE/QUFBP0FB" +
					"Qz9BQUE/QUFBP0FBQT9BQUM/QUFBP0FBQT9BQUM/QUFBP0FBQT9BQUE/Q" +
					"UFDP0FBQT9BQUFBQT9BQUE/QUFBP0FBQT9BQUE/QUFBQUE/QUFBP0FBQT" +
					"9BQUE/QUFjQWdAb0BfQFdPa0NhQn1AZ0BPSz8/Vn1AVGFBckBlQ05lQFh" +
					"fQVRxQEZT",
			},
		},
	}
	for i, test := range tests {
		coords := make([]string, len(test.points))
		for p, point := range test.points {
			coords[p] = fmt.Sprintf("%f,%f", point[1], point[0])
		}

		apiKey := os.Getenv("GOOGLE_API_KEY")

		if apiKey == "" {
			t.Skip()
		}
		c, err := maps.NewClient(maps.WithAPIKey(apiKey))
		if err != nil {
			panic(err)
		}

		rPoly := &maps.DirectionsRequest{
			Origin:      coords[0],
			Destination: coords[len(coords)-1],
			Waypoints:   coords[1 : len(coords)-1],
		}

		fullPoly, polyLegs, err := Polylines(c, rPoly)
		if err != nil {
			panic(err)
		}

		// Assert polylines are correct.
		gotPolyFull := base64.StdEncoding.EncodeToString([]byte(fullPoly))
		wantPolyFull := test.expectedPolyFull
		if gotPolyFull != wantPolyFull {
			t.Errorf("test %d: got %s, want %s", i, gotPolyFull, wantPolyFull)
		}

		encodedPolyLegs := make([]string, len(polyLegs))
		for i, p := range polyLegs {
			encodedPolyLegs[i] = base64.StdEncoding.EncodeToString([]byte(p))
		}

		gotPolyLegs := encodedPolyLegs
		wantPolyLegs := test.expectedPolyLegs
		if !reflect.DeepEqual(gotPolyLegs, wantPolyLegs) {
			t.Errorf("test %d: got %s, want %s", i, gotPolyLegs, wantPolyLegs)
		}
	}
}
