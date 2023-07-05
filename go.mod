module github.com/nextmv-io/sdk

go 1.19

replace github.com/nextmv-io/sdk/route/osrm => ./route/osrm

require (
	github.com/danielgtaylor/huma v1.14.1
	github.com/nextmv-io/sdk/route/osrm v0.21.1
	github.com/xeipuuv/gojsonschema v1.2.0
)

require (
	github.com/google/uuid v1.3.0
	github.com/gorilla/schema v1.2.0
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/itzg/go-flagsfiller v1.9.1
	github.com/twpayne/go-polyline v1.1.1 // indirect
)

require (
	github.com/iancoleman/strcase v0.2.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
)
