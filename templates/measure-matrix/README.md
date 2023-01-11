# Nextmv measure matrix input generation for the routing template

`measure-matrix` serves as an input generator for the `routing-matrix-input`
app. The generated input data includes matrices for distance and travel time
using OSRM. You can simply exchange the client for one of our other matrix
providers: Routingkit, Google and Here.

To run both apps combined, you should look at two things:

* add the OSRM server to `main.go` (or use a different provider of your choice)
* use locations in the input data in `main.go` that can be consumed by your OSRM
  server.
  
Run the command below to see if everything works as
expected:

```bash
nextmv run local .
```

One file should have been created: `routing-input.json`. This file can be used
as a direct input for the `routing-matrix-input` app.

## Next steps

* For more information about our platform, please visit: <https://docs.nextmv.io>.
* Need more assistance? Send us an [email](mailto:support@nextmv.io)!
