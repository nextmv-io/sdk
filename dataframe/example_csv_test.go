package dataframe_test

import (
	"fmt"
	"github.com/nextmv-io/sdk/dataframe"
	"reflect"
	"strings"
)

func ExampleFromCSV() {
	//TODO remove this code
	input := `Duration,Pulse,Maxpulse,Calories,Name
		60,110,130,409.1, "A"
		70,100,120,408.1, "B"`

	decoder := dataframe.FromCSV()

	var df dataframe.DataFrame

	dataFrameValueAsInterface := reflect.New(reflect.TypeOf(&df).
		Elem()).
		Interface()

	err := decoder.Decode(strings.NewReader(input), dataFrameValueAsInterface)

	if err != nil {
		panic(err)
	}

	df = *(dataFrameValueAsInterface).(*dataframe.DataFrame)

	columns := df.Columns()

	for _, c := range columns {
		fmt.Println(c.Name())
	}
	// Unordered output:
	// Duration
	// Pulse
	// Maxpulse
	// Calories
	// Name
}
