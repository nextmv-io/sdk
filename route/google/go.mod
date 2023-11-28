module github.com/nextmv-io/sdk/route/google

go 1.19

replace github.com/nextmv-io/sdk => ../../.

replace github.com/nextmv-io/sdk/measure/google => ../../measure/google

require (
	github.com/nextmv-io/sdk v1.0.2
	github.com/nextmv-io/sdk/measure/google v1.0.2
	googlemaps.github.io/maps v1.4.0
)

require (
	github.com/google/uuid v1.3.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/time v0.3.0 // indirect
)
