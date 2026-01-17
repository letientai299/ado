// Package gitcli provides utilities to query git repo or query data via Git CLI.
package gitcli

import (
	"os"

	"github.com/go-git/go-git/v5"
)

// Root finds the git repo root, or fallback to current working if fail
func Root() string {
	wd, _ := os.Getwd()
	repo, err := Open()
	if err != nil {
		return wd
	}

	wt, err := repo.Worktree()
	if err != nil {
		return wd
	}

	return wt.Filesystem.Root()
}

func Open() (*git.Repository, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return git.PlainOpenWithOptions(wd, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
}

// RemoteURL returns the first URL of the specified remote.
func RemoteURL(remoteName string) (string, error) {
	repo, err := Open()
	if err != nil {
		return "", err
	}

	remote, err := repo.Remote(remoteName)
	if err != nil {
		return "", err
	}

	if len(remote.Config().URLs) == 0 {
		return "", git.ErrRemoteNotFound
	}

	return remote.Config().URLs[0], nil
}
