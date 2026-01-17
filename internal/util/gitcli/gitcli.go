// Package gitcli provides utilities to query git repo or query data via Git CLI.
package gitcli

import (
	"os"

	"github.com/go-git/go-git/v5"
)

// TryRoot finds the git repo root, or fallback to current working if fail
func TryRoot() string {
	gitRoot, err := Root()
	if err == nil {
		return gitRoot
	}
	wd, _ := os.Getwd()
	return wd
}

func Root() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	repo, err := git.PlainOpenWithOptions(wd, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return "", err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return "", err
	}

	return wt.Filesystem.Root(), nil
}
