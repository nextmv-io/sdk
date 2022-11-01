# Routing Example Inputs

Running the Routing app requires access to Nextmv's remote platform. If you do
not have access, you can [contact support][support] to request access.

The routing example inputs in this folder are to be used with the Nextmv
Routing app and are meant to give examples of how to use the different features
offered in that app for varying use cases.

Below is a list and a short description of features utilized along with an
explanation of use case.

## General

General routing examples

- [fleet_tiny.json](./fleet_tiny.json) - A basic fleet with 2 vehicles and 10
stops.
- [fleet_pd.json](./fleet_pd.json) - A larger fleet utilizing
[pickup and delivery precedence](https://docs.nextmv.io/cloud/features/precedence).

## Delivery

Routes that include pickups and dropoffs. Vehicles can start and end anywhere.
Example use cases: meal delivery, ridesharing

- [delivery-tiny.json](./delivery-tiny.json) - A small fleet of 2 vehicles and
10 stops with no starting and ending location and delivery precedence.
- [delivery-advanced.json](./delivery-advanced.json) - A delivery fleet input
utilizing
[compatability attributes](https://docs.nextmv.io/cloud/features/compatibility-attributes)
and [target times](https://docs.nextmv.io/cloud/features/time-settings#target-times).

## Distribution

Routes include dropoffs only. All vehichles start from the same location.
Sample use cases: package/home goods distribution

- [distribution-tiny.json](./distribution-tiny.json) - A small fleet of vehicles
starting from a central location.
- [distribution-route-limit.json](./distribution-route-limit.json) - A fleet of
vehicles starting from a central location and utilizing the
[max distance](https://docs.nextmv.io/cloud/features/route-limits#max-distance)
feature.

## Sourcing

Routes include pickups only. Vehichles end at the same location. Sample use
cases: agriculture sourcing, waste management

- [sourcing-tiny.json](./sourcing-tiny.json) - A small fleet with the same
ending location.

For a complete list of features, please see
[our docs][docs].

[support]: https://www.nextmv.io/contact
[docs]: https://docs.nextmv.io/cloud/features
