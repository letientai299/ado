package gitcli

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/letientai299/ado/internal/util"
)

const Origin = "origin"

// cachedRepo holds the cached repository and the working directory it was opened from.
var (
	repoMu     sync.Mutex
	cachedRepo *git.Repository
	cachedWd   string
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

// Open returns a cached git repository handle for the current working directory.
// The repository is cached and reused for subsequent calls from the same directory.
func Open() (*git.Repository, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	repoMu.Lock()
	defer repoMu.Unlock()

	// Return cached repo if working directory hasn't changed
	if cachedRepo != nil && cachedWd == wd {
		return cachedRepo, nil
	}

	repo, err := git.PlainOpenWithOptions(wd, &git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: true,
	})
	if err != nil {
		return nil, err
	}

	cachedRepo = repo
	cachedWd = wd
	return repo, nil
}

// ClearCache clears the cached repository. Useful for testing.
func ClearCache() {
	repoMu.Lock()
	defer repoMu.Unlock()
	cachedRepo = nil
	cachedWd = ""
}

// RemoteURL returns the first URL of the specified remote.
func RemoteURL() (string, error) {
	repo, err := Open()
	if err != nil {
		return "", err
	}

	remote, err := repo.Remote(Origin)
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

type Divergence struct {
	Target string
	Source string
	Ahead  []Commit
	Behind []Commit
}

func (d Divergence) NoChanges() bool {
	return len(d.Ahead) == 0 && len(d.Behind) == 0
}

func (d Divergence) IsBehind() bool {
	return len(d.Behind) > 0 && len(d.Ahead) == 0
}

func (d Divergence) IsAhead() bool {
	return len(d.Ahead) > 0 && len(d.Behind) == 0
}

func (d Divergence) IsDiverged() bool {
	return len(d.Ahead) > 0 && len(d.Behind) > 0
}

// CompareRevision returns the divergence between target and source branch.
func CompareRevision(target, source string) (Divergence, error) {
	repo, err := Open()
	if err != nil {
		return Divergence{}, err
	}

	targetHash, err := repo.ResolveRevision(plumbing.Revision(target))
	if err != nil {
		return Divergence{}, fmt.Errorf("target branch %s not found: %w", target, err)
	}

	sourceHash, err := repo.ResolveRevision(plumbing.Revision(source))
	if err != nil {
		return Divergence{}, fmt.Errorf("source branch %s not found: %w", source, err)
	}

	sourceCommit, err := repo.CommitObject(*sourceHash)
	if err != nil {
		return Divergence{}, err
	}

	targetCommit, err := repo.CommitObject(*targetHash)
	if err != nil {
		return Divergence{}, err
	}

	bases, err := sourceCommit.MergeBase(targetCommit)
	if err != nil {
		return Divergence{}, err
	}

	baseHashes := make(map[plumbing.Hash]struct{}, len(bases))
	for _, b := range bases {
		baseHashes[b.Hash] = struct{}{}
	}

	ahead, err := collectCommits(repo, *sourceHash, baseHashes)
	if err != nil {
		return Divergence{}, err
	}

	behind, err := collectCommits(repo, *targetHash, baseHashes)
	if err != nil {
		return Divergence{}, err
	}

	return Divergence{Target: target, Source: source, Ahead: ahead, Behind: behind}, nil
}

func collectCommits(
	repo *git.Repository,
	from plumbing.Hash,
	stopAt map[plumbing.Hash]struct{},
) ([]Commit, error) {
	iter, err := repo.Log(&git.LogOptions{From: from})
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

	return commits, err
}

// RemoteBranchExists checks if a branch exists on the specified remote.
func RemoteBranchExists(branch string) bool {
	repo, err := Open()
	if err != nil {
		return false
	}

	remote, err := repo.Remote(Origin)
	if err != nil {
		return false
	}

	err = remote.Fetch(&git.FetchOptions{
		Auth: getAuth(),
		RefSpecs: []config.RefSpec{
			config.RefSpec(fmt.Sprintf("refs/heads/%[1]s:refs/remotes/%[2]s/%[1]s", branch, Origin)),
		},
		Depth: 1,
	})

	return err == nil || errors.Is(err, git.NoErrAlreadyUpToDate)
}
