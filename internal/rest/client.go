package rest

import (
	"net/http"

	"github.com/charmbracelet/log"
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

func New(tenant string) *Client {
	token, err := util.GetToken(tenant)
	if err != nil {
		log.Fatal(err)
	}

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
