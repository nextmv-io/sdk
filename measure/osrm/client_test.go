package osrm_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nextmv-io/sdk/measure"
	"github.com/nextmv-io/sdk/measure/osrm"
)

type testServer struct {
	s        *httptest.Server
	reqCount int
}

func newTestServer(t *testing.T, endpoint osrm.Endpoint) *testServer {
	ts := &testServer{}
	responseOk := ""
	if endpoint == osrm.TableEndpoint {
		responseOk = tableResponseOK
	} else {
		responseObject := osrm.RouteResponse{
			Code: "Ok",
			Routes: []osrm.Route{
				{
					Geometry: "mfp_I__vpAqJ`@wUrCa\\dCgGig@{DwW",
					Legs: []osrm.Leg{
						{
							Steps: []osrm.Step{
								{Geometry: "mfp_I__vpAWBQ@K@[BuBRgBLK@UBMMC?AA" +
									"KAe@FyBTC@E?IDKDA@K@]BUBSBA?E@E@A@KFUBK@mA" +
									"L{CZQ@qBRUBmAFc@@}@Fu@DG?a@B[@qAF}@JA?[D_" +
									"E`@SBO@ODA@UDA?]JC?uBNE?OAKA"},
								{Geometry: "yer_IcuupACa@AI]mCCUE[AK[iCWqB[{Bk" +
									"@sE_@_DAICSAOIm@AIQuACOQyAG[Gc@]wBw@aFKu@" +
									"y@oFCMAOIm@?K"},
								{Geometry: "}sr_IevwpA"},
							},
						},
					},
				},
			},
		}
		resp, err := json.Marshal(responseObject)
		if err != nil {
			t.Errorf("could not marshal response object, %v", err)
		}
		responseOk = string(resp)
	}
	ts.s = httptest.NewServer(
		http.AllowQuerySemicolons(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				ts.reqCount++
				_, err := io.WriteString(w, responseOk)
				if err != nil {
					t.Errorf("could not write resp: %v", err)
				}
			}),
		),
	)
	return ts
}

func newErrorServer(_ *testing.T, statusCode int, statusMsg string, causeErr bool) *testServer {
	ts := &testServer{}
	ts.s = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			ts.reqCount++
			if causeErr {
				panic("causing error")
			}
			w.WriteHeader(statusCode)
			_, _ = io.WriteString(w, statusMsg)
		}),
	)
	return ts
}

func TestCacheHit(t *testing.T) {
	ts := newTestServer(t, osrm.TableEndpoint)
	defer ts.s.Close()

	c := osrm.DefaultClient(ts.s.URL, true)

	_, err := c.Get(ts.s.URL)
	if err != nil {
		t.Errorf("get failed: %v", err)
	}

	_, err = c.Get(ts.s.URL)
	if err != nil {
		t.Errorf("get failed: %v", err)
	}

	_, err = c.Get(ts.s.URL)
	if err != nil {
		t.Errorf("get failed: %v", err)
	}

	if ts.reqCount != 1 {
		t.Errorf("want: 1; got: %v", ts.reqCount)
	}
}

func TestCacheMiss(t *testing.T) {
	ts := newTestServer(t, osrm.TableEndpoint)
	defer ts.s.Close()

	c := osrm.DefaultClient(ts.s.URL, false)
	_, err := c.Get(ts.s.URL)
	if err != nil {
		t.Errorf("get failed: %v", err)
	}

	_, err = c.Get(ts.s.URL)
	if err != nil {
		t.Errorf("get failed: %v", err)
	}

	_, err = c.Get(ts.s.URL)
	if err != nil {
		t.Errorf("get failed: %v", err)
	}

	if ts.reqCount != 3 {
		t.Errorf("want: 3; got: %v", ts.reqCount)
	}
}

func TestMatrixCall(t *testing.T) {
	ts := newTestServer(t, osrm.TableEndpoint)
	defer ts.s.Close()

	c := osrm.DefaultClient(ts.s.URL, true)
	m, err := osrm.DurationMatrix(c, []measure.Point{{0, 0}, {1, 1}}, 0)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	if v := m.Cost(0, 1); v != 17699.1 {
		t.Errorf("want: 0; got: %v", v)
	}
}

func TestTableErrorHandling(t *testing.T) {
	smallTable := []measure.Point{{0, 0}, {1, 1}}
	largeTable := make([]measure.Point, 101)
	for i := range largeTable {
		largeTable[i] = measure.Point{float64(i), float64(i)}
	}

	tests := map[string]struct {
		statusCode int
		statusMsg  string
		causeErr   bool
		isUserErr  bool
		table      []measure.Point
	}{
		"bad request": {
			statusCode: http.StatusBadRequest,
			statusMsg:  "Invalid segment",
			causeErr:   false,
			isUserErr:  true,
			table:      smallTable,
		},
		"internal server error": {
			statusCode: http.StatusInternalServerError,
			statusMsg:  "",
			causeErr:   false,
			isUserErr:  false,
			table:      smallTable,
		},
		"cause request err": {
			statusCode: 0, // not used
			statusMsg:  "",
			causeErr:   true,
			isUserErr:  false,
			table:      smallTable,
		},
		"cause multiple user errs": {
			statusCode: http.StatusBadRequest,
			statusMsg:  "Invalid segment",
			causeErr:   false,
			isUserErr:  true,
			table:      largeTable,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts := newErrorServer(t, tc.statusCode, tc.statusMsg, tc.causeErr)
			defer ts.s.Close()
			c := osrm.DefaultClient(ts.s.URL, true)
			_, _, err := c.Table(tc.table)
			if err == nil {
				t.Errorf("expected error, got nil")
			}
			uerr, ok := err.(osrm.Error)
			if !ok {
				t.Errorf("expected osrm.Error, got %T", err)
			}
			if uerr.IsInputError() != tc.isUserErr {
				t.Errorf("expected input error, got %T", err)
			}
			if tc.isUserErr && !strings.Contains(uerr.Error(), tc.statusMsg) {
				t.Errorf("want: %v; got: %v", tc.statusMsg, uerr.Error())
			}
		})
	}
}

func TestPolylineCall(t *testing.T) {
	ts := newTestServer(t, osrm.RouteEndpoint)
	defer ts.s.Close()

	c := osrm.DefaultClient(ts.s.URL, true)

	polyline, polyLegs, err := osrm.Polyline(
		c,
		[]measure.Point{
			{13.388860, 52.517037},
			{13.397634, 52.529407},
		},
	)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	println(polyline)
	println(polyLegs)
}

const tableResponseOK = `{
	"code": "Ok",
	"sources": [{
		"hint": "",
		"distance": 9.215349,
		"name": "",
		"location": [-105.050583, 39.762548]
	}, {
		"hint": "",
		"distance": 11740767.450958,
		"name": "Prairie Hill Road",
		"location": [-104.095128, 38.21453]
	}],
	"destinations": [{
		"hint": "",
		"distance": 9.215349,
		"name": "",
		"location": [-105.050583, 39.762548]
	}, {
		"hint": "",
		"distance": 11740767.450958,
		"name": "Prairie Hill Road",
		"location": [-104.095128, 38.21453]
	}],
	"durations": [
		[0, 17699.1],
		[17732.3, 0]
	],
	"distances": [
		[0, 245976.4],
		[245938.6, 0]
	]
}`

func TestEmptyPoints(t *testing.T) {
	ts := newTestServer(t, osrm.TableEndpoint)
	defer ts.s.Close()

	c := osrm.DefaultClient(ts.s.URL, true)

	_, _, err := c.Table([]measure.Point{})
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	_, _, err = c.Polyline([]measure.Point{})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
