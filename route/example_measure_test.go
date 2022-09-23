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
	measure := byPoint.Cost(
		route.Point{135.772695, 34.967146},
		route.Point{135.78506, 34.994857},
	)
	fmt.Println(int(measure))
	// Output:
	// 3280
}

func ExampleConstantByPoint() {
	byPoint := route.ConstantByPoint(1.234)
	measure := byPoint.Cost(
		route.Point{135.772695, 34.967146},
		route.Point{135.78506, 34.994857},
	)
	fmt.Println(measure)
	// Output:
	// 1.234
}

func ExampleIndexed() {
	haversineByPoint := route.HaversineByPoint()
	points := []route.Point{
		{135.772695, 34.967146},
		{135.78506, 34.994857},
	}
	indexed := route.Indexed(haversineByPoint, points)
	fmt.Println(int(indexed.Cost(0, 1)))
	// Output:
	// 3280
}

func ExampleScale() {
	byIndex := route.Constant(1.234)
	scaled := route.Scale(byIndex, 2.0)
	fmt.Println(scaled.Cost(0, 1))
	// Output:
	// 2.468
}

func ExampleBin() {
	constant1 := route.Constant(1.234)
	constant2 := route.Constant(4.321)

	bin := route.Bin(
		[]route.ByIndex{constant1, constant2},
		func(from, to int) int {
			if from == 0 && to == 1 {
				return 0
			}
			return 1
		},
	)
	fmt.Println(bin.Cost(0, 1))
	fmt.Println(bin.Cost(1, 0))
	// Output:
	// 1.234
	// 4.321
}
