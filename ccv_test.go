package ccv_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/smlx/ccv"
)

func TestNextVersion(t *testing.T) {
	var testCases = map[string]struct {
		gitCmds           [][]string
		expectVersion     string
		expectVersionType string
	}{
		"none": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "chore: not much"},
		},
			expectVersion:     "v0.1.0",
			expectVersionType: "",
		},
		"patch": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "fix: minor bug"},
		},
			expectVersion:     "v0.1.1",
			expectVersionType: "patch",
		},
		"minor": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "feat: cool new feature"},
		},
			expectVersion:     "v0.2.0",
			expectVersionType: "minor",
		},
		"major": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "feat: major refactor\nBREAKING CHANGE: new stuff"},
		},
			expectVersion:     "v1.0.0",
			expectVersionType: "major",
		},
		"major fix bang": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "fix!: major bug"},
		},
			expectVersion:     "v1.0.0",
			expectVersionType: "major",
		},
		"major feat bang": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "feat!: major change"},
		},
			expectVersion:     "v1.0.0",
			expectVersionType: "major",
		},
		"patch with scope": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "fix(lasers): minor bug"},
		},
			expectVersion:     "v0.1.1",
			expectVersionType: "patch",
		},
		"minor with scope": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "feat(phasers): cool new feature"},
		},
			expectVersion:     "v0.2.0",
			expectVersionType: "minor",
		},
		"major with scope": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "feat(blasters): major refactor\nBREAKING CHANGE: new stuff"},
		},
			expectVersion:     "v1.0.0",
			expectVersionType: "major",
		},
		"major fix bang with scope": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "fix(lightsaber)!: major bug"},
		},
			expectVersion:     "v1.0.0",
			expectVersionType: "major",
		},
		"major feat bang with scope": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"commit", "--allow-empty", "-m", "feat(bowcaster)!: major change"},
		},
			expectVersion:     "v1.0.0",
			expectVersionType: "major",
		},
		"no existing tags feat": {gitCmds: [][]string{
			{"commit", "--allow-empty", "-m", "feat: new change"},
		},
			expectVersion:     "v0.1.0",
			expectVersionType: "minor",
		},
		"no existing tags chore": {gitCmds: [][]string{
			{"commit", "--allow-empty", "-m", "chore: boring change"},
		},
			expectVersion:     "v0.1.0",
			expectVersionType: "minor",
		},
		"on a branch": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"checkout", "-b", "new-branch"},
			{"commit", "--allow-empty", "-m", "fix: minor change"},
		},
			expectVersion:     "v0.1.1",
			expectVersionType: "patch",
		},
		"tag on a branch": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"checkout", "-b", "new-branch"},
			{"commit", "--allow-empty", "-m", "fix: minor change"},
			{"tag", "v0.1.1"},
			{"checkout", "main"},
			{"commit", "--allow-empty", "-m", "feat: minor change"},
		},
			expectVersion:     "v0.2.0",
			expectVersionType: "minor",
		},
		"on a branch again": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"checkout", "-b", "new-branch"},
			{"commit", "--allow-empty", "-m", "fix: minor change"},
			{"tag", "v0.1.1"},
			{"checkout", "main"},
			{"commit", "--allow-empty", "-m", "feat: minor change"},
			{"tag", "v0.2.0"},
			{"commit", "--allow-empty", "-m", "fix: minor change"},
		},
			expectVersion:     "v0.2.1",
			expectVersionType: "patch",
		},
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
		},
			expectVersion:     "v0.1.2",
			expectVersionType: "patch",
		},
		"main after merge": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"checkout", "-b", "new-branch"},
			{"commit", "--allow-empty", "-m", "chore: boring change"},
			{"commit", "--allow-empty", "-m", "fix: minor change"},
			{"commit", "--allow-empty", "-m", "chore: boring change"},
			{"checkout", "main"},
			{"merge", "--no-ff", "new-branch", "-m", "chore: merge"},
		},
			expectVersion:     "v0.1.1",
			expectVersionType: "patch",
		},
		"branch after merge": {gitCmds: [][]string{
			{"tag", "v0.1.0"},
			{"checkout", "-b", "new-branch"},
			{"commit", "--allow-empty", "-m", "fix: minor change"},
			{"checkout", "main"},
			{"merge", "--no-ff", "new-branch", "-m", "chore: merge"},
			{"tag", "v0.1.2"},
			{"checkout", "-b", "new-branch-2"},
			{"commit", "--allow-empty", "-m", "feat: major change"},
		},
			expectVersion:     "v0.2.0",
			expectVersionType: "minor",
		},
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
		},
			expectVersion:     "v0.2.0",
			expectVersionType: "minor",
		},
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
		},
			expectVersion:     "v0.1.2",
			expectVersionType: "patch",
		},
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
				tt.Fatalf("error from main.NextVersion(): %v", err)
			}
			if next != tc.expectVersion {
				tt.Fatalf("expected: %v, got: %v", tc.expectVersion, next)
			}
			nextType, err := ccv.NextVersionType(dir)
			if err != nil {
				tt.Fatalf("error from main.NextVersionType(): %v", err)
			}
			if nextType != tc.expectVersionType {
				tt.Fatalf("expected: %v, got: %v", tc.expectVersionType, nextType)
			}
		})
	}
}
