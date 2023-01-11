module github.com/nextmv-io/sdk

go 1.19

replace github.com/nextmv-io/sdk/route/osrm => ./route/osrm

require github.com/nextmv-io/sdk/route/osrm v0.20.9

require (
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/twpayne/go-polyline v1.1.1 // indirect
)
