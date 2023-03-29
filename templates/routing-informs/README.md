# Nextmv routing template

We are in the business of connecting customers in urban areas of Colorado to
local farms who offer farm shares (1 farm share = 1 order). We currently have a
fleet of vans that pick up orders from farms and deliver those orders to
customers’ homes. As a business committed to the environment, we are working to
reduce emissions from our fleet by reducing our reliance on vans and adding
bicycles to our fleet. Our goal is to make use of bicycles whenever we can still
meet our delivery time windows, and only utilize vans when that’s not possible.

The objective of the exercise is to use the Nextmv platform to help the business
decide how much of their fleet they can convert to bicycles in the Colorado
markets they operate in: Denver and Aurora.

`main.go` implements a VRP solver with the features needed for this problem
already configured. The input files hold different data for different markets
and sizes.

Before you start customizing run the command below to see if everything works as
expected:

```bash
nextmv sdk run . -- -runner.input.path input.json\
  -runner.output.path output.json -limits.duration 10s
```

A file `output.json` should have been created with a VRP solution.

## Next steps

* For more information about our platform, please visit: <https://docs.nextmv.io>.
* Need more assistance? Send us an [email](mailto:support@nextmv.io)!
