package rest

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/goccy/go-json"
	"github.com/letientai299/ado/internal/rest/_shared"
)

func httpGet[T any](
	ctx context.Context,
	c Client,
	url string,
	queries ..._shared.Querier,
) (*T, error) {
	var sb strings.Builder
	sb.WriteString(url)
	sb.WriteString("?api-version=")
	sb.WriteString(apiVersion)
	_shared.Queriers(queries).AppendTo(&sb)
	url = sb.String()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("fail to create HTTP request: %v", err)
		return nil, err
	}

	return call[T](c, req)
}

func call[T any](c Client, req *http.Request) (*T, error) {
	log.Debugf("HTTP request: %s %s", req.Method, req.URL)

	req.Header.Set("Authorization", "Bearer "+c.token)
	resp, err := c.http.Do(req)
	if err != nil {
		log.Error("fail to call HTTP request", "url", req.URL, "err", err)
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if err = validateResponse(resp); err != nil {
		return nil, err
	}

	return decode[T](resp.Body)
}

func validateResponse(resp *http.Response) error {
	code := resp.StatusCode
	uri := resp.Request.URL.RequestURI()
	if code == 404 {
		return ErrNotFound.WithURI(uri)
	}

	if code == 203 {
		return ErrAuth.WithURI(uri)
	}

	if code == 403 {
		return ErrForbidden.WithURI(uri)
	}

	if code >= 500 {
		return ErrInternalServer.WithCode(code).WithURI(uri)
	}

	if code >= 400 {
		return ErrInvalidRequest.WithCode(code).WithURI(uri)
	}

	return nil
}

func decode[T any](body io.ReadCloser) (*T, error) {
	t := new(T)

	err := json.NewDecoder(body).Decode(t)
	if err == nil {
		return t, nil
	}

	log.Error("fail to decode response body", "target_type", reflect.TypeFor[T](), "err", err)
	return nil, err
}
