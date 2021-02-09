# Conventional Commits Versioner

![Tag and Release](https://github.com/smlx/ccv/workflows/Tag%20and%20Release/badge.svg)

`ccv` does one thing: it walks git commit history back from the current `HEAD` to find the most recent tag, taking note of commit messages along the way.
When it reaches the most recent tag, it uses the commit messages it saw to figure out how the tag should be incremented, and prints the incremented tag.

`ccv` is intended for use in continuous delivery automation.

The ideas behind `ccv` are described by [Conventional Commits](https://www.conventionalcommits.org/) and [Semantic Versioning](https://semver.org/). Currently parts 1 to 3 of the Conventional Commits specification summary are recognized when incrementing versions.

## Get it

```
go get github.com/smlx/ccv
```

## Use it

For a full example, see the [`tag-release` workflow](https://github.com/smlx/ccv/blob/main/.github/workflows/tag-release.yaml) in this repository.

Simple example:

```
# add an incremented tag if necessary
if [ -z $(git tag -l $(ccv)) ]; then
	git tag $(ccv)
fi
```

`ccv` takes no arguments or options\*.

\* Yet!

## Prior art

* [caarlos0/svu](https://github.com/caarlos0/svu) does pretty much the same thing, but it has more features and shells out to git. `ccv` uses [go-git/go-git](https://github.com/go-git/go-git) instead.
