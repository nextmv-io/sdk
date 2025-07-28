export MODULES=$(cat workflow-configuration.yml | yq '.nested_modules[]' -r)

for module in $MODULES; do
    echo "Bumping $module to $VERSION"
    pushd ../$module
    go get github.com/nextmv-io/sdk@$VERSION
    go mod tidy
    popd
done

git add --all
git checkout -b feature/bump-nested-$VERSION
git commit -S -m "Bump nested modules after $VERSION release"
git push origin --set-upstream feature/bump-nested-$VERSION
gh pr create --base develop --title "Bump nested modules after $VERSION release" --body "Automated bump of nested modules to version $VERSION"
