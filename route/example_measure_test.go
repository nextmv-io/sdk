package route_test

import (
	"fmt"

	"github.com/nextmv-io/sdk/route"
)

func ExampleConstant() {
	constant := route.Constant(1.234)
	fmt.Println(constant.Cost(1, 0))
	fmt.Println(constant.Cost(0, 0))
	fmt.Println(constant.Cost(0, 1))
	// Output:
	// 1.234
	// 1.234
	// 1.234
}

func ExampleOverride() {
	defaultByIndex := route.Constant(1.234)
	overrideByIndex := route.Constant(4.321)
	condition := func(from, to int) bool {
		return from > 0
	}
	override := route.Override(
		defaultByIndex,
		overrideByIndex,
		condition,
	)
	fmt.Println(override.Cost(0, 0))
	fmt.Println(override.Cost(0, 1))
	fmt.Println(override.Cost(1, 0))
	// Output:
	// 1.234
	// 1.234
	// 4.321
}

func ExampleHaversineByPoint() {
	byPoint := route.HaversineByPoint()
	measure := byPoint.Cost(route.Point{1, 2}, route.Point{4, 5})
	fmt.Println(int(measure))
	// Output:
	// 471293
}

func ExampleConstantByPoint() {
	byPoint := route.ConstantByPoint(1.234)
	measure := byPoint.Cost(route.Point{1, 2}, route.Point{4, 5})
	fmt.Println(measure)
	// Output:
	// 1.234
}

func ExampleIndexed() {
	haversineByPoint := route.HaversineByPoint()
	points := []route.Point{
		{1, 2},
		{4, 5},
	}
	indexed := route.Indexed(haversineByPoint, points)
	fmt.Println(int(indexed.Cost(0, 1)))
	// Output:
	// 471293
}

func ExampleScale() {
	byIndex := route.Constant(1.234)
	scaled := route.Scale(byIndex, 2.0)
	fmt.Println(scaled.Cost(0, 1))
	// Output:
	// 2.468
}
