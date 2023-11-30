module github.com/nextmv-io/sdk/route/osrm

go 1.19

replace github.com/nextmv-io/sdk => ../../.

replace github.com/nextmv-io/sdk/measure/osrm => ../../measure/osrm

require (
	github.com/nextmv-io/sdk/measure/osrm v1.0.2
)

require (
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/twpayne/go-polyline v1.1.1 // indirect
	go.uber.org/mock v0.3.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/nextmv-io/sdk v1.0.3
	github.com/twpayne/go-polyline v1.1.1
)
