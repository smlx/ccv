package main_test

import (
	"io/ioutil"
	"os/exec"
	"testing"

	ccv "github.com/smlx/ccv"
)

func TestNextVersion(t *testing.T) {
	var testCases = map[string]struct {
		gitCmds [][]string
		expect  string
	}{
		"none": {gitCmds: [][]string{
			{"commit", "--allow-empty", "-m", "chore: not much"},
		}, expect: "v0.1.0"},
		"patch": {gitCmds: [][]string{
			{"commit", "--allow-empty", "-m", "fix: minor bug"},
		}, expect: "v0.1.1"},
		"minor": {gitCmds: [][]string{
			{"commit", "--allow-empty", "-m", "feat: cool new feature"},
		}, expect: "v0.2.0"},
		"major": {gitCmds: [][]string{
			{"commit", "--allow-empty", "-m", "feat: major refactor\nBREAKING CHANGE: new stuff"},
		}, expect: "v1.0.0"},
		"major fix bang": {gitCmds: [][]string{
			{"commit", "--allow-empty", "-m", "fix!: major bug"},
		}, expect: "v1.0.0"},
		"major feat bang": {gitCmds: [][]string{
			{"commit", "--allow-empty", "-m", "feat!: major change"},
		}, expect: "v1.0.0"},
		"patch with scope": {gitCmds: [][]string{
			{"commit", "--allow-empty", "-m", "fix(lasers): minor bug"},
		}, expect: "v0.1.1"},
		"minor with scope": {gitCmds: [][]string{
			{"commit", "--allow-empty", "-m", "feat(phasers): cool new feature"},
		}, expect: "v0.2.0"},
		"major with scope": {gitCmds: [][]string{
			{"commit", "--allow-empty", "-m", "feat(blasters): major refactor\nBREAKING CHANGE: new stuff"},
		}, expect: "v1.0.0"},
		"major fix bang with scope": {gitCmds: [][]string{
			{"commit", "--allow-empty", "-m", "fix(lightsaber)!: major bug"},
		}, expect: "v1.0.0"},
		"major feat bang with scope": {gitCmds: [][]string{
			{"commit", "--allow-empty", "-m", "feat(bowcaster)!: major change"},
		}, expect: "v1.0.0"},
	}
	for name, tc := range testCases {
		t.Run(name, func(tt *testing.T) {
			// create test dir
			dir, err := ioutil.TempDir("", "example")
			if err != nil {
				tt.Fatalf("couldn't get a tempdir: %v", err)
			}
			// init git repo
			initCmds := [][]string{
				{"init"},
				{"commit", "--allow-empty", "-m", "feat: initial commit"},
				{"tag", "v0.1.0"},
			}
			for _, c := range initCmds {
				cmd := exec.Command("git", c...)
				cmd.Dir = dir
				if output, err := cmd.CombinedOutput(); err != nil {
					tt.Fatalf("couldn't run init cmd git %v: %v (%s)", c, err, output)
				}
			}
			for _, c := range tc.gitCmds {
				cmd := exec.Command("git", c...)
				cmd.Dir = dir
				if output, err := cmd.CombinedOutput(); err != nil {
					tt.Fatalf("couldn't run git %v: %v (%s)", c, err, output)
				}
			}
			next, err := ccv.NextVersion(dir)
			if err != nil {
				tt.Fatalf("error from main.nextVersion(): %v", err)
			}
			if next != tc.expect {
				tt.Fatalf("expected: %v, got: %v", tc.expect, next)
			}
		})
	}
}
