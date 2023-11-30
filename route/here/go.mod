module github.com/nextmv-io/sdk/route/here

go 1.19

replace github.com/nextmv-io/sdk => ../../.

replace github.com/nextmv-io/sdk/measure/here => ../../measure/here

require (
	github.com/google/go-cmp v0.5.6
	github.com/nextmv-io/sdk v1.0.3
)

require github.com/nextmv-io/sdk/measure/here v1.0.2

require github.com/nextmv-io/sdk v1.0.2 // indirect
