package kmeans_test

import (
	"fmt"

	"github.com/nextmv-io/sdk/cluster/kmeans"
	"github.com/nextmv-io/sdk/measure"
)

func ExampleModel() {
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
	// Print the number of points in the model.
	fmt.Println(len(model.Points()))
	// Print the number of cluster models in the model.
	fmt.Println(len(model.ClusterModels()))

	// Output:
	// 3
	// 2
}
