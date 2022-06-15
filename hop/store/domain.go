package store

// // Domain creates a domain in a store.
// func Domain(s Store, ranges ...model.Range) Domain {
// 	return domainProxy{domain: Var(s, model.NewDomain(ranges...))}
// }

// // Singleton creates a domain with one value.
// func Singleton(s Store, value int) Domain {
// 	return domainProxy{domain: Var(s, model.Singleton(value))}
// }

// // Multiple creates a domain with multiple values.
// func Multiple(s Store, values ...int) Domain {
// 	return domainProxy{domain: Var(s, model.Multiple(values...))}
// }

// type domainProxy struct {
// 	domain types.Var[model.Domain]
// }

// // Implements Domain

// func (d domainProxy) Add(values ...int) Change {
// 	return func(s Store) { d.domain.Set(d.domain.Get(s).Add(values...)) }
// }

// func (d domainProxy) AtLeast(value int) Change {
// 	return func(s Store) { d.domain.Set(d.domain.Get(s).AtLeast(value)) }
// }

// func (d domainProxy) AtMost(value int) Change {
// 	return func(s Store) { d.domain.Set(d.domain.Get(s).AtMost(value)) }
// }

// func (d domainProxy) Contains(s Store, value int) bool {
// 	return d.domain.Get(s).Contains(value)
// }

// func (d domainProxy) Domain(s Store) model.Domain {
// 	return d.domain.Get(s)
// }

// func (d domainProxy) Empty(s Store) bool {
// 	return d.domain.Get(s).Empty()
// }

// func (d domainProxy) Len(s Store) int {
// 	return d.domain.Get(s).Len()
// }

// func (d domainProxy) Max(s Store) (int, bool) {
// 	return d.domain.Get(s).Max()
// }

// func (d domainProxy) Min(s Store) (int, bool) {
// 	return d.domain.Get(s).Min()
// }

// func (d domainProxy) Remove(values ...int) Change {
// 	return func(s Store) { d.domain.Set(d.domain.Get(s).Remove(values...)) }
// }

// func (d domainProxy) Slice(s Store) []int {
// 	return d.domain.Get(s).Slice()
// }

// func (d domainProxy) Value(s Store) (int, bool) {
// 	return d.domain.Get(s).Value()
// }
