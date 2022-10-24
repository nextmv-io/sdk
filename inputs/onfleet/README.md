# Routing with Onfleet Integration Example Inputs

Running the Routing app requires access to Nextmv's remote platform. If you do
not have access, you can [contact support][support] to request access.

It is recommended that you start with [our docs][docs] before using this
feature.

In order to use the Onfleet Integration with our Routing app, you must first
set up a run profile with an attached Onfleet integration. You can follow the
steps
[here](https://docs.nextmv.io/cloud/integrations/onfleet#using-the-integration)
or [create an integration on our
console](https://cloud.nextmv.io/config/integrations/nextmv-routing) and then
[attach it to a run profile](https://cloud.nextmv.io/config/run-profiles).

This run profile id should replace the id in the
[example input json](./onfleet.json).

Additionally, you must create [tasks](https://docs.onfleet.com/reference/tasks)
and [workers](https://docs.onfleet.com/reference/workers) using OnFleet and
include their ids in the [example input json](./onfleet.json).

[support]: https://www.nextmv.io/contact
[docs]: https://docs.nextmv.io/cloud/features
