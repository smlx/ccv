package ccv_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/smlx/ccv"
)

func TestNextVersion(t *testing.T) {
	var testCases = map[string]struct {
		gitCmds [][]string
		expect  string
	}{
		"none": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "chore: not much"},
		}, expect: "v0.1.0"},
		"patch": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "fix: minor bug"},
		}, expect: "v0.1.1"},
		"minor": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "feat: cool new feature"},
		}, expect: "v0.2.0"},
		"major": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "feat: major refactor\nBREAKING CHANGE: new stuff"},
		}, expect: "v1.0.0"},
		"major fix bang": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "fix!: major bug"},
		}, expect: "v1.0.0"},
		"major feat bang": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "feat!: major change"},
		}, expect: "v1.0.0"},
		"patch with scope": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "fix(lasers): minor bug"},
		}, expect: "v0.1.1"},
		"minor with scope": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "feat(phasers): cool new feature"},
		}, expect: "v0.2.0"},
		"major with scope": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "feat(blasters): major refactor\nBREAKING CHANGE: new stuff"},
		}, expect: "v1.0.0"},
		"major fix bang with scope": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "fix(lightsaber)!: major bug"},
		}, expect: "v1.0.0"},
		"major feat bang with scope": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "feat(bowcaster)!: major change"},
		}, expect: "v1.0.0"},
		"no existing tags feat": {gitCmds: [][]string{
			{"commit", "--allow-empty", "-m", "feat: new change"},
		}, expect: "v0.1.0"},
		"no existing tags chore": {gitCmds: [][]string{
			{"commit", "--allow-empty", "-m", "chore: boring change"},
		}, expect: "v0.1.0"},
		"on a branch": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"checkout", "-b", "new-branch"},
			{"commit", "--allow-empty", "-m", "fix: minor change"},
		}, expect: "v0.1.1"},
		"tag on a branch": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"checkout", "-b", "new-branch"},
			{"commit", "--allow-empty", "-m", "fix: minor change"},
			{"tag", "v0.1.1"},
			{"checkout", "main"},
			{"commit", "--allow-empty", "-m", "feat: minor change"},
		}, expect: "v0.2.0"},
		"on a branch again": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"checkout", "-b", "new-branch"},
			{"commit", "--allow-empty", "-m", "fix: minor change"},
			{"tag", "v0.1.1"},
			{"checkout", "main"},
			{"commit", "--allow-empty", "-m", "feat: minor change"},
			{"tag", "v0.2.0"},
			{"commit", "--allow-empty", "-m", "fix: minor change"},
		}, expect: "v0.2.1"},
		"back on a branch": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"checkout", "-b", "new-branch"},
			{"commit", "--allow-empty", "-m", "fix: minor change"},
			{"tag", "v0.1.1"},
			{"checkout", "main"},
			{"commit", "--allow-empty", "-m", "feat: minor change"},
			{"tag", "v0.2.0"},
			{"checkout", "new-branch"},
			{"commit", "--allow-empty", "-m", "fix: minor change"},
		}, expect: "v0.1.2"},
		"main after merge": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"checkout", "-b", "new-branch"},
			{"commit", "--allow-empty", "-m", "chore: boring change"},
			{"commit", "--allow-empty", "-m", "fix: minor change"},
			{"commit", "--allow-empty", "-m", "chore: boring change"},
			{"checkout", "main"},
			{"merge", "--no-ff", "new-branch", "-m", "chore: merge"},
		}, expect: "v0.1.1"},
		"branch after merge": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"checkout", "-b", "new-branch"},
			{"commit", "--allow-empty", "-m", "fix: minor change"},
			{"checkout", "main"},
			{"merge", "--no-ff", "new-branch", "-m", "chore: merge"},
			{"tag", "v0.1.2"},
			{"checkout", "-b", "new-branch-2"},
			{"commit", "--allow-empty", "-m", "feat: major change"},
		}, expect: "v0.2.0"},
		"main after merge again": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"checkout", "-b", "new-branch"},
			{"commit", "--allow-empty", "-m", "fix: minor change"},
			{"checkout", "main"},
			{"merge", "--no-ff", "new-branch", "-m", "chore: merge"},
			{"tag", "v0.1.2"},
			{"checkout", "-b", "new-branch-2"},
			{"commit", "--allow-empty", "-m", "feat: major change"},
			{"checkout", "main"},
			{"merge", "--no-ff", "new-branch-2", "-m", "chore: merge"},
		}, expect: "v0.2.0"},
		"branch before tag and merge": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"checkout", "-b", "new-branch-1"},
			{"checkout", "-b", "new-branch-2"},
			{"commit", "--allow-empty", "-m", "fix: minor change"},
			{"checkout", "main"},
			{"merge", "--no-ff", "new-branch-2", "-m", "chore: merge"},
			{"tag", "v0.1.1"},
			{"checkout", "new-branch-1"},
			{"commit", "--allow-empty", "-m", "fix: another minor change"},
			{"checkout", "main"},
			{"merge", "--no-ff", "new-branch-1", "-m", "chore: merge"},
		}, expect: "v0.1.2"},
	}
	for name, tc := range testCases {
		t.Run(name, func(tt *testing.T) {
			// create test dir
			dir, err := os.MkdirTemp("", "example")
			if err != nil {
				tt.Fatalf("couldn't get a tempdir: %v", err)
			}
			// init git repo
			initCmds := [][]string{
				{"init", "-b", "main"},
				{"commit", "--allow-empty", "-m", "feat: initial commit"},
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
