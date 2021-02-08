package main_test

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"testing"

	ccv "github.com/smlx/ccv"
	"go.uber.org/zap"
)

func TestNextVersion(t *testing.T) {
	var testCases = map[string]struct {
		gitCmds [][]string
		expect  string
	}{
		"none": {gitCmds: [][]string{
			{"commit", "--allow-empty", "-m", "feat: initial commit"},
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "chore: not much"},
		}, expect: "v0.1.0"},
		"patch": {gitCmds: [][]string{
			{"commit", "--allow-empty", "-m", "feat: initial commit"},
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "fix: minor bug"},
		}, expect: "v0.1.1"},
		"minor": {gitCmds: [][]string{
			{"commit", "--allow-empty", "-m", "feat: initial commit"},
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "feat: cool new feature"},
		}, expect: "v0.2.0"},
		"major": {gitCmds: [][]string{
			{"commit", "--allow-empty", "-m", "feat: initial commit"},
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "feat: major refactor\nBREAKING CHANGE: new stuff"},
		}, expect: "v1.0.0"},
	}
	log, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("couldn't get logger: %v", err)
	}
	for name, tc := range testCases {
		t.Run(name, func(tt *testing.T) {
			// create test dir
			dir, err := ioutil.TempDir("", "example")
			if err != nil {
				tt.Fatalf("couldn't get a tempdir: %v", err)
			}
			// init git repo
			initCmd := exec.Command("git", "init")
			initCmd.Dir = dir
			initCmd.Env = []string{fmt.Sprintf("HOME=%s", dir)}
			if err = initCmd.Run(); err != nil {
				tt.Fatalf("couldn't git init: %v", err)
			}
			for _, c := range tc.gitCmds {
				cmd := exec.Command("git", c...)
				cmd.Dir = dir
				cmd.Env = []string{fmt.Sprintf("HOME=%s", dir)}
				if err = cmd.Run(); err != nil {
					tt.Fatalf("couldn't run git command `%s`: %v", c, err)
				}
			}
			tt.Log(dir)
			next, err := ccv.NextVersion(log, dir)
			if err != nil {
				tt.Fatalf("error from main.nextVersion(): %v", err)
			}
			if next != tc.expect {
				tt.Fatalf("expected: %v, got: %v", tc.expect, next)
			}
		})
	}
}
