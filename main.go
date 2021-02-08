package main

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"go.uber.org/zap"
)

// nextVersion returns a string containing the next version based on the state
// of the git repository in path.
func nextVersion(log *zap.Logger, path string) (string, error) {
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
	// walk commit hashes back from HEAD
	commits, err := r.Log(&git.LogOptions{Order: git.LogOrderDFSPost})
	if err != nil {
		return "", fmt.Errorf("couldn't get commits: %w", err)
	}
	var major, minor, patch bool
	var stopIter error = fmt.Errorf("stop commit iteration")
	var latestTag string
	err = commits.ForEach(func(c *object.Commit) error {
		if latestTag = tagRefs[c.Hash.String()]; latestTag != "" {
			return stopIter
		}
		// analyze commit message
		if strings.HasPrefix(c.Message, "fix: ") {
			patch = true
		}
		if strings.HasPrefix(c.Message, "feat: ") {
			minor = true
		}
		if strings.Contains(c.Message, "BREAKING CHANGE: ") {
			major = true
		}
		return nil
	})
	if (err != nil && err != stopIter) || latestTag == "" {
		return "", fmt.Errorf("couldn't determine latest tag: %w", err)
	}
	// found a tag: parse, increment, and return.
	latestVersion, err := semver.NewVersion(latestTag)
	if err != nil {
		return "", fmt.Errorf(`couldn't parse tag "%v": %w`, latestTag, err)
	}
	var newVersion semver.Version
	switch {
	case major:
		newVersion = latestVersion.IncMajor()
	case minor:
		newVersion = latestVersion.IncMinor()
	case patch:
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
	next, err := nextVersion(log, `.`)
	if err != nil {
		log.Fatal("couldn't get next version", zap.Error(err))
	}
	fmt.Println(next)
}
