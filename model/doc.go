/*
Package model provides modeling components, such as domains and ranges.

Create individual integer domains.

    d1 := model.NewDomain() // empty domain
    d2 := model.NewDomain(model.NewRange(1, 10))
    d3 := model.NewDomain(model.NewRange(1, 10), model.MewRange(21, 30))
    d4 := model.Singleton(42)
    d5 := model.Multiple(2, 4, 6, 8)

Create sequences of domains.

    domains1 := model.NewDomains(d1, d2, d3, d4, d5)
    domains2 := model.Repeat(5, d1) // 5 empty domains
*/
package model
