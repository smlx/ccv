package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"go.uber.org/zap"
)

func main() {
	// init logger
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	// open repository
	r, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		log.Fatal("couldn't open git repository", zap.Error(err))
	}
	tags, err := r.Tags()
	if err != nil {
		log.Fatal("couldn't get tags", zap.Error(err))
	}
	// map tags to commit hashes
	tagRefs := map[string]string{}
	err = tags.ForEach(func(r *plumbing.Reference) error {
		tagRefs[r.Hash().String()] = r.Name().Short()
		return nil
	})
	if err != nil {
		log.Fatal("couldn't iterate tags", zap.Error(err))
	}
	// walk commit hashes back from HEAD
	commits, err := r.Log(&git.LogOptions{})
	if err != nil {
		log.Fatal("couldn't get commits", zap.Error(err))
	}
	var major, minor, patch bool
	commits.ForEach(func(c *object.Commit) error {
		if tag := tagRefs[c.Hash.String()]; tag != "" {
			// found a tag
			t, err := semver.NewVersion(tag)
			if err != nil {
				log.Fatal("couldn't parse tag", zap.Error(err), zap.String("tag", tag))
			}
			var v semver.Version
			switch {
			case major:
				v = t.IncMajor()
			case minor:
				v = t.IncMinor()
			case patch:
				v = t.IncPatch()
			}
			fmt.Printf("%s%s\n", "v", v.String())
			os.Exit(0)
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
}
