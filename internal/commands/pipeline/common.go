package pipeline

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/rest"
	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/ui"
	"github.com/letientai299/ado/internal/util"
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
func (c *common[T]) list() ([]models.BuildDefinition, error) {
	return c.client.Pipelines().
		Definitions(c.cfg.Repository).
		List(c.ctx, rest.ListOptions{
			RepositoryID:         c.repoID,
			IncludeAllProperties: true,
		})
}

// filter returns pipelines matching all keywords.
func (c *common[T]) filter(all []models.BuildDefinition) []models.BuildDefinition {
	keywords := c.opts.Keywords()
	var result []models.BuildDefinition
	for _, p := range all {
		if p.QueueStatus == models.DefinitionQueueStatusDisabled {
			continue
		}

		if c.containsAll(p, keywords) {
			result = append(result, p)
		}
	}
	return result
}

// containsAll checks if a pipeline's name or path contains all keywords.
func (c *common[T]) containsAll(p models.BuildDefinition, keywords []string) bool {
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

const pipelinePickTpl = `{{.Name}} {{if .Process}}({{.Process.YamlFilename | const}}){{end}}`

func pick(pipelines []models.BuildDefinition) fp.Optional[models.BuildDefinition] {
	selected := ui.Pick(pipelines, ui.PickConfig[models.BuildDefinition]{
		Render: func(w io.Writer, p models.BuildDefinition, matches []int) {
			p.Name = styles.HighlightMatch(p.Name, matches)
			util.PanicIf(styles.Render(w, pipelinePickTpl, p))
		},
		FilterValue: func(p models.BuildDefinition) string { return strings.ToLower(p.Name) },
	})
	return selected
}
