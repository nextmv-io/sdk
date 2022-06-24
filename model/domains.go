package model

// NewDomains creates a sequence of domains.
func NewDomains(domains ...Domain) Domains {
	connect()
	return newDomainsFunc(domains...)
}

// Repeat a domain n times.
func Repeat(n int, d Domain) Domains {
	connect()
	return repeatFunc(n, d)
}

var (
	newDomainsFunc func(...Domain) Domains
	repeatFunc     func(int, Domain) Domains
)
