/*
Package here provides a client for measuring distances and durations.

A HERE client requests distance and duration data using HERE Maps API. It
makes requests to construct a matrix measure.

The client can construct a distance matrix, a duration matrix, or both.

Each of these functions will use a synchronous request flow if the number
of points requested is below HERE's size limit for synchronous API calls -
otherwise, an asynchronous flow will be used. The functions all take a
context which can be used to cancel the request flow while it is in progress.

These measures implement measure.ByIndex.

These matrix-generating functions can also take one or more options
that allow you to configure the routes that will be included in the matrices.
*/
package here
