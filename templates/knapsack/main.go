// package main holds the implementation of the knapsack template.
package main

import (
	"log"
	"math"
	"sort"
	"time"

	"github.com/nextmv-io/sdk/run"
	"github.com/nextmv-io/sdk/store"
)

func main() {
	_, err := run.Run(solver)
	if err != nil {
		log.Fatal(err)
	}
}

// A Knapsack holds the most valuable set of items possible while not exceeding
// its carrying capacity.
type input struct {
	Capacity int    `json:"capacity"` // weight the knapsack can hold
	Items    []item `json:"items"`    // potential items for packing
}

// item represents an item than can be put into the knapsack.
type item struct {
	ID     string `json:"id"`
	Value  int    `json:"value"`
	Weight int    `json:"weight"`
}

func solver(input input, opts store.Options) (store.Solver, error) {
	// The model works under the assumptions that the items are sorted by
	// efficiency = value / weight i.e. the most value per weight. Hence, we
	// sort them here.

	sort.SliceStable(input.Items, func(i, j int) bool {
		value := func(x item) float64 {
			return float64(x.Value) / float64(x.Weight)
		}
		return value(input.Items[i]) > value(input.Items[j])
	})

	// We create a new store that stores everything we need to construct valid
	// knapsack solutions.
	knapsack := store.New()

	// We store the value of the knapsack (i.e. the sum of all item values).
	value := store.NewVar(knapsack, 0)
	// We store the weight of the knapsack (i.e. the sum of all item weights).
	weight := store.NewVar(knapsack, 0)
	// itemIdx points to the most recent item we either put in our knapsack
	// or not.
	itemIdx := store.NewVar(knapsack, -1)
	// `trace`` stores our decisions along the way recording for each item
	// until `itemIdx` if it is part of the knapsack or not.
	trace := make([]store.Var[bool], len(input.Items))
	for i := 0; i < len(input.Items); i++ {
		trace[i] = store.NewVar(knapsack, false)
	}

	// Now we configure how we value, format and bound stores and, most
	// importantly, how we generate new stores.
	knapsack = knapsack.Generate(func(s store.Store) store.Generator {
		next := itemIdx.Get(s) + 1
		// i is a counter that ensures we only generate two new stores: one
		// with the item and one without.
		i := 0
		return store.Lazy(
			// Define if we want to branch off more stores from store s
			// We continue with a new branch for including and excluding
			// an item (i <= 2), if there are items left
			// (next < len(input.Items) and if there is capacity for the
			// item in the knapsack (weight.Get(s) <= input.Capacity).
			func() bool {
				return i <= 2 && next < len(input.Items) &&
					weight.Get(s) <= input.Capacity
			},
			// Define how to make a new store from store s
			// If we invoke the first time (i=0) we create a new store
			// out of s by changing it so the item is included in the
			// knapsack. The second time invoked (i=1) we change s such
			// that it is not included in the knapsack.
			func() store.Store {
				changes := make([]store.Change, 0)

				i++
				takeItem := i == 1
				if takeItem {
					newWeight := weight.Get(s) + input.Items[next].Weight
					changes = append(changes,
						weight.Set(newWeight),
						value.Set(value.Get(s)+input.Items[next].Value),
						trace[next].Set(takeItem),
						itemIdx.Set(next),
					)
				} else {
					changes = append(changes,
						trace[next].Set(takeItem),
						itemIdx.Set(next),
					)
				}
				return s.Apply(changes...)
			},
		)
	}).Validate(
		// The store is operationally valid if the capacity is not exceeded.
		func(s store.Store) bool {
			return weight.Get(s) <= input.Capacity
		}).Value(
		// Define the value of a store for optimization, the value is the sum
		// of the values of the individual items in the knapsack.
		value.Get,
	).Bound(
		// The last important ingredient is the `Bound` function.
		// Given a store we try to estimate what the best possible value could be
		// if we continue the search from here.
		// Since we sorted our items by efficiency we sum the values of all
		// remaining items until there is no more space left in the knapsack.
		// In case the last item does not fit fully,
		// we add it partially.
		// This forms a valid upper bound of what the final value could be.
		// The solver can use that information to more efficiently decide
		// what stores are being generated and how the search tree is traversed.
		// In general: the better the bounds, the better the search and the faster
		// we can prove optimality.
		func(s store.Store) store.Bounds {
			upperBound := float64(value.Get(s))
			currentWeight := float64(weight.Get(s))
			lastItemIdx := itemIdx.Get(s)
			for i := lastItemIdx + 1; i < len(input.Items); i++ {
				spaceLeft := float64(input.Capacity) - currentWeight
				weightNeeded := float64(input.Items[i].Weight)
				val := float64(input.Items[i].Value)
				if weightNeeded <= spaceLeft {
					upperBound += val
					currentWeight += weightNeeded
				} else {
					// in this case we just take the part of the item
					// after that our knapsack is full
					fraction := spaceLeft / weightNeeded
					upperBound += fraction * val
					break
				}
			}

			return store.Bounds{
				Lower: value.Get(s),
				Upper: int(math.Ceil(upperBound)),
			}
		}).Format(format(trace, value, weight, input))
	// A duration limit of 0 is treated as infinity. For cloud runs you need to
	// set an explicit duration limit which is why it is currently set to 10s
	// here in case no duration limit is set. For local runs there is no time
	// limitation. If you want to make cloud runs for longer than 5 minutes,
	// please contact: support@nextmv.io
	if opts.Limits.Duration == 0 {
		opts.Limits.Duration = 10 * time.Second
	}

	// We are optimizing to search for the largest value of the knapsack.
	return knapsack.Maximizer(opts), nil
}

// format returns a function to format the solution output.
func format(
	trace []store.Var[bool],
	value store.Var[int],
	weight store.Var[int],
	input input,
) func(s store.Store) any {
	return func(s store.Store) any {
		selectedItems := make([]item, 0, len(trace))
		for i, v := range trace {
			if v.Get(s) {
				selectedItems = append(selectedItems, input.Items[i])
			}
		}
		return map[string]any{
			"items":  selectedItems,
			"value":  value.Get(s),
			"weight": weight.Get(s),
		}
	}
}
