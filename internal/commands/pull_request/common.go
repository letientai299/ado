package pull_request

import (
	"context"
	"net/url"

	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/rest"
	"github.com/spf13/cobra"
)

type common[T any] struct {
	ctx     context.Context
	cfg     *config.Config
	client  *rest.Client
	baseURL string
	opts    T
}

func newCommon[T any](cmd *cobra.Command, opts T) (*common[T], error) {
	ctx := cmd.Context()
	cfg := config.From(ctx)
	token, err := cfg.Token()
	if err != nil {
		return nil, err
	}

	client := rest.New(token)
	baseURL, _ := url.JoinPath(cfg.Repository.WebURL(), "pullRequest")
	return &common[T]{
		ctx:     ctx,
		cfg:     cfg,
		client:  client,
		baseURL: baseURL,
		opts:    opts,
	}, nil
}

func copyCommon[A, B any](a *common[A], mod func(*common[B]) *common[B]) *common[B] {
	b := &common[B]{
		ctx:     a.ctx,
		cfg:     a.cfg,
		client:  a.client,
		baseURL: a.baseURL,
	}
	return mod(b)
}

type filterConfig struct {
	mine     bool // find my PRs
	draft    bool // find draft PRs
	keywords []string
}

func (f *filterConfig) RegisterFlags(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()
	flags.BoolVarP(&f.mine, "mine", "m", false, "show only your PRs")
	flags.BoolVarP(&f.draft, "draft", "d", false, "include draft PRs")
}
