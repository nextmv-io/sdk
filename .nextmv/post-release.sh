export MODULES=$(cat "$(dirname "$0")/workflow-configuration.yml" | yq '.nested_modules[]' -r)

git checkout -b feature/bump-nested-$VERSION

for module in $MODULES; do
    echo "Bumping $module to $VERSION"
    pushd $module
    go get github.com/nextmv-io/sdk@$VERSION
    go mod tidy
    git add go.mod go.sum
    popd
done

git commit -S -m "Bump nested modules after $VERSION release"
git push origin --set-upstream feature/bump-nested-$VERSION

OUTPUT=$(gh pr create --base $BRANCH --title "Bump nested modules after $VERSION release" --body "Automated bump of nested modules to version $VERSION")

echo "# :rocket: PR created" >> $GITHUB_STEP_SUMMARY
echo "" >> $GITHUB_STEP_SUMMARY
echo "Bump nested modules :arrow_right: [PR Link](${OUTPUT})" >> $GITHUB_STEP_SUMMARY
