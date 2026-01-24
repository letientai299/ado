package gitcli

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/letientai299/ado/internal/util"
	"github.com/letientai299/ado/internal/util/sh"
)

var (
	authOnce sync.Once
	auth     transport.AuthMethod
	tokenFn  func() (string, error)
)

func SetTokenProvider(fn func() (string, error)) {
	tokenFn = fn
}

func getAuth() transport.AuthMethod {
	authOnce.Do(func() {
		token, err := tokenFn()
		if err != nil {
			panic(fmt.Errorf("fail to get token: %w", err))
		}
		auth = &http.BasicAuth{
			Password: token, // no need username
		}
	})
	return auth
}

// SyncToRemote ensures the local branch is pushed to remote.
func SyncToRemote(branch string, confirmFn func(ask string) bool) error {
	if !RemoteBranchExists(branch) {
		ask := fmt.Sprintf("Remote branch %s does not exist. Push it?", branch)
		if !confirmFn(ask) {
			return util.StrErr("remote branch does not exist")
		}

		if err := Push(branch); err != nil {
			return fmt.Errorf("fail to push branch: %w", err)
		}

		return nil
	}

	// Check if the local branch is ahead of remote
	remote := Origin + "/" + branch
	div, err := CompareRevision(remote, branch)
	if err != nil {
		return err
	}

	if div.NoChanges() {
		return nil
	}

	if div.IsAhead() {
		return confirmAndPush(div, confirmFn)
	}

	if div.IsBehind() {
		return confirmAndPull(div, confirmFn)
	}

	if div.IsDiverged() {
		return fmt.Errorf(
			"local and remote branches have diverged, %d ahead and %d behind",
			len(div.Ahead),
			len(div.Behind),
		)
	}

	return nil
}

func confirmAndPull(div Divergence, confirmFn func(ask string) bool) error {
	ask := fmt.Sprintf("Local branch is %d commit(s) behind remote. Pull?", len(div.Behind))
	if !confirmFn(ask) {
		return nil
	}

	if err := Pull(div.Source); err != nil {
		return fmt.Errorf("fail to pull branch: %w", err)
	}

	return nil
}

func confirmAndPush(div Divergence, confirmFn func(ask string) bool) error {
	ask := fmt.Sprintf("Local branch is %d commit(s) ahead of remote. Push?", len(div.Ahead))
	if !confirmFn(ask) {
		return nil
	}

	if err := Push(div.Source); err != nil {
		return fmt.Errorf("fail to push branch: %w", err)
	}

	return nil
}

// Push pushes the branch to the specified remote and sets upstream tracking.
func Push(branch string) error {
	repo, err := Open()
	if err != nil {
		return err
	}

	refSpec := config.RefSpec(fmt.Sprintf("refs/heads/%s:refs/heads/%s", branch, branch))
	if err = repo.Push(&git.PushOptions{
		RemoteName: Origin,
		RefSpecs:   []config.RefSpec{refSpec},
		Auth:       getAuth(),
	}); err != nil {
		return err
	}

	// Set upstream tracking (-u flag equivalent)
	cfg, err := repo.Config()
	if err != nil {
		return err
	}

	cfg.Branches[branch] = &config.Branch{
		Name:   branch,
		Remote: Origin,
		Merge:  plumbing.ReferenceName("refs/heads/" + branch),
	}
	return repo.SetConfig(cfg)
}

func Pull(branch string) error {
	repo, err := Open()
	if err != nil {
		return err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = wt.Pull(&git.PullOptions{
		RemoteName:    Origin,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		Auth:          getAuth(),
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}

	return nil
}

// FetchBranch fetches the specified branch from remote.
func FetchBranch(branch string) error {
	repo, err := Open()
	if err != nil {
		return err
	}

	refSpec := config.RefSpec(fmt.Sprintf(
		"refs/heads/%[1]s:refs/remotes/%[2]s/%[1]s",
		branch, Origin,
	))

	err = repo.Fetch(&git.FetchOptions{
		RemoteName: Origin,
		RefSpecs:   []config.RefSpec{refSpec},
		Auth:       getAuth(),
	})

	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}

	return nil
}

// Rebase rebases the current branch onto the target branch.
// Returns ErrRebaseConflict if conflicts occur.
func Rebase(target string) error {
	_, err := runGit("rebase", Origin+"/"+target)
	if err != nil {
		// Abort the rebase to leave the repo in a clean state
		_, _ = runGit("rebase", "--abort")
		return ErrRebaseConflict
	}
	return nil
}

// runGit executes a git command and returns its output.
func runGit(args ...string) (string, error) {
	root := Root()
	cmd := fmt.Sprintf("cd %q && git %s", root, strings.Join(args, " "))
	return sh.Run(cmd)
}
