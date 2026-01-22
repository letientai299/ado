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
	"github.com/letientai299/ado/internal/styles"
)

var apiVersionQuery = _shared.KV[string]{
	Key:   "api-version",
	Value: apiVersion,
}

func httpGet[T any](ctx context.Context, c Client, url string, qs ..._shared.Querier) (*T, error) {
	qs = append(qs, apiVersionQuery)
	url = _shared.AppendQueries(url, qs...)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("fail to create HTTP request: %v", err)
		return nil, err
	}

	return call[T](c, req)
}

func httpPost[T any](
	ctx context.Context,
	c Client,
	url string,
	body any,
	qs ..._shared.Querier,
) (*T, error) {
	return httpX[T](ctx, c, http.MethodPost, url, body, qs...)
}

func httpPatch[T any](
	ctx context.Context,
	c Client,
	url string,
	body any,
	qs ..._shared.Querier,
) (*T, error) {
	return httpX[T](ctx, c, http.MethodPatch, url, body, qs...)
}

func httpPut[T any](
	ctx context.Context,
	c Client,
	url string,
	body any,
	qs ..._shared.Querier,
) (*T, error) {
	return httpX[T](ctx, c, http.MethodPut, url, body, qs...)
}

func httpX[T any](
	ctx context.Context,
	c Client,
	method, url string,
	body any,
	qs ..._shared.Querier,
) (*T, error) {
	var b io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		b = strings.NewReader(string(jsonBody))
	}

	qs = append(qs, apiVersionQuery)
	url = _shared.AppendQueries(url, qs...)
	req, err := http.NewRequestWithContext(ctx, method, url, b)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
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
		logErrResponse(resp)
		return nil, err
	}

	return decode[T](resp.Body)
}

func logErrResponse(resp *http.Response) {
	all, _ := io.ReadAll(resp.Body)
	log.Errorf("HTTP response: %s %s", resp.Status, resp.Request.URL.RequestURI())

	contentType := resp.Header.Get("content-type")
	if !strings.HasPrefix(contentType, "application/json") {
		log.Error("response body", "body", string(all))
		return
	}

	var m map[string]any
	if err := json.Unmarshal(all, &m); err != nil {
		log.Warn("failed to unmarshal JSON response", "err", err, "body", string(all))
		return
	}

	_ = styles.DumpYAML(m)
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
