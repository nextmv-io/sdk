# Release

A reusable workflow is used to release the package. Nextmv team members: please
go to the corresponding repository for more information.

## Stable release

Open a PR against the `develop` branch with the following change:

* Update the version in the `VERSION` file.

After the PR is merged, the `release.yml` workflow will be triggered and it
will automatically create a release.

After the release is created, a new PR will be opened against `develop` bumping
the nested modules to the new version. Make sure you approve and merge this new
PR.

## Pre-release

Update the version in the `VERSION` file to a dev tag. When a commit is pushed,
the `release.yml` workflow will be triggered and it will automatically create a
release.

After the release is created, a new PR will be opened against `develop` bumping
the nested modules to the new version. Make sure you approve and merge this new
PR.
