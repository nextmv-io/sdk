package model

import "github.com/nextmv-io/sdk/hop/model/types"

// Domains creates a sequence of domains.
func Domains(domains ...types.Domain) types.Domains {
	connect()
	return domainsFunc(domains...)
}

// Repeat a domain n times.
func Repeat(n int, d types.Domain) types.Domains {
	connect()
	return repeatFunc(n, d)
}

var domainsFunc func(...types.Domain) types.Domains
var repeatFunc func(int, types.Domain) types.Domains
