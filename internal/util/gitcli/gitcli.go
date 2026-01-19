package gitcli

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/letientai299/ado/internal/util"
)

// Root finds the git repo root or fallback to current working if fail
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

const ErrNotOnBranch = util.StrErr("not on a branch")

// CurrentBranch returns the name of the current branch.
func CurrentBranch() (string, error) {
	repo, err := Open()
	if err != nil {
		return "", err
	}

	head, err := repo.Head()
	if err != nil {
		return "", err
	}

	if !head.Name().IsBranch() {
		return "", ErrNotOnBranch
	}

	return head.Name().Short(), nil
}

type Commit struct {
	Subject string
	Body    string
}

// CommitsAhead returns the commits between target and source branch.
func CommitsAhead(target, source string) ([]Commit, error) {
	repo, err := Open()
	if err != nil {
		return nil, err
	}

	targetHash, err := repo.ResolveRevision(plumbing.Revision(target))
	if err != nil {
		return nil, fmt.Errorf("target branch %s not found: %w", target, err)
	}

	sourceHash, err := repo.ResolveRevision(plumbing.Revision(source))
	if err != nil {
		return nil, fmt.Errorf("source branch %s not found: %w", source, err)
	}

	sourceCommit, err := repo.CommitObject(*sourceHash)
	if err != nil {
		return nil, err
	}

	targetCommit, err := repo.CommitObject(*targetHash)
	if err != nil {
		return nil, err
	}

	bases, err := sourceCommit.MergeBase(targetCommit)
	if err != nil {
		return nil, err
	}

	stopAt := make(map[plumbing.Hash]struct{}, len(bases)+1)
	for _, b := range bases {
		stopAt[b.Hash] = struct{}{}
	}
	stopAt[*targetHash] = struct{}{}

	iter, err := repo.Log(&git.LogOptions{From: *sourceHash})
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var commits []Commit
	err = iter.ForEach(func(c *object.Commit) error {
		if _, stop := stopAt[c.Hash]; stop {
			return storer.ErrStop
		}

		subject, body, _ := strings.Cut(c.Message, "\n")
		commits = append(commits, Commit{
			Subject: strings.TrimSpace(subject),
			Body:    strings.TrimSpace(body),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return commits, nil
}

// Push pushes the branch to the specified remote and sets upstream tracking.
func Push(remoteName, branch string) error {
	repo, err := Open()
	if err != nil {
		return err
	}

	refSpec := config.RefSpec(fmt.Sprintf("refs/heads/%s:refs/heads/%s", branch, branch))
	if err = repo.Push(&git.PushOptions{
		RemoteName: remoteName,
		RefSpecs:   []config.RefSpec{refSpec},
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
		Remote: remoteName,
		Merge:  plumbing.ReferenceName("refs/heads/" + branch),
	}
	return repo.SetConfig(cfg)
}

// RemoteBranchExists checks if a branch exists on the specified remote.
func RemoteBranchExists(remoteName, branch string) bool {
	repo, err := Open()
	if err != nil {
		return false
	}

	remote, err := repo.Remote(remoteName)
	if err != nil {
		return false
	}

	refs, err := remote.List(&git.ListOptions{})
	if err != nil {
		return false
	}

	branchRef := plumbing.ReferenceName("refs/heads/" + branch)
	return slices.ContainsFunc(refs, func(ref *plumbing.Reference) bool {
		return ref.Name() == branchRef
	})
}
