package common

import (
	"strings"
	"testing"
)

type item struct {
	name  string
	count int
	value float64
}

func (i item) Name() string {
	return i.name
}

func (i item) Count() int {
	return i.count
}

func (i item) Value() float64 {
	return i.value
}

func newItems() []item {
	return []item{
		{"a", 1, 1.0},
		{"b", 1, 2.0},
		{"c", 2, 3.0},
		{"d", 2, 3.0},
		{"e", 2, 4.0},
	}
}

func verifyOrder(original, items []item, order []int) bool {
	for i, j := range order {
		if items[i] != original[j] {
			return false
		}
	}
	return true
}

func verifyValues[T comparable](items []item, f func(item) T, v []T) bool {
	for i, item := range items {
		if f(item) != v[i] {
			return false
		}
	}
	return true
}
func TestSort(t *testing.T) {
	items := newItems()

	itemsByCount := Sort(
		items,
		CompareBy(item.Count),
	)
	if !verifyValues(itemsByCount, item.Count, []int{1, 1, 2, 2, 2}) {
		t.Error("itemsByCount")
	}

	itemsByValue := Sort(
		items,
		CompareBy(item.Value),
	)
	if !verifyValues(itemsByValue, item.Value, []float64{1.0, 2.0, 3.0, 3.0, 4.0}) {
		t.Error("itemsByValue")
	}

	itemsByName := Sort(
		items,
		func(a, b item) int {
			return strings.Compare(a.name, b.name)
		},
	)
	if !verifyOrder(items, itemsByName, []int{0, 1, 2, 3, 4}) {
		t.Error("itemsByName")
	}

	itemsByValueThenOppositeName := Sort(
		items,
		CompareBy(item.Value),
		CompareOpposite(func(a, b item) int {
			return strings.Compare(a.name, b.name)
		}),
	)
	if !verifyOrder(items, itemsByValueThenOppositeName, []int{0, 1, 3, 2, 4}) {
		t.Error("itemsByValueThenOppositeName")
	}

	emptyItems := []item{}

	sorted := Sort(
		emptyItems,
		CompareBy(item.Value),
	)

	emptyItems = append(emptyItems, item{"a", 1, 1.0})
	_ = emptyItems

	if len(sorted) != 0 {
		t.Error("emptyItems")
	}

	singleItem := []item{{"a", 1, 1.0}}
	_ = singleItem

	sorted = Sort(
		singleItem,
		CompareBy(item.Value),
	)

	singleItem = append(singleItem, item{"a", 1, 1.0})

	if len(sorted) != 1 {
		t.Error("singleItem")
	}
}
