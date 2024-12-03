// Package ccv implements the conventional commits versioner logic.
package ccv

import (
	"fmt"
	"regexp"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var patchRegex = regexp.MustCompile(`^fix(\(.+\))?: `)
var minorRegex = regexp.MustCompile(`^feat(\(.+\))?: `)
var majorRegex = regexp.MustCompile(`^(fix|feat)(\(.+\))?!: |BREAKING CHANGE: `)

// walkCommits walks the git history in the defined order until it reaches a
// tag, analysing the commits it finds.
func walkCommits(r *git.Repository, tagRefs map[string]string, order git.LogOrder) (*semver.Version, bool, bool, bool, error) {
	var major, minor, patch bool
	var stopIter = fmt.Errorf("stop commit iteration")
	var latestTag string
	// walk commit hashes back from HEAD via main
	commits, err := r.Log(&git.LogOptions{Order: order})
	if err != nil {
		return nil, false, false, false, fmt.Errorf("couldn't get commits: %w", err)
	}
	err = commits.ForEach(func(c *object.Commit) error {
		if latestTag = tagRefs[c.Hash.String()]; latestTag != "" {
			return stopIter
		}
		// analyze commit message
		if patchRegex.MatchString(c.Message) {
			patch = true
		}
		if minorRegex.MatchString(c.Message) {
			minor = true
		}
		if majorRegex.MatchString(c.Message) {
			major = true
		}
		return nil
	})
	if err != nil && err != stopIter {
		return nil, false, false, false,
			fmt.Errorf("couldn't determine latest tag: %w", err)
	}
	// not tagged yet. this can happen if we are on a branch with no tags.
	if latestTag == "" {
		return nil, false, false, false, nil
	}
	// found a tag: parse, increment, and return.
	latestVersion, err := semver.NewVersion(latestTag)
	if err != nil {
		return nil, false, false, false,
			fmt.Errorf(`couldn't parse tag "%v": %w`, latestTag, err)
	}
	return latestVersion, major, minor, patch, nil
}

// NextVersion returns a string containing the next version number based on the
// state of the git repository in path. It inspects the most recent tag, and
// the commits made after that tag.
func NextVersion(path string) (string, error) {
	return nextVersion(path, false)
}

// NextVersionType returns a string containing the next version type (major,
// minor, patch) based on the state of the git repository in path. It inspects
// the most recent tag, and the commits made after that tag.
func NextVersionType(path string) (string, error) {
	return nextVersion(path, true)
}

// nextVersion returns a string containing either the next version number, or
// the next version type (major, minor, patch) based on the state of the git
// repository in path. It inspects the most recent tag, and the commits made
// after that tag.
func nextVersion(path string, versionType bool) (string, error) {
	// open repository
	r, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return "", fmt.Errorf("couldn't open git repository: %w", err)
	}
	tags, err := r.Tags()
	if err != nil {
		return "", fmt.Errorf("couldn't get tags: %w", err)
	}
	// map tags to commit hashes
	tagRefs := map[string]string{}
	err = tags.ForEach(func(r *plumbing.Reference) error {
		tagRefs[r.Hash().String()] = r.Name().Short()
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("couldn't iterate tags: %w", err)
	}
	if len(tagRefs) == 0 {
		// no existing tags
		if versionType {
			return "minor", nil
		}
		return "v0.1.0", nil
	}
	// now we check both main and branch to figure out what the tag should be.
	// this logic is required for branches which split before the latest tag on
	// main. See the "branch before tag and merge" test.
	latestMain, majorMain, minorMain, patchMain, err :=
		walkCommits(r, tagRefs, git.LogOrderDFS)
	if err != nil {
		return "", fmt.Errorf("couldn't walk commits on main: %w", err)
	}
	latestBranch, majorBranch, minorBranch, patchBranch, err :=
		walkCommits(r, tagRefs, git.LogOrderDFSPost)
	if err != nil {
		return "", fmt.Errorf("couldn't walk commits on branch: %w", err)
	}
	if latestMain == nil && latestBranch == nil {
		return "",
			fmt.Errorf("tags exist in the repository, but not in ancestors of HEAD")
	}
	// figure out the latest version in either parent
	var latestVersion *semver.Version
	switch {
	case latestMain == nil:
		latestVersion = latestBranch
	case latestBranch == nil || latestMain.GreaterThan(latestBranch):
		latestVersion = latestMain
	default:
		latestVersion = latestBranch
	}
	// figure out the highest increment in either parent
	var newVersion semver.Version
	var newVersionType string
	switch {
	case majorMain || majorBranch:
		newVersion = latestVersion.IncMajor()
		newVersionType = "major"
	case minorMain || minorBranch:
		newVersion = latestVersion.IncMinor()
		newVersionType = "minor"
	case patchMain || patchBranch:
		newVersion = latestVersion.IncPatch()
		newVersionType = "patch"
	default:
		newVersion = *latestVersion
	}
	if versionType {
		return newVersionType, nil
	}
	return fmt.Sprintf("%s%s", "v", newVersion.String()), nil
}
