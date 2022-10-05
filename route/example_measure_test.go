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

func ExampleDebugOverride() {
	defaultByIndex := route.Constant(1.234)
	overrideByIndex := route.Constant(4.321)
	condition := func(from, to int) bool {
		return from > 0
	}
	override := route.DebugOverride(
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

func ExampleEuclideanByPoint() {
	byPoint := route.EuclideanByPoint()
	cost := byPoint.Cost(
		route.Point{135.772695, 34.967146},
		route.Point{135.78506, 34.994857},
	)
	fmt.Println(int(cost * 1000))
	// Output:
	// 30
}

func ExamplePower() {
	byPoint := route.Power(route.Constant(2), 2)
	cost := byPoint.Cost(0, 1)
	fmt.Println(int(cost))
	// Output:
	// 4
}

func ExampleSum() {
	byPoint := route.Sum(route.Constant(2), route.Constant(1))
	cost := byPoint.Cost(0, 1)
	fmt.Println(int(cost))
	// Output:
	// 3
}

func ExampleTaxicabByPoint() {
	byPoint := route.TaxicabByPoint()
	cost := byPoint.Cost(
		route.Point{135.772695, 34.967146},
		route.Point{135.78506, 34.994857},
	)
	fmt.Println(int(cost * 1000))
	// Output:
	// 40
}

func ExampleTruncate() {
	truncatedByUpper := route.Truncate(route.Constant(10), 1, 8)
	cost := truncatedByUpper.Cost(0, 1)
	fmt.Println(int(cost))

	truncatedByLower := route.Truncate(route.Constant(0), 1, 8)
	cost = truncatedByLower.Cost(0, 1)
	fmt.Println(int(cost))
	// Output:
	// 8
	// 1
}

func ExampleSparse() {
	costMap := make(map[int]map[int]float64)
	costMap[0] = map[int]float64{}
	costMap[0][1] = 10
	byPoint := route.Sparse(route.Constant(1), costMap)
	cost := byPoint.Cost(0, 1)
	fmt.Println(int(cost))
	cost = byPoint.Cost(1, 2)
	fmt.Println(int(cost))
	// Output:
	// 10
	// 1
}

func ExampleMatrix() {
	byPoint := route.Matrix([][]float64{
		{
			0,
			1,
		},
	})
	cost := byPoint.Cost(0, 0)
	fmt.Println(int(cost))
	cost = byPoint.Cost(0, 1)
	fmt.Println(int(cost))
	// Output:
	// 0
	// 1
}
