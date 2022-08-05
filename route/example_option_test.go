package route_test

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nextmv-io/sdk/model"
	"github.com/nextmv-io/sdk/route"
	"github.com/nextmv-io/sdk/store"
)

// Create routes to visit seven landmarks in Kyoto using two vehicles. The
// vehicles have starting locations.
func ExampleStarts() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
		"v2",
	}

	// Define vehicle start locations.
	starts := []route.Position{
		{Lon: 135.737230, Lat: 35.043810}, // v1
		{Lon: 135.771716, Lat: 34.951317}, // v2
	}

	// Declare the router and its solver.
	router, err := route.NewRouter(stops, vehicles, route.Starts(starts))
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "unassigned": [],
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [
	//         {
	//           "id": "v1-start",
	//           "position": {
	//             "lon": 135.73723,
	//             "lat": 35.04381
	//           }
	//         },
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           }
	//         },
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           }
	//         }
	//       ],
	//       "route_duration": 664
	//     },
	//     {
	//       "id": "v2",
	//       "route": [
	//         {
	//           "id": "v2-start",
	//           "position": {
	//             "lon": 135.771716,
	//             "lat": 34.951317
	//           }
	//         },
	//         {
	//           "id": "Fushimi Inari Taisha",
	//           "position": {
	//             "lon": 135.772695,
	//             "lat": 34.967146
	//           }
	//         },
	//         {
	//           "id": "Kiyomizu-dera",
	//           "position": {
	//             "lon": 135.78506,
	//             "lat": 34.994857
	//           }
	//         },
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           }
	//         },
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           }
	//         },
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           }
	//         }
	//       ],
	//       "route_duration": 1085
	//     }
	//   ]
	// }
}

// Create routes to visit seven landmarks in Kyoto using two vehicles. The
// vehicles have starting and ending locations. Vehicle v1 starts at a point
// with no ending being set. Vehicle v2 starts and ends at the same geographical
// position. Endings could also be set as a standalone option.
func ExampleEnds() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
		"v2",
	}

	// Define vehicle start and end locations.
	starts := []route.Position{
		{Lon: 135.737230, Lat: 35.043810}, // v1
		{Lon: 135.758794, Lat: 34.986080}, // v2
	}
	ends := []route.Position{
		{},                                // v1
		{Lon: 135.758794, Lat: 34.986080}, // v2
	}

	// Declare the router and its solver.
	router, err := route.NewRouter(
		stops,
		vehicles,
		route.Starts(starts),
		route.Ends(ends),
	)
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "unassigned": [],
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [
	//         {
	//           "id": "v1-start",
	//           "position": {
	//             "lon": 135.73723,
	//             "lat": 35.04381
	//           }
	//         },
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           }
	//         },
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           }
	//         }
	//       ],
	//       "route_duration": 664
	//     },
	//     {
	//       "id": "v2",
	//       "route": [
	//         {
	//           "id": "v2-start",
	//           "position": {
	//             "lon": 135.758794,
	//             "lat": 34.98608
	//           }
	//         },
	//         {
	//           "id": "Fushimi Inari Taisha",
	//           "position": {
	//             "lon": 135.772695,
	//             "lat": 34.967146
	//           }
	//         },
	//         {
	//           "id": "Kiyomizu-dera",
	//           "position": {
	//             "lon": 135.78506,
	//             "lat": 34.994857
	//           }
	//         },
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           }
	//         },
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           }
	//         },
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           }
	//         },
	//         {
	//           "id": "v2-end",
	//           "position": {
	//             "lon": 135.758794,
	//             "lat": 34.98608
	//           }
	//         }
	//       ],
	//       "route_duration": 1484
	//     }
	//   ]
	// }
}

// Create routes to visit seven landmarks in Kyoto using two vehicles. The stops
// have quantities that must be fulfilled. The vehicles have starting locations
// and a maximum capacity that they can service.
func ExampleCapacity() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
		"v2",
	}

	// Define vehicle start locations.
	starts := []route.Position{
		{Lon: 135.737230, Lat: 35.043810}, // v1
		{Lon: 135.771716, Lat: 34.951317}, // v2
	}

	// Defines stop quantities and vehicle capacities.
	quantities := []int{
		-1, // "Fushimi Inari Taisha"
		-1, // "Kiyomizu-dera"
		-3, // "Nijō Castle"
		-1, // "Kyoto Imperial Palace"
		-1, // "Gionmachi"
		-3, // "Kinkaku-ji"
		-3, // "Arashiyama Bamboo Forest"
	}
	capacities := []int{
		9, // v1
		4, // v2
	}

	// Declare the router and its solver.
	router, err := route.NewRouter(
		stops,
		vehicles,
		route.Starts(starts),
		route.Capacity(quantities, capacities),
	)
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "unassigned": [],
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [
	//         {
	//           "id": "v1-start",
	//           "position": {
	//             "lon": 135.73723,
	//             "lat": 35.04381
	//           }
	//         },
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           }
	//         },
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           }
	//         },
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           }
	//         }
	//       ],
	//       "route_duration": 1116
	//     },
	//     {
	//       "id": "v2",
	//       "route": [
	//         {
	//           "id": "v2-start",
	//           "position": {
	//             "lon": 135.771716,
	//             "lat": 34.951317
	//           }
	//         },
	//         {
	//           "id": "Fushimi Inari Taisha",
	//           "position": {
	//             "lon": 135.772695,
	//             "lat": 34.967146
	//           }
	//         },
	//         {
	//           "id": "Kiyomizu-dera",
	//           "position": {
	//             "lon": 135.78506,
	//             "lat": 34.994857
	//           }
	//         },
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           }
	//         },
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           }
	//         }
	//       ],
	//       "route_duration": 908
	//     }
	//   ]
	// }
}

// Create routes to visit seven landmarks in Kyoto using two vehicles. In
// addition precedences for stops are defined. The vehicles have no starting
// locations and no maximum capacity that they can service.
func ExamplePrecedence() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
		"v2",
	}

	// Defines precedences for stops. In each couple the first ID precedes the
	// second ID.
	precedences := []route.Job{
		{PickUp: "Fushimi Inari Taisha", DropOff: "Kiyomizu-dera"},
		{PickUp: "Nijō Castle", DropOff: "Kiyomizu-dera"},
	}
	// Declare the router and its solver.
	router, err := route.NewRouter(
		stops,
		vehicles,
		route.Precedence(precedences),
	)
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "unassigned": [],
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           }
	//         },
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           }
	//         },
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           }
	//         },
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           }
	//         },
	//         {
	//           "id": "Fushimi Inari Taisha",
	//           "position": {
	//             "lon": 135.772695,
	//             "lat": 34.967146
	//           }
	//         },
	//         {
	//           "id": "Kiyomizu-dera",
	//           "position": {
	//             "lon": 135.78506,
	//             "lat": 34.994857
	//           }
	//         }
	//       ],
	//       "route_duration": 1517
	//     },
	//     {
	//       "id": "v2",
	//       "route": [
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           }
	//         }
	//       ],
	//       "route_duration": 0
	//     }
	//   ]
	// }
}

// Create routes to visit seven landmarks in Kyoto using two vehicles. The
// stops have unassigned penalties.
func ExampleUnassigned() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
		"v2",
	}

	// Define unassigned penalties.
	penalties := []int{
		0,      // "Fushimi Inari Taisha"
		0,      // "Kiyomizu-dera"
		200000, // "Nijō Castle"
		200000, // "Kyoto Imperial Palace"
		200000, // "Gionmachi"
		200000, // "Kinkaku-ji"
		200000, // "Arashiyama Bamboo Forest"
	}

	// Declare the router and its solver.
	router, err := route.NewRouter(
		stops,
		vehicles,
		route.Unassigned(penalties),
		route.Threads(1),
	)
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "unassigned": [
	//     {
	//       "id": "Fushimi Inari Taisha",
	//       "position": {
	//         "lon": 135.772695,
	//         "lat": 34.967146
	//       }
	//     },
	//     {
	//       "id": "Kiyomizu-dera",
	//       "position": {
	//         "lon": 135.78506,
	//         "lat": 34.994857
	//       }
	//     }
	//   ],
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           }
	//         },
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           }
	//         },
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           }
	//         },
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           }
	//         }
	//       ],
	//       "route_duration": 795
	//     },
	//     {
	//       "id": "v2",
	//       "route": [
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           }
	//         }
	//       ],
	//       "route_duration": 0
	//     }
	//   ]
	// }
}

// Create routes to visit seven landmarks in Kyoto using two vehicles with
// service times.
func ExampleServices() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
		"v2",
	}

	serviceTimes := []route.Service{
		{
			ID:       "Fushimi Inari Taisha",
			Duration: 900,
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Duration: 900,
		},
	}

	// Declare the router and its solver.
	router, err := route.NewRouter(
		stops,
		vehicles,
		route.Services(serviceTimes),
		route.Threads(1),
	)
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "unassigned": [],
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [
	//         {
	//           "id": "Fushimi Inari Taisha",
	//           "position": {
	//             "lon": 135.772695,
	//             "lat": 34.967146
	//           }
	//         },
	//         {
	//           "id": "Kiyomizu-dera",
	//           "position": {
	//             "lon": 135.78506,
	//             "lat": 34.994857
	//           }
	//         },
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           }
	//         },
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           }
	//         },
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           }
	//         },
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           }
	//         }
	//       ],
	//       "route_duration": 2143
	//     },
	//     {
	//       "id": "v2",
	//       "route": [
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           }
	//         }
	//       ],
	//       "route_duration": 900
	//     }
	//   ]
	// }
}

// Create routes to visit seven landmarks in Kyoto using two vehicles with
// shifts and service time.
func ExampleShifts() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
		"v2",
	}

	serviceTimes := []route.Service{
		{
			ID:       "Fushimi Inari Taisha",
			Duration: 900,
		},
	}

	// Define shifts for every vehicle
	shifts := []route.TimeWindow{
		{
			Start: time.Date(2020, 10, 17, 9, 0, 0, 0, time.UTC),
			End:   time.Date(2020, 10, 17, 17, 0, 0, 0, time.UTC),
		},
		{
			Start: time.Date(2020, 10, 17, 9, 0, 0, 0, time.UTC),
			End:   time.Date(2020, 10, 17, 17, 0, 0, 0, time.UTC),
		},
	}
	// Declare the router and its solver.
	router, err := route.NewRouter(
		stops,
		vehicles,
		route.Services(serviceTimes),
		route.Shifts(shifts),
		route.Threads(1),
	)
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "unassigned": [],
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [
	//         {
	//           "id": "Fushimi Inari Taisha",
	//           "position": {
	//             "lon": 135.772695,
	//             "lat": 34.967146
	//           },
	//           "estimated_arrival": "2020-10-17T09:00:00Z",
	//           "estimated_departure": "2020-10-17T09:15:00Z"
	//         },
	//         {
	//           "id": "Kiyomizu-dera",
	//           "position": {
	//             "lon": 135.78506,
	//             "lat": 34.994857
	//           },
	//           "estimated_arrival": "2020-10-17T09:20:28Z",
	//           "estimated_departure": "2020-10-17T09:20:28Z"
	//         },
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           },
	//           "estimated_arrival": "2020-10-17T09:22:28Z",
	//           "estimated_departure": "2020-10-17T09:22:28Z"
	//         },
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           },
	//           "estimated_arrival": "2020-10-17T09:27:12Z",
	//           "estimated_departure": "2020-10-17T09:27:12Z"
	//         },
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           },
	//           "estimated_arrival": "2020-10-17T09:30:10Z",
	//           "estimated_departure": "2020-10-17T09:30:10Z"
	//         },
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           },
	//           "estimated_arrival": "2020-10-17T09:35:43Z",
	//           "estimated_departure": "2020-10-17T09:35:43Z"
	//         }
	//       ],
	//       "route_duration": 2143
	//     },
	//     {
	//       "id": "v2",
	//       "route": [
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           },
	//           "estimated_arrival": "2020-10-17T09:00:00Z",
	//           "estimated_departure": "2020-10-17T09:00:00Z"
	//         }
	//       ],
	//       "route_duration": 0
	//     }
	//   ]
	// }
}

// Create routes to visit seven landmarks in Kyoto using two vehicles with
// shifts. The stops have time windows.
func ExampleWindows() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
		"v2",
	}

	serviceTimes := []route.Service{
		{
			ID:       "Gionmachi",
			Duration: 900,
		},
	}

	// Define time windows for every stop.
	windows := []route.Window{
		{
			TimeWindow: route.TimeWindow{
				Start: time.Date(2020, 10, 17, 7, 0, 0, 0, time.UTC),
				End:   time.Date(2020, 10, 17, 12, 0, 0, 0, time.UTC),
			},
			MaxWait: 900,
		},
		{},
		{},
		{},
		{},
		{},
		{},
	}

	// Define shifts for every vehicle
	shifts := []route.TimeWindow{
		{
			Start: time.Date(2020, 10, 17, 9, 0, 0, 0, time.UTC),
			End:   time.Date(2020, 10, 17, 17, 0, 0, 0, time.UTC),
		},
		{
			Start: time.Date(2020, 10, 17, 9, 0, 0, 0, time.UTC),
			End:   time.Date(2020, 10, 17, 17, 0, 0, 0, time.UTC),
		},
	}
	// Declare the router and its solver.
	router, err := route.NewRouter(
		stops,
		vehicles,
		route.Services(serviceTimes),
		route.Shifts(shifts),
		route.Windows(windows),
		route.Threads(1),
	)
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "unassigned": [],
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [
	//         {
	//           "id": "Fushimi Inari Taisha",
	//           "position": {
	//             "lon": 135.772695,
	//             "lat": 34.967146
	//           },
	//           "estimated_arrival": "2020-10-17T09:00:00Z",
	//           "estimated_departure": "2020-10-17T09:00:00Z"
	//         },
	//         {
	//           "id": "Kiyomizu-dera",
	//           "position": {
	//             "lon": 135.78506,
	//             "lat": 34.994857
	//           },
	//           "estimated_arrival": "2020-10-17T09:05:28Z",
	//           "estimated_departure": "2020-10-17T09:05:28Z"
	//         },
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           },
	//           "estimated_arrival": "2020-10-17T09:07:28Z",
	//           "estimated_departure": "2020-10-17T09:22:28Z"
	//         },
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           },
	//           "estimated_arrival": "2020-10-17T09:27:12Z",
	//           "estimated_departure": "2020-10-17T09:27:12Z"
	//         },
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           },
	//           "estimated_arrival": "2020-10-17T09:30:10Z",
	//           "estimated_departure": "2020-10-17T09:30:10Z"
	//         },
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           },
	//           "estimated_arrival": "2020-10-17T09:35:43Z",
	//           "estimated_departure": "2020-10-17T09:35:43Z"
	//         }
	//       ],
	//       "route_duration": 2143
	//     },
	//     {
	//       "id": "v2",
	//       "route": [
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           },
	//           "estimated_arrival": "2020-10-17T09:00:00Z",
	//           "estimated_departure": "2020-10-17T09:00:00Z"
	//         }
	//       ],
	//       "route_duration": 0
	//     }
	//   ]
	// }
}

// Create routes to visit seven landmarks in Kyoto using two vehicles. One
// vehicle has a backlog.
func ExampleBacklogs() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
		"v2",
	}
	// Define backlog for vehicle one.
	backlog := []route.Backlog{
		{
			VehicleID: "v1",
			Stops:     []string{"Kinkaku-ji", "Kyoto Imperial Palace"},
		},
	}

	// Declare the router and its solver.
	router, err := route.NewRouter(stops, vehicles, route.Backlogs(backlog))
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "unassigned": [],
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           }
	//         },
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           }
	//         },
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           }
	//         },
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           }
	//         },
	//         {
	//           "id": "Kiyomizu-dera",
	//           "position": {
	//             "lon": 135.78506,
	//             "lat": 34.994857
	//           }
	//         },
	//         {
	//           "id": "Fushimi Inari Taisha",
	//           "position": {
	//             "lon": 135.772695,
	//             "lat": 34.967146
	//           }
	//         }
	//       ],
	//       "route_duration": 1243
	//     },
	//     {
	//       "id": "v2",
	//       "route": [
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           }
	//         }
	//       ],
	//       "route_duration": 0
	//     }
	//   ]
	// }
}

// Create routes to visit seven landmarks in Kyoto using one vehicle. The route
// distance is minimized.
func ExampleMinimize() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
	}

	// Declare the router and its solver.
	router, err := route.NewRouter(stops, vehicles, route.Minimize())
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           }
	//         },
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           }
	//         },
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           }
	//         },
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           }
	//         },
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           }
	//         },
	//         {
	//           "id": "Kiyomizu-dera",
	//           "position": {
	//             "lon": 135.78506,
	//             "lat": 34.994857
	//           }
	//         },
	//         {
	//           "id": "Fushimi Inari Taisha",
	//           "position": {
	//             "lon": 135.772695,
	//             "lat": 34.967146
	//           }
	//         }
	//       ],
	//       "route_duration": 1818
	//     }
	//   ]
	// }
}

// Create routes to visit seven landmarks in Kyoto using one vehicle. The route
// distance is maximized.
func ExampleMaximize() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
	}

	// Declare the router and its solver.
	router, err := route.NewRouter(stops, vehicles, route.Maximize())
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           }
	//         },
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           }
	//         },
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           }
	//         },
	//         {
	//           "id": "Kiyomizu-dera",
	//           "position": {
	//             "lon": 135.78506,
	//             "lat": 34.994857
	//           }
	//         },
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           }
	//         },
	//         {
	//           "id": "Fushimi Inari Taisha",
	//           "position": {
	//             "lon": 135.772695,
	//             "lat": 34.967146
	//           }
	//         },
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           }
	//         }
	//       ],
	//       "route_duration": 4569
	//     }
	//   ]
	// }
}

// Create routes to visit seven landmarks in Kyoto using two vehicles. A
// service group and starting locations are configured.
func ExampleServiceGroups() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
		"v2",
	}

	serviceGroups := []route.ServiceGroup{
		{
			Group:    []string{"Gionmachi", "Kinkaku-ji"},
			Duration: 300,
		},
	}
	starts := []route.Position{
		{Lon: 135.672009, Lat: 35.017209},
		{Lon: 135.672009, Lat: 35.017209},
	}

	// Declare the router and its solver.
	router, err := route.NewRouter(
		stops,
		vehicles,
		route.ServiceGroups(serviceGroups),
		route.Starts(starts),
	)
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "unassigned": [],
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [
	//         {
	//           "id": "v1-start",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           }
	//         },
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           }
	//         },
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           }
	//         },
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           }
	//         },
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           }
	//         },
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           }
	//         },
	//         {
	//           "id": "Kiyomizu-dera",
	//           "position": {
	//             "lon": 135.78506,
	//             "lat": 34.994857
	//           }
	//         },
	//         {
	//           "id": "Fushimi Inari Taisha",
	//           "position": {
	//             "lon": 135.772695,
	//             "lat": 34.967146
	//           }
	//         }
	//       ],
	//       "route_duration": 2418
	//     },
	//     {
	//       "id": "v2",
	//       "route": [
	//         {
	//           "id": "v2-start",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           }
	//         }
	//       ],
	//       "route_duration": 0
	//     }
	//   ]
	// }
}

// Create routes to visit seven landmarks in Kyoto using two vehicles. A
// limit constraint is configured.
func ExampleLimits() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
		"v2",
	}

	// Define route limits.
	routeLimits := []route.Limit{
		{
			Measure: route.Constant(42.0),
			Value:   1000000,
		},
		{
			Measure: route.Constant(42.0),
			Value:   float64(model.MaxInt),
		},
	}
	ignoreTriangularity := true

	// Declare the router and its solver.
	router, err := route.NewRouter(
		stops,
		vehicles,
		route.Limits(routeLimits, ignoreTriangularity),
	)
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "unassigned": [],
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [
	//         {
	//           "id": "Fushimi Inari Taisha",
	//           "position": {
	//             "lon": 135.772695,
	//             "lat": 34.967146
	//           }
	//         },
	//         {
	//           "id": "Kiyomizu-dera",
	//           "position": {
	//             "lon": 135.78506,
	//             "lat": 34.994857
	//           }
	//         },
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           }
	//         },
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           }
	//         },
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           }
	//         },
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           }
	//         }
	//       ],
	//       "route_duration": 1243
	//     },
	//     {
	//       "id": "v2",
	//       "route": [
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           }
	//         }
	//       ],
	//       "route_duration": 0
	//     }
	//   ]
	// }
}

// Create routes to visit seven landmarks in Kyoto using two vehicles. A
// durations limit constraint is configured.
func ExampleLimitDurations() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
		"v2",
	}

	// Define route limits.
	routeLimits := []float64{1000.0, 1000.0}
	ignoreTriangularity := true

	// Declare the router and its solver.
	router, err := route.NewRouter(
		stops,
		vehicles,
		route.LimitDurations(routeLimits, ignoreTriangularity),
		route.Threads(1),
	)
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "unassigned": [],
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [
	//         {
	//           "id": "Fushimi Inari Taisha",
	//           "position": {
	//             "lon": 135.772695,
	//             "lat": 34.967146
	//           }
	//         },
	//         {
	//           "id": "Kiyomizu-dera",
	//           "position": {
	//             "lon": 135.78506,
	//             "lat": 34.994857
	//           }
	//         },
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           }
	//         },
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           }
	//         },
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           }
	//         }
	//       ],
	//       "route_duration": 909
	//     },
	//     {
	//       "id": "v2",
	//       "route": [
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           }
	//         },
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           }
	//         }
	//       ],
	//       "route_duration": 575
	//     }
	//   ]
	// }
}

// Create routes to visit seven landmarks in Kyoto using two vehicles. A
// distance limit constraint is configured.
func ExampleLimitDistances() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
		"v2",
	}

	// Define route limits.
	routeLimits := []float64{10000.0, 10000.0}
	ignoreTriangularity := true

	// Declare the router and its solver.
	router, err := route.NewRouter(
		stops,
		vehicles,
		route.LimitDistances(routeLimits, ignoreTriangularity),
		route.Threads(1),
	)
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "unassigned": [],
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [
	//         {
	//           "id": "Fushimi Inari Taisha",
	//           "position": {
	//             "lon": 135.772695,
	//             "lat": 34.967146
	//           }
	//         },
	//         {
	//           "id": "Kiyomizu-dera",
	//           "position": {
	//             "lon": 135.78506,
	//             "lat": 34.994857
	//           }
	//         },
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           }
	//         },
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           }
	//         },
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           }
	//         }
	//       ],
	//       "route_duration": 909
	//     },
	//     {
	//       "id": "v2",
	//       "route": [
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           }
	//         },
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           }
	//         }
	//       ],
	//       "route_duration": 575
	//     }
	//   ]
	// }
}

// Create routes to visit seven landmarks in Kyoto using two vehicles. Groups
// are configured.
func ExampleGrouper() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
		"v2",
	}

	// Define groups.
	groups := [][]string{
		{"Fushimi Inari Taisha", "Kiyomizu-dera", "Nijō Castle"},
		{"Gionmachi", "Kinkaku-ji", "Arashiyama Bamboo Forest"},
	}

	// Declare the router and its solver.
	router, err := route.NewRouter(stops, vehicles, route.Grouper(groups))
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "unassigned": [],
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           }
	//         }
	//       ],
	//       "route_duration": 0
	//     },
	//     {
	//       "id": "v2",
	//       "route": [
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           }
	//         },
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           }
	//         },
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           }
	//         },
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           }
	//         },
	//         {
	//           "id": "Kiyomizu-dera",
	//           "position": {
	//             "lon": 135.78506,
	//             "lat": 34.994857
	//           }
	//         },
	//         {
	//           "id": "Fushimi Inari Taisha",
	//           "position": {
	//             "lon": 135.772695,
	//             "lat": 34.967146
	//           }
	//         }
	//       ],
	//       "route_duration": 1639
	//     }
	//   ]
	// }
}

// Create routes to visit seven landmarks in Kyoto using two vehicles with
// shifts. Each vehicle has a value function measure.
func ExampleValueFunctionMeasures() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
		"v2",
	}

	// Define time windows for every stop.
	windows := []route.Window{
		{},
		{},
		{},
		{},
		{
			TimeWindow: route.TimeWindow{
				Start: time.Date(2020, 10, 17, 9, 0, 0, 0, time.UTC),
				End:   time.Date(2020, 10, 17, 9, 15, 0, 0, time.UTC),
			},
			MaxWait: -1,
		},
		{},
		{},
	}

	// Define shifts for every vehicle
	shifts := []route.TimeWindow{
		{
			Start: time.Date(2020, 10, 17, 9, 0, 0, 0, time.UTC),
			End:   time.Date(2020, 10, 17, 9, 30, 0, 0, time.UTC),
		},
		{
			Start: time.Date(2020, 10, 17, 9, 0, 0, 0, time.UTC),
			End:   time.Date(2020, 10, 17, 9, 30, 0, 0, time.UTC),
		},
	}

	count := len(stops)
	points := make([]route.Point, count+2*len(vehicles))
	for s, stop := range stops {
		point := route.Point{
			stop.Position.Lon,
			stop.Position.Lat,
		}

		points[s] = point
	}

	measures := make([]route.ByIndex, len(vehicles))

	// Haversine measure and override cost of going to/from an empty
	// point.
	m := route.Indexed(route.HaversineByPoint(), points)
	m = route.Override(
		m,
		route.Constant(0),
		func(from, to int) bool {
			return points[from] == nil || points[to] == nil
		},
	)

	for v := range vehicles {
		// v1 and v2 have a speed of 7.0 m/s
		measures[v] = route.Scale(m, 1/7.0)
	}

	// Declare the router and its solver.
	router, err := route.NewRouter(
		stops,
		vehicles,
		route.ValueFunctionMeasures(measures),
		route.Shifts(shifts),
		route.Windows(windows),
		route.Threads(1),
	)
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "unassigned": [],
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [
	//         {
	//           "id": "Fushimi Inari Taisha",
	//           "position": {
	//             "lon": 135.772695,
	//             "lat": 34.967146
	//           },
	//           "estimated_arrival": "2020-10-17T09:00:00Z",
	//           "estimated_departure": "2020-10-17T09:00:00Z"
	//         },
	//         {
	//           "id": "Kiyomizu-dera",
	//           "position": {
	//             "lon": 135.78506,
	//             "lat": 34.994857
	//           },
	//           "estimated_arrival": "2020-10-17T09:05:28Z",
	//           "estimated_departure": "2020-10-17T09:05:28Z"
	//         },
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           },
	//           "estimated_arrival": "2020-10-17T09:07:28Z",
	//           "estimated_departure": "2020-10-17T09:07:28Z"
	//         },
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           },
	//           "estimated_arrival": "2020-10-17T09:12:12Z",
	//           "estimated_departure": "2020-10-17T09:12:12Z"
	//         },
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           },
	//           "estimated_arrival": "2020-10-17T09:15:10Z",
	//           "estimated_departure": "2020-10-17T09:15:10Z"
	//         },
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           },
	//           "estimated_arrival": "2020-10-17T09:20:43Z",
	//           "estimated_departure": "2020-10-17T09:20:43Z"
	//         }
	//       ],
	//       "route_duration": 1243
	//     },
	//     {
	//       "id": "v2",
	//       "route": [
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           },
	//           "estimated_arrival": "2020-10-17T09:00:00Z",
	//           "estimated_departure": "2020-10-17T09:00:00Z"
	//         }
	//       ],
	//       "route_duration": 0
	//     }
	//   ]
	// }
}

// Create routes to visit seven landmarks in Kyoto using two vehicles with
// shifts. Each vehicle has a travel time measure.
func ExampleTravelTimeMeasures() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
		"v2",
	}

	// Define time windows for every stop.
	windows := []route.Window{
		{},
		{},
		{},
		{},
		{
			TimeWindow: route.TimeWindow{
				Start: time.Date(2020, 10, 17, 9, 0, 0, 0, time.UTC),
				End:   time.Date(2020, 10, 17, 9, 15, 0, 0, time.UTC),
			},
			MaxWait: -1,
		},
		{},
		{},
	}

	// Define shifts for every vehicle
	shifts := []route.TimeWindow{
		{
			Start: time.Date(2020, 10, 17, 9, 0, 0, 0, time.UTC),
			End:   time.Date(2020, 10, 17, 9, 30, 0, 0, time.UTC),
		},
		{
			Start: time.Date(2020, 10, 17, 9, 0, 0, 0, time.UTC),
			End:   time.Date(2020, 10, 17, 9, 30, 0, 0, time.UTC),
		},
	}

	count := len(stops)
	points := make([]route.Point, count+2*len(vehicles))
	for s, stop := range stops {
		point := route.Point{
			stop.Position.Lon,
			stop.Position.Lat,
		}

		points[s] = point
	}

	measures := make([]route.ByIndex, len(vehicles))

	// Haversine measure and override cost of going to/from an empty
	// point.
	m := route.Indexed(route.HaversineByPoint(), points)
	m = route.Override(
		m,
		route.Constant(0),
		func(from, to int) bool {
			return points[from] == nil || points[to] == nil
		},
	)

	for v := range vehicles {
		// v1 and v2 have a speed of 7.0 m/s
		measures[v] = route.Scale(m, 1/7.0)
	}

	// Declare the router and its solver.
	router, err := route.NewRouter(
		stops,
		vehicles,
		route.TravelTimeMeasures(measures),
		route.Shifts(shifts),
		route.Windows(windows),
		route.Threads(1),
	)
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "unassigned": [],
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [
	//         {
	//           "id": "Fushimi Inari Taisha",
	//           "position": {
	//             "lon": 135.772695,
	//             "lat": 34.967146
	//           },
	//           "estimated_arrival": "2020-10-17T09:00:00Z",
	//           "estimated_departure": "2020-10-17T09:00:00Z"
	//         },
	//         {
	//           "id": "Kiyomizu-dera",
	//           "position": {
	//             "lon": 135.78506,
	//             "lat": 34.994857
	//           },
	//           "estimated_arrival": "2020-10-17T09:07:49Z",
	//           "estimated_departure": "2020-10-17T09:07:49Z"
	//         },
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           },
	//           "estimated_arrival": "2020-10-17T09:10:41Z",
	//           "estimated_departure": "2020-10-17T09:10:41Z"
	//         },
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           },
	//           "estimated_arrival": "2020-10-17T09:17:27Z",
	//           "estimated_departure": "2020-10-17T09:17:27Z"
	//         },
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           },
	//           "estimated_arrival": "2020-10-17T09:21:41Z",
	//           "estimated_departure": "2020-10-17T09:21:41Z"
	//         },
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           },
	//           "estimated_arrival": "2020-10-17T09:29:37Z",
	//           "estimated_departure": "2020-10-17T09:29:37Z"
	//         }
	//       ],
	//       "route_duration": 1777
	//     },
	//     {
	//       "id": "v2",
	//       "route": [
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           },
	//           "estimated_arrival": "2020-10-17T09:00:00Z",
	//           "estimated_departure": "2020-10-17T09:00:00Z"
	//         }
	//       ],
	//       "route_duration": 0
	//     }
	//   ]
	// }
}

// Create routes to visit seven landmarks in Kyoto using two vehicles. An
// attribute compatibility filter is configured.
func ExampleAttribute() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
		"v2",
	}

	// Define compatibility attributes.
	vehicleAttributes := []route.Attributes{
		{
			ID:         "v1",
			Attributes: []string{"Cooling System"},
		},
		{
			ID:         "v2",
			Attributes: []string{"Large"},
		},
	}
	stopAttributes := []route.Attributes{
		{
			ID:         "Fushimi Inari Taisha",
			Attributes: []string{"Cooling System"},
		},
		{
			ID:         "Arashiyama Bamboo Forest",
			Attributes: []string{"Large"},
		},
		{
			ID:         "Kinkaku-ji",
			Attributes: []string{"Large"},
		},
	}

	// Declare the router and its solver.
	router, err := route.NewRouter(
		stops,
		vehicles,
		route.Attribute(vehicleAttributes, stopAttributes),
	)
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "unassigned": [],
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [
	//         {
	//           "id": "Fushimi Inari Taisha",
	//           "position": {
	//             "lon": 135.772695,
	//             "lat": 34.967146
	//           }
	//         },
	//         {
	//           "id": "Kiyomizu-dera",
	//           "position": {
	//             "lon": 135.78506,
	//             "lat": 34.994857
	//           }
	//         },
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           }
	//         },
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           }
	//         },
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           }
	//         }
	//       ],
	//       "route_duration": 909
	//     },
	//     {
	//       "id": "v2",
	//       "route": [
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           }
	//         },
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           }
	//         }
	//       ],
	//       "route_duration": 575
	//     }
	//   ]
	// }
}

// Create routes to visit seven landmarks in Kyoto using two vehicles. Use a
// single thread.
func ExampleThreads() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
		"v2",
	}

	// Declare the router and its solver.
	router, err := route.NewRouter(stops, vehicles, route.Threads(1))
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "unassigned": [],
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [
	//         {
	//           "id": "Fushimi Inari Taisha",
	//           "position": {
	//             "lon": 135.772695,
	//             "lat": 34.967146
	//           }
	//         },
	//         {
	//           "id": "Kiyomizu-dera",
	//           "position": {
	//             "lon": 135.78506,
	//             "lat": 34.994857
	//           }
	//         },
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           }
	//         },
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           }
	//         },
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           }
	//         },
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           }
	//         }
	//       ],
	//       "route_duration": 1243
	//     },
	//     {
	//       "id": "v2",
	//       "route": [
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           }
	//         }
	//       ],
	//       "route_duration": 0
	//     }
	//   ]
	// }
}

// Create routes to visit seven landmarks in Kyoto using two vehicles. Alternate
// stops are configured.
func ExampleAlternates() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
		"v2",
	}

	alt := []route.Alternate{
		{
			VehicleID: "v1",
			Stops:     []string{"Kiyomizu-dera", "Gionmachi"},
		},
	}

	// Declare the router and its solver.
	router, err := route.NewRouter(stops, vehicles, route.Alternates(alt))
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "unassigned": [],
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           }
	//         },
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           }
	//         },
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           }
	//         },
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           }
	//         },
	//         {
	//           "id": "Fushimi Inari Taisha",
	//           "position": {
	//             "lon": 135.772695,
	//             "lat": 34.967146
	//           }
	//         }
	//       ],
	//       "route_duration": 1189
	//     },
	//     {
	//       "id": "v2",
	//       "route": [
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           }
	//         }
	//       ],
	//       "route_duration": 0
	//     }
	//   ]
	// }
}

// Create routes to visit seven landmarks in Kyoto using two vehicles.
// Velocities for the travel time measure are configured.
func ExampleVelocities() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
		"v2",
	}

	velocities := []float64{5, 7}

	// Declare the router and its solver.
	router, err := route.NewRouter(stops, vehicles, route.Velocities(velocities))
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "unassigned": [],
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [
	//         {
	//           "id": "Fushimi Inari Taisha",
	//           "position": {
	//             "lon": 135.772695,
	//             "lat": 34.967146
	//           }
	//         },
	//         {
	//           "id": "Kiyomizu-dera",
	//           "position": {
	//             "lon": 135.78506,
	//             "lat": 34.994857
	//           }
	//         },
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           }
	//         },
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           }
	//         },
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           }
	//         },
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           }
	//         }
	//       ],
	//       "route_duration": 2485
	//     },
	//     {
	//       "id": "v2",
	//       "route": [
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           }
	//         }
	//       ],
	//       "route_duration": 0
	//     }
	//   ]
	// }
}

// Create routes to visit seven landmarks in Kyoto using two vehicles.
// Initialization costs are configured.
func ExampleInitializationCosts() {
	// Define stops and vehicles.
	stops := []route.Stop{
		{
			ID:       "Fushimi Inari Taisha",
			Position: route.Position{Lon: 135.772695, Lat: 34.967146},
		},
		{
			ID:       "Kiyomizu-dera",
			Position: route.Position{Lon: 135.785060, Lat: 34.994857},
		},
		{
			ID:       "Nijō Castle",
			Position: route.Position{Lon: 135.748134, Lat: 35.014239},
		},
		{
			ID:       "Kyoto Imperial Palace",
			Position: route.Position{Lon: 135.762057, Lat: 35.025431},
		},
		{
			ID:       "Gionmachi",
			Position: route.Position{Lon: 135.775682, Lat: 35.002457},
		},
		{
			ID:       "Kinkaku-ji",
			Position: route.Position{Lon: 135.728898, Lat: 35.039705},
		},
		{
			ID:       "Arashiyama Bamboo Forest",
			Position: route.Position{Lon: 135.672009, Lat: 35.017209},
		},
	}
	vehicles := []string{
		"v1",
		"v2",
	}
	initializationCosts := []float64{100000, 0}

	// Declare the router and its solver.
	router, err := route.NewRouter(
		stops,
		vehicles,
		route.Threads(1),
		route.InitializationCosts(initializationCosts),
	)
	if err != nil {
		panic(err)
	}
	solver, err := router.Solver(store.DefaultOptions())
	if err != nil {
		panic(err)
	}

	// Get the last solution of the problem and print it.
	last := solver.Last(context.Background())
	b, err := json.MarshalIndent(last.Store, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output:
	// {
	//   "unassigned": [],
	//   "vehicles": [
	//     {
	//       "id": "v1",
	//       "route": [],
	//       "route_duration": 0
	//     },
	//     {
	//       "id": "v2",
	//       "route": [
	//         {
	//           "id": "Fushimi Inari Taisha",
	//           "position": {
	//             "lon": 135.772695,
	//             "lat": 34.967146
	//           }
	//         },
	//         {
	//           "id": "Kiyomizu-dera",
	//           "position": {
	//             "lon": 135.78506,
	//             "lat": 34.994857
	//           }
	//         },
	//         {
	//           "id": "Gionmachi",
	//           "position": {
	//             "lon": 135.775682,
	//             "lat": 35.002457
	//           }
	//         },
	//         {
	//           "id": "Kyoto Imperial Palace",
	//           "position": {
	//             "lon": 135.762057,
	//             "lat": 35.025431
	//           }
	//         },
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           }
	//         },
	//         {
	//           "id": "Kinkaku-ji",
	//           "position": {
	//             "lon": 135.728898,
	//             "lat": 35.039705
	//           }
	//         },
	//         {
	//           "id": "Arashiyama Bamboo Forest",
	//           "position": {
	//             "lon": 135.672009,
	//             "lat": 35.017209
	//           }
	//         }
	//       ],
	//       "route_duration": 1818
	//     }
	//   ]
	// }
}
