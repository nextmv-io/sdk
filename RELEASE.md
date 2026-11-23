# Release

A reusable workflow is used to release the package. Nextmv team members: please
go to the corresponding repository for more information.

The nested modules (`measure/*`) are released together with the sdk: the same
commit is tagged `vX.Y.Z` and `measure/<module>/vX.Y.Z`. Bump the version with:

```bash
echo vX.Y.Z > VERSION
.nextmv/check-nested-modules.sh --fix
```

The `nested modules` check fails if a nested module does not require the sdk
version in `VERSION`.

## Stable release

Open a PR against the `develop` branch with the version bump above. After the
PR is merged, the `release.yml` workflow will be triggered and it will
automatically create a release.

## Pre-release

Apply the version bump above with a dev tag. When a commit is pushed, the
`release.yml` workflow will be triggered and it will automatically create a
release.
