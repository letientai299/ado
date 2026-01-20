package pipeline

import (
	"context"
	"fmt"
	"strings"

	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/rest"
	"github.com/letientai299/ado/internal/util/fp"
	"github.com/spf13/cobra"
)

// keywordProvider is implemented by config types that support keyword filtering.
type keywordProvider interface {
	Keywords() []string
}

type common[T keywordProvider] struct {
	ctx     context.Context
	cfg     *config.Config
	client  *rest.Client
	baseURL string
	repoID  string
	opts    T
}

func newCommon[T keywordProvider](cmd *cobra.Command, opts T) (*common[T], error) {
	ctx := cmd.Context()
	cfg := config.From(ctx)
	token, err := cfg.Token()
	if err != nil {
		return nil, err
	}

	client := rest.New(token)
	baseURL := fmt.Sprintf(
		"https://dev.azure.com/%s/%s/_build?definitionId=",
		cfg.Repository.Org,
		cfg.Repository.Project,
	)

	// Get repository ID for filtering pipelines
	repoInfo, err := client.Git().RepoInfo(ctx, cfg.Repository)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository info: %w", err)
	}

	return &common[T]{
		ctx:     ctx,
		cfg:     cfg,
		client:  client,
		baseURL: baseURL,
		repoID:  repoInfo.Id,
		opts:    opts,
	}, nil
}

type filterConfig struct {
	keywords []string
}

// Keywords return the filter keywords.
func (f *filterConfig) Keywords() []string {
	return f.keywords
}

// list fetches all pipelines from the API and converts them to Pipeline structs.
func (c *common[T]) list() ([]Pipeline, error) {
	all, err := c.client.Pipelines().
		Definitions(c.cfg.Repository).
		List(c.ctx, rest.ListOptions{RepositoryID: c.repoID})
	if err != nil {
		return nil, err
	}

	return fp.Map(all, func(m models.BuildDefinition) Pipeline {
		return toPipeline(m, c.baseURL)
	}), nil
}

// filter returns pipelines matching all keywords.
func (c *common[T]) filter(all []Pipeline) []Pipeline {
	keywords := c.opts.Keywords()
	if len(keywords) == 0 {
		return all
	}

	var result []Pipeline
	for _, p := range all {
		if c.containsAll(p, keywords) {
			result = append(result, p)
		}
	}
	return result
}

// containsAll checks if a pipeline's name or path contains all keywords.
func (c *common[T]) containsAll(p Pipeline, keywords []string) bool {
	name := strings.ToLower(p.Name)
	path := strings.ToLower(p.Path)
	for _, pattern := range keywords {
		pat := strings.ToLower(pattern)
		if !strings.Contains(name, pat) && !strings.Contains(path, pat) {
			return false
		}
	}
	return true
}
