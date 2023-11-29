module github.com/nextmv-io/sdk/measure/osrm

go 1.19

replace github.com/nextmv-io/sdk => ../../.

require (
	github.com/hashicorp/golang-lru v0.5.4
	github.com/nextmv-io/sdk v1.0.2
	github.com/nextmv-io/sdk/route/osrm v0.21.1
	github.com/twpayne/go-polyline v1.1.1
	go.uber.org/mock v0.3.0
)
