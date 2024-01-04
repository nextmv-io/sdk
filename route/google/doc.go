/*
Package google provides functions for measuring distances and durations using
the Google Distance Matrix API and polylines from Google Maps Distance API. A
Google Maps client and request are required. The client uses an API key for
authentication.

  - Matrix API: At a minimum, the request requires the origins and destinations
    to estimate.

  - Distance API: At minimum, the request requires the origin and destination.
    But it is recommended to pass in waypoints encoded as a polyline with "enc:"
    as a prefix to get a more precise polyline for each leg of the route.

Deprecated: This package is deprecated and will be removed in the next major release.
It is used with the router engine which was replaced by
[github.com/nextmv-io/sdk/measure/google].
*/
package google
