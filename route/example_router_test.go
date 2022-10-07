package route_test

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nextmv-io/sdk/route"
	"github.com/nextmv-io/sdk/store"
)

// Create routes to visit seven landmarks in Kyoto using two vehicles.

func Example() {
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

// Create routes to visit seven landmarks in Kyoto using two vehicles. The
// vehicles have starting locations. Endings are configured via the Options
// function in a separate step.
func ExampleRouter_options() {
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
		{Lon: 135.758794, Lat: 34.986080}, // v2
	}

	// Declare the router. Ends are omitted from the initial set of options.
	router, err := route.NewRouter(stops, vehicles, route.Starts(starts))
	if err != nil {
		panic(err)
	}

	// Define ending locations and configure them as a separate step.
	ends := []route.Position{
		{},                                // v1
		{Lon: 135.758794, Lat: 34.986080}, // v2
	}
	err = router.Options(route.Ends(ends))
	if err != nil {
		panic(err)
	}

	// Declare the solver.
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

// Use the Plan func to get direct access to the underlying Store and variables.
func ExampleRouter_plan() {
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

	// Get the last solution of the problem.
	last := solver.Last(context.Background())
	plan := router.Plan().Get(last.Store)

	// Extract unassigned stops from the plan.
	unassigned := plan.Unassigned
	b, err := json.MarshalIndent(unassigned, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println("--- unassigned")
	fmt.Println(string(b))

	// Extract the routes from the plan.
	plannedVehicles := plan.Vehicles
	b, err = json.MarshalIndent(plannedVehicles, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println("--- vehicles")
	fmt.Println(string(b))
	// Output:
	// --- unassigned
	// [
	//   {
	//     "id": "Fushimi Inari Taisha",
	//     "position": {
	//       "lon": 135.772695,
	//       "lat": 34.967146
	//     }
	//   },
	//   {
	//     "id": "Kiyomizu-dera",
	//     "position": {
	//       "lon": 135.78506,
	//       "lat": 34.994857
	//     }
	//   }
	// ]
	// --- vehicles
	// [
	//   {
	//     "id": "v1",
	//     "route": [
	//       {
	//         "id": "Gionmachi",
	//         "position": {
	//           "lon": 135.775682,
	//           "lat": 35.002457
	//         }
	//       },
	//       {
	//         "id": "Kyoto Imperial Palace",
	//         "position": {
	//           "lon": 135.762057,
	//           "lat": 35.025431
	//         }
	//       },
	//       {
	//         "id": "Nijō Castle",
	//         "position": {
	//           "lon": 135.748134,
	//           "lat": 35.014239
	//         }
	//       },
	//       {
	//         "id": "Kinkaku-ji",
	//         "position": {
	//           "lon": 135.728898,
	//           "lat": 35.039705
	//         }
	//       }
	//     ],
	//     "route_duration": 795
	//   },
	//   {
	//     "id": "v2",
	//     "route": [
	//       {
	//         "id": "Arashiyama Bamboo Forest",
	//         "position": {
	//           "lon": 135.672009,
	//           "lat": 35.017209
	//         }
	//       }
	//     ],
	//     "route_duration": 0
	//   }
	// ]
}

func ExampleNewRouter() {
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
	}
	vehicles := []string{
		"v1",
		"v2",
	}
	quantities := []int{-1, -1, -1}
	capacities := []int{2, 2}

	// Declare the router and its solver.
	router, err := route.NewRouter(
		stops,
		vehicles,
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
	//       "route_duration": 328
	//     },
	//     {
	//       "id": "v2",
	//       "route": [
	//         {
	//           "id": "Nijō Castle",
	//           "position": {
	//             "lon": 135.748134,
	//             "lat": 35.014239
	//           }
	//         }
	//       ],
	//       "route_duration": 0
	//     }
	//   ]
	// }
}
