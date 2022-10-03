package model_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/nextmv-io/sdk/model"
)

type TestInt float64

func (t TestInt) ID() string {
	return strconv.Itoa(int(t))
}

func TestMultiMapOneSet(t *testing.T) {
	ints := []TestInt{1, 2, 3, 4, 5}

	// creates an element to be stored in the multimap based on the passed in
	// index
	create := func(index ...TestInt) float64 {
		returnValue := 0.0
		for i, v := range index {
			returnValue += math.Pow(2, float64(i)) * float64(v)
		}
		return returnValue
	}
	x := model.NewMultiMap(create, ints)

	// initially length is 0 - there are no elements in the multimap
	if x.Length() != 0 {
		t.Error("Expected length to be 0")
	}

	for _, v := range ints {
		// by executing "Get" the element as defined by the create function gets
		// stored in the multimap
		actual := x.Get(v)
		expected := float64(v)
		if expected != actual {
			t.Errorf("Expected %v, Actual: %v", expected, actual)
		}
	}

	// there are as many values in the multimap as there are values in 'ints'
	if x.Length() != len(ints) {
		t.Errorf("Expected length to be %d", len(ints))
	}

	// repeating should return the same values
	for _, v := range ints {
		actual := x.Get(v)
		expected := float64(v)
		if expected != actual {
			t.Errorf("Expected %v, Actual: %v", expected, actual)
		}
	}

	// but there should be no new values in the map
	if x.Length() != len(ints) {
		t.Errorf("Expected length to be %d", len(ints))
	}
}

func TestMultiMapTwoSets(t *testing.T) {
	ints := []TestInt{1, 2, 3, 4, 5}

	// creates an element to be stored in the multimap based on the passed in
	// index
	create := func(index ...TestInt) float64 {
		returnValue := 0.0
		for i, v := range index {
			returnValue += math.Pow(2, float64(i)) * float64(v)
		}
		return returnValue
	}
	x := model.NewMultiMap(create, ints, ints)

	// initially length is 0 - there are no elements in the multimap
	if x.Length() != 0 {
		t.Error("Expected length to be 0")
	}

	for _, v := range ints {
		// by executing "Get" the element as defined by the create function gets
		// stored in the multimap
		actual := x.Get(v, v)
		expected := float64(v) + 2*float64(v)
		if expected != actual {
			t.Errorf("Expected %v, Actual: %v", expected, actual)
		}
	}

	// there are as many values in the multimap as there are values in 'ints'
	if x.Length() != len(ints) {
		t.Errorf("Expected length to be %d", len(ints))
	}

	// repeating should return the same values
	for _, v := range ints {
		actual := x.Get(v, v)
		expected := float64(v) + 2*float64(v)
		if expected != actual {
			t.Errorf("Expected %v, Actual: %v", expected, actual)
		}
	}

	// but there should be no new values in the map
	if x.Length() != len(ints) {
		t.Errorf("Expected length to be %d", len(ints))
	}
}
