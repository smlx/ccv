package main

import (
	"fmt"
	"regexp"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"go.uber.org/zap"
)

var patchRegex = regexp.MustCompile(`^fix(\(.+\))?: `)
var minorRegex = regexp.MustCompile(`^feat(\(.+\))?: `)
var majorRegex = regexp.MustCompile(`^(fix|feat)(\(.+\))?!: |BREAKING CHANGE: `)

// walkCommits walks the git history in the defined order until it reaches a
// tag, analysing the commits it finds.
func walkCommits(r *git.Repository, tagRefs map[string]string, order git.LogOrder) (*semver.Version, bool, bool, bool, error) {
	var major, minor, patch bool
	var stopIter error = fmt.Errorf("stop commit iteration")
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

// NextVersion returns a string containing the next version based on the state
// of the git repository in path.
func NextVersion(path string) (string, error) {
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
	if latestMain == nil {
		latestVersion = latestBranch
	} else if latestBranch == nil {
		latestVersion = latestMain
	} else if latestMain.GreaterThan(latestBranch) {
		latestVersion = latestMain
	} else {
		latestVersion = latestBranch
	}
	// figure out the highest increment in either parent
	var newVersion semver.Version
	switch {
	case majorMain || majorBranch:
		newVersion = latestVersion.IncMajor()
	case minorMain || minorBranch:
		newVersion = latestVersion.IncMinor()
	case patchMain || patchBranch:
		newVersion = latestVersion.IncPatch()
	default:
		newVersion = *latestVersion
	}
	return fmt.Sprintf("%s%s", "v", newVersion.String()), nil
}

func main() {
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	next, err := NextVersion(`.`)
	if err != nil {
		log.Fatal("couldn't get next version", zap.Error(err))
	}
	fmt.Println(next)
}
