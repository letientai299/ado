package rest

import (
	"context"
	"net/http"
	"net/url"

	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/util"
)

const (
	apiVersion = "7.1"
	adoHost    = "https://dev.azure.com"
)

const (
	ErrNotFound       util.StrErr = "not found"
	ErrInternalServer util.StrErr = "ADO server internal error"
	ErrInvalidRequest util.StrErr = "invalid request"
)

func New(token string) *Client {
	return &Client{
		token: token,
		http:  http.DefaultClient,
	}
}

type Client struct {
	token string
	http  *http.Client
}

func (c Client) Git() Git {
	return Git{client: c}
}

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

	body, err := call(c, req)
	if err != nil {
		return nil, err
	}

	type Temp struct {
		AuthenticatedUser *models.Identity `json:"authenticatedUser"`
	}

	t, err := decode[Temp](body)
	if err != nil {
		return nil, err
	}

	return t.AuthenticatedUser, nil
}
