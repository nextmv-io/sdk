package kmeans_test

import (
	"fmt"

	"github.com/nextmv-io/sdk/cluster/kmeans"
	"github.com/nextmv-io/sdk/measure"
)

func ExampleClusterModel() {
	points := []measure.Point{
		{2.5, 2.5},
		{7.5, 7.5},
		{5.0, 7.5},
	}
	// Create a model.
	model, err := kmeans.NewModel(points, 2)
	if err != nil {
		panic(err)
	}

	cm1 := model.ClusterModels()[0]
	cm2 := model.ClusterModels()[1]

	// Set maximum points in first cluster to 2.
	cm1.SetMaximumPoints(2)
	// Exclude the first point from the first cluster.
	err = cm1.SetExcludedPointIndices([]int{0})

	if err != nil {
		panic(err)
	}

	// The values of the points in the second cluster
	// must sum to 10 or less. The values of the points
	// used are 6,7 and 8 in order of the points in the model.
	msv, err := cm1.SetMaximumSumValue(
		10.0,
		[]int{6, 7, 4},
	)
	if err != nil {
		panic(err)
	}

	// Set maximum points in second cluster to 1.
	cm2.SetMaximumPoints(1)
	// Exclude the second and third point from the
	// second cluster.
	err = cm2.SetExcludedPointIndices([]int{1, 2})

	if err != nil {
		panic(err)
	}

	// Print the maximum points in the first cluster.
	fmt.Println(cm1.MaximumPoints())
	// Print the excluded point indices in the first cluster.
	fmt.Println(cm1.ExcludedPointIndices())
	// Print the maximum sum value in the first cluster.
	fmt.Println(msv.MaximumValue())
	// Print the values of the points used in the maximum
	// sum value constraint.
	fmt.Println(msv.Values())
	// Print the maximum points in the second cluster.
	fmt.Println(cm2.MaximumPoints())
	// Print the excluded point indices in the second cluster.
	fmt.Println(cm2.ExcludedPointIndices())

	// Output:
	// 2
	// [0]
	// 10
	// [6 7 4]
	// 1
	// [1 2]
}
