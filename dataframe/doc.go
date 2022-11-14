/*
Package dataframe provides a general interface for managing tabular data
that support filtering, aggregation and data manipulation.

This package is beta and its API is subject to change.

This package includes a decoder for CSV files which can be used to create
dataframes from CSV files. The decoder is not a streaming decoder and will
load the entire CSV file into memory.

In addition there exist an Apache Feather decoder which can be used to read
dataframes from Apache Arrow Feather files (IPC format). Like the CSV decoder
this decoder is not a streaming decoder and will load the entire file into
memory.

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
		-hop.runner.input.path ./data.csv.gz \
		-hop.runner.output.path output.json
*/
package dataframe
