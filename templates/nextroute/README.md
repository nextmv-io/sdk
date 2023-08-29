# Nextmv nextroute template

`nextroute` is a modeling kit for vehicle routing problems (VRP). This template
will get you up to speed deploying your own solution.

The most important files created are `main.go` and `input.json`.

`main.go` implements a VRP solver with many real world features already
configured. `input.json` is a sample input file that follows the input
definition in `main.go`.

You should be able to run the following command. It assumes that you gave your
app the app-id `shift-scheduling`:

```bash
nextmv app run -a "shift-scheduling" --options max_hours_per_day="9h"
```

## Next steps

* For more information about our platform, please visit: <https://docs.nextmv.io>.
* Need more assistance? Send us an [email](mailto:support@nextmv.io)!
