#!/bin/sh
set -eu
# if the first argument to the script is "true", it will push a tag to the repository.
WRITE_TAG="$1"
# the runner workspace will be mounted here, and git complains otherwise
git config --global --add safe.directory /github/workspace
# if the ccv tag exists, just exit
if [ "$(git tag -l "$(ccv)")" ]; then
	echo "new-tag=false" >>"$GITHUB_OUTPUT"
	exit
fi
# if it doesn't, tag and push
if [ "$WRITE_TAG" = "true" ]; then
	git tag "$(ccv)"
	git push --tags
fi
{
	echo "new-tag=true"
	echo "new-tag-version=$(ccv)"
	echo "new-tag-version-type=$(ccv --version-type)"
} >>"$GITHUB_OUTPUT"
