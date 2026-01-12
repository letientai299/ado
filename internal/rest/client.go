package rest

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/goccy/go-json"
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

func escape(v any) string {
	return url.QueryEscape(fmt.Sprint(v))
}

func (c Client) httpGet(ctx context.Context, url string, keyVals ...any) (io.ReadCloser, error) {
	fullUrl := getFullURL(url, keyVals)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullUrl, nil)
	if err != nil {
		log.Error("fail to create HTTP request", "url", url, "err", err)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		log.Error("fail to call HTTP request", "url", url, "err", err)
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, errors.Join(ErrNotFound, resp.Body.Close())
	}

	if resp.StatusCode >= 400 {
		return nil, extractError(resp.StatusCode, resp.Body)
	}

	return resp.Body, nil
}

func getFullURL(url string, keyVals []any) string {
	var sb strings.Builder
	sb.WriteString(url)
	sb.WriteString("?api-version=")
	sb.WriteString(apiVersion)

	for i := 0; i < len(keyVals); i++ {
		sb.WriteByte('&')
		sb.WriteString(escape(keyVals[i]))
		if i+1 >= len(keyVals) {
			panic("invalid keyVals length")
		}
		sb.WriteByte('=')
		sb.WriteString(escape(keyVals[i+1]))
		i++
	}

	fullUrl := sb.String()
	return fullUrl
}

func extractError(code int, body io.ReadCloser) error {
	msg, err := io.ReadAll(body)
	if code >= 500 {
		return errors.Join(ErrInternalServer, err, util.StrErr(msg))
	}

	return errors.Join(ErrInvalidRequest, err, util.StrErr(msg))
}

func (c Client) Git() Git {
	return Git{client: c}
}

func httpGet[T any](ctx context.Context, c Client, url string, keyVals ...any) (*T, error) {
	rc, err := c.httpGet(ctx, url, keyVals...)
	if err != nil {
		return nil, err
	}

	t := new(T)
	if err = json.NewDecoder(rc).Decode(t); err != nil {
		log.Error("fail to decode response body", "target_type", reflect.TypeFor[T](), "err", err)
		return nil, err
	}

	return t, nil
}
