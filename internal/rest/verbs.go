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

// apiVersionQuery is the query parameter added to all API requests.
var apiVersionQuery = _shared.KV[string]{
	Key:   "api-version",
	Value: apiVersion,
}

// httpGet performs an HTTP GET request and decodes the JSON response.
// Query parameters are appended using the provided Querier implementations.
func httpGet[T any](ctx context.Context, c Client, url string, qs ..._shared.Querier) (*T, error) {
	qs = append(qs, apiVersionQuery)
	url = appendQueries(url, qs...)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("fail to create HTTP request: %v", err)
		return nil, err
	}

	return call[T](c, req)
}

// httpPost performs an HTTP POST request with a JSON body.
func httpPost[T any](ctx context.Context, c Client, url string, body any) (*T, error) {
	return httpX[T](ctx, c, http.MethodPost, url, body)
}

// httpPatch performs an HTTP PATCH request with a JSON body.
func httpPatch[T any](ctx context.Context, c Client, url string, body any) (*T, error) {
	return httpX[T](ctx, c, http.MethodPatch, url, body)
}

// httpPut performs an HTTP PUT request with a JSON body.
func httpPut[T any](ctx context.Context, c Client, url string, body any) (*T, error) {
	return httpX[T](ctx, c, http.MethodPut, url, body)
}

// httpX is a helper that performs HTTP requests with JSON body for POST, PUT, PATCH.
func httpX[T any](ctx context.Context, c Client, method, url string, body any) (*T, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	url = appendQueries(url, apiVersionQuery)
	var b io.Reader = strings.NewReader(string(jsonBody))
	req, err := http.NewRequestWithContext(ctx, method, url, b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return call[T](c, req)
}

// appendQueries appends query parameters to a URL.
// The first parameter is converted from & to ? to start the query string.
func appendQueries(url string, queries ..._shared.Querier) string {
	var sb strings.Builder
	sb.WriteString(url)
	_shared.Queriers(queries).AppendTo(&sb)
	s := sb.String()
	return strings.Replace(s, "&", "?", 1)
}

// call executes an HTTP request and decodes the JSON response.
// Handles authentication, error responses, and response decoding.
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

// logErrResponse logs details of an error response for debugging.
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

// validateResponse checks the HTTP response status and returns an appropriate error.
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

// decode unmarshals JSON from the response body into the specified type.
func decode[T any](body io.ReadCloser) (*T, error) {
	t := new(T)

	err := json.NewDecoder(body).Decode(t)
	if err == nil {
		return t, nil
	}

	log.Error("fail to decode response body", "target_type", reflect.TypeFor[T](), "err", err)
	return nil, err
}
