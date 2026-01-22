package rest

import (
	"context"
	"net/http"
	"net/url"

	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/models"
)

const (
	// apiVersion is the Azure DevOps REST API version used by this client.
	apiVersion = "7.1"

	// adoHost is the base URL for Azure DevOps REST APIs.
	adoHost = "https://dev.azure.com"
)

// New creates a new Azure DevOps REST API client.
// The token should be a Personal Access Token (PAT) or Azure AD token
// with appropriate permissions for the APIs being called.
func New(token string) *Client {
	return &Client{
		token: token,
		http:  http.DefaultClient,
	}
}

// Client is the main entry point for Azure DevOps REST API operations.
// It manages authentication and provides access to various API categories.
type Client struct {
	token string
	http  *http.Client
}

// Git returns a client for Git repository and pull request operations.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/git
func (c Client) Git() Git {
	return Git{client: c}
}

// Pipelines returns a client for pipeline operations.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/pipelines
func (c Client) Pipelines() Pipelines {
	return Pipelines{client: c}
}

// Builds returns a client for build operations.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/build
func (c Client) Builds() Builds {
	return Builds{client: c}
}

// Policy returns a client for branch policy operations.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/policy
func (c Client) Policy() Policy {
	return Policy{client: c}
}

// Identity retrieves the currently authenticated user's identity information.
// This is useful for operations that require the current user's ID,
// such as setting PR votes.
//
// See:
// https://learn.microsoft.com/en-us/rest/api/azure/devops/account/accounts/list
func (c Client) Identity(ctx context.Context, org string) (*models.Identity, error) {
	api, err := url.JoinPath(adoHost, org, "_apis/connectionData")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, api, nil)
	if err != nil {
		log.Errorf("fail to create HTTP request: %v", err)
		return nil, err
	}

	type Temp struct {
		AuthenticatedUser *models.Identity `json:"authenticatedUser"`
	}

	// manually call instead of httpGet, to not include apiVersion automatically,
	// because this API doesn't share the same version with others.
	// this is not an ADO-specific API.
	t, err := call[Temp](c, req)
	if err != nil {
		return nil, err
	}

	return t.AuthenticatedUser, nil
}
