#!/usr/bin/env bash
# Checks that every nested module requires the sdk version in VERSION and is
# tagged by the release workflow. Nested modules are released together with
# the sdk (same commit, tags <module path>/<version>).
#
# Usage: .nextmv/check-nested-modules.sh [--fix]
set -euo pipefail

cd "$(dirname "$0")/.."
VERSION="$(cat VERSION)"
FIX="${1:-}"
FAILED=0

for gomod in $(find . -mindepth 2 -name go.mod -not -path "./.git/*" | sort); do
    module="$(dirname "${gomod#./}")"

    required="$(go mod edit -json "$gomod" | jq -r '.Require[]? | select(.Path == "github.com/nextmv-io/sdk") | .Version')"
    if [ "$required" != "$VERSION" ]; then
        if [ "$FIX" = "--fix" ]; then
            go mod edit -require "github.com/nextmv-io/sdk@$VERSION" "$gomod"
            echo "$module: bumped github.com/nextmv-io/sdk $required -> $VERSION"
        else
            echo "::error file=${gomod#./}::$module requires github.com/nextmv-io/sdk $required, expected $VERSION (run .nextmv/check-nested-modules.sh --fix)"
            FAILED=1
        fi
    fi

    if ! grep -qE "EXTRA_TAG_PREFIXES:.*(^|[ \"'])$module/([ \"']|$)" .github/workflows/release.yml; then
        echo "::error file=.github/workflows/release.yml::$module/ is missing from EXTRA_TAG_PREFIXES"
        FAILED=1
    fi
done

exit $FAILED
