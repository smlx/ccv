# Conventional Commits Versioner

[![Release](https://github.com/smlx/ccv/actions/workflows/release.yaml/badge.svg)](https://github.com/smlx/ccv/actions/workflows/release.yaml)
[![coverage](https://raw.githubusercontent.com/smlx/ccv/badges/.badges/main/coverage.svg)](https://github.com/smlx/ccv/actions/workflows/coverage.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/smlx/ccv)](https://goreportcard.com/report/github.com/smlx/ccv)

`ccv` does one thing: it walks git commit history back from the current `HEAD` to find the most recent tag, taking note of commit messages along the way.
When it reaches the most recent tag, it uses the commit messages it saw to figure out how the tag should be incremented, and prints the incremented tag.

`ccv` is intended for use in continuous delivery automation.

The ideas behind `ccv` are described by [Conventional Commits](https://www.conventionalcommits.org/) and [Semantic Versioning](https://semver.org/). Currently parts 1 to 3 of the Conventional Commits specification summary are recognized when incrementing versions.

## Use as a Github Action

This repository is also a [Github Action](https://docs.github.com/en/actions).

Inputs:

* `write-tag`: If true, and ccv determines that a new version is required, the action will automatically write the new version tag to the repository. Default `true`.

Outputs:

* `new-tag`: Either "true" or "false" depending on whether a new tag was pushed. Example: `true`.
* `new-tag-version`: The new version that was tagged. This will only be set if new_tag=true. Example: `v0.1.2`.

### Example: automatic tagging

The main use-case of this action is to automatically tag and build new releases in a fully automated release workflow.

```yaml
name: release
on:
  push:
    branches:
    - main
permissions: {}
jobs:
  release-tag:
    permissions:
      # create tag
      contents: write
    runs-on: ubuntu-latest
    outputs:
      new-tag: ${{ steps.ccv.outputs.new-tag }}
    steps:
    - uses: actions/checkout@0ad4b8fadaa221de15dcec353f45205ec38ea70b # v4.1.4
      with:
        fetch-depth: 0
    - name: Bump tag if necessary
      id: ccv
      uses: smlx/ccv@c5f6769c943c082c4e8d8ccf2ec4b6f5f517e1f2 # v0.7.3
  release-build:
    permissions:
      # create release
      contents: write
      # push docker images to registry
      packages: write
    needs: release-tag
    if: needs.release-tag.outputs.new-tag == 'true'
    runs-on: ubuntu-latest
    steps:
    # ... build and release steps here
```

For a fully-functional example, see the [release workflow of this repository](https://github.com/smlx/ccv/blob/main/.github/workflows/release.yaml).

### Example: read-only

You can also check the tag your PR will generate by running with `write-tag: false`. Note that the permissions on this job are read-only.

```yaml
name: build
on:
  pull_request:
    branches:
    - main
permissions: {}
jobs:
  check-tag:
    permissions:
      contents: read
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@0ad4b8fadaa221de15dcec353f45205ec38ea70b # v4.1.4
      with:
        fetch-depth: 0
    - id: ccv
      uses: smlx/ccv@c5f6769c943c082c4e8d8ccf2ec4b6f5f517e1f2 # v0.7.3
      with:
        write-tag: false
    - run: |
        echo "new-tag=$NEW_TAG"
        echo "new-tag-version=$NEW_TAG_VERSION"
      env:
        NEW_TAG: ${{steps.ccv.outputs.new-tag}}
        NEW_TAG_VERSION: ${{steps.ccv.outputs.new-tag-version}}
```

Gives this output:

```
new-tag=true
new-tag-version=v0.16.0
```

For a fully-functional example, see the [build workflow of this repository](https://github.com/smlx/ccv/blob/main/.github/workflows/build.yaml).

## Use locally

Download the latest [release](https://github.com/smlx/ccv/releases) on github, or:

```
go install github.com/smlx/ccv/cmd/ccv@latest
```

Run `ccv` in the directory containing your git repository.

## Prior art

* [caarlos0/svu](https://github.com/caarlos0/svu) does pretty much the same thing, but it has more features and shells out to git. `ccv` uses [go-git/go-git](https://github.com/go-git/go-git) instead.
