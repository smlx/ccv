#!/bin/sh
set -eu
# https://docs.github.com/en/actions/creating-actions/creating-a-docker-container-action#accessing-files-created-by-a-container-action
cd /github/workspace
# if the ccv tag exists, just exit
if [ "$(git tag -l "$(ccv)")" ]; then
	echo "new_tag=false" >>"$GITHUB_OUTPUT"
	exit
fi
# if it doesn't, tag and push
git tag "$(ccv)"
git push --tags
echo "new_tag=true" >>"$GITHUB_OUTPUT"
echo "new_tag_version=$(git tag --points-at HEAD)" >>"$GITHUB_OUTPUT"
