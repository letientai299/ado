package rest

import (
	"context"
	"fmt"
	"time"

	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/util/cache"
)

// Core provides access to Azure DevOps Core APIs (Projects, Teams, etc.).
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/core
type Core struct {
	client Client
}

// Project retrieves details for a team project by name or ID.
//
// See: https://learn.microsoft.com/en-us/rest/api/azure/devops/core/projects/get
func (c Core) Project(ctx context.Context, org, project string) (*models.TeamProject, error) {
	cacheKey := fmt.Sprintf("project_%s_%s", org, project)
	if cached, ok := cache.Get[models.TeamProject](cacheKey); ok {
		return cached, nil
	}

	url := fmt.Sprintf("%s/%s/_apis/projects/%s", adoHost, org, project)
	ctx = WithAPIVersion(ctx, apiVersion7_1)
	res, err := httpGet[models.TeamProject](ctx, c.client, url)
	if err != nil {
		return nil, err
	}

	_ = cache.Set(cacheKey, res, 24*time.Hour)
	return res, nil
}
