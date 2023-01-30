# Nextmv PagerDuty scheduling template

The PagerDuty scheduling problem assigns users to shifts (days) for incident
response on-call duty. In our example problem we have the following setting:

* We have a number of days to schedule.
* We have a number of users. Each day should have a user, but not every
  user needs to be assigned to a day.
* In addition users can be unavailable for certain full days.
* Users can have preferences as to what days they prefer.
* Our objective is to find a plan that balances on-call responsibility across
  users and maximizes happiness. Happiness is measured as the number of times
  the user had their preferred day assigned.

This template uses our *custom modelling* framework to model and solve
such a scheduling problem.

The most important files created are `main.go` and `input.json`.

* `main.go` implements a shift scheduling solver for PagerDuty.
* `input.json` is a sample input file that follows the input definition in
`main.go`.

Before you start customizing, run the command below to see if everything works
as expected:

```bash
nextmv sdk run . -- -runner.input.path input.json \
  -runner.output.path output.json -limits.duration 10s
```

A file `output.json` should have been created with the best found schedule.

## Next steps

* For more information about our platform, please visit: <https://docs.nextmv.io>.
* Need more assistance? Send us an [email](mailto:support@nextmv.io)!
