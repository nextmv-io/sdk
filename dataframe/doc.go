/*
Package dataframe provides a general interface for managing tabular data
that support filtering, aggregation and data manipulation.

This package is beta and its API is subject to change.

This package includes a decoder for CSV files which can be used to create
dataframes from CSV files. The decoder is not a streaming decoder and will
load the entire CSV file into memory.

The following example shows how to create a dataframe from a CSV file:

	package main

	import ...

	var (
		PickupDay = dataframe.Strings("PICKUP_DAY")
		PickupTime = dataframe.Ints("PICKUP_TIME")
		// ...
	)
	func main() {
		run.Run(handler, run.Decode(dataframe.FromCSV))
	}

	func handler(d dataframe.DataFrame, o store.Options) (store.Solver, error) {
		d = d.Filter(
				PickupDay.Equals("Monday").Or(PickupDay.Equals("Tuesday")
			).Filter(
				PickupTime.IsInRange(16, 19),
			)
		// ...
	}

Running the example will create a dataframe from the CSV file and pass it
to the handler function. The handler function can then use the dataframe
to solve the problem. To pass a CSV file to the example use the
input path flag:

	nextmv run local . -- \
		-runner.input.path ./data.csv.gz \
		-runner.output.path output.json

Deprecated: This package is deprecated and will be removed in the future.
*/
package dataframe
