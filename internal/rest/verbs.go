package rest

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/goccy/go-json"
	"github.com/letientai299/ado/internal/rest/_shared"
	"github.com/letientai299/ado/internal/styles"
)

type ctxKey string

const (
	ctxApiVersion ctxKey = apiVersion7_1
)

func WithAPIVersion(ctx context.Context, ver string) context.Context {
	return context.WithValue(ctx, ctxApiVersion, ver)
}

func ApiVersion(ctx context.Context) (string, bool) {
	v := ctx.Value(ctxApiVersion)
	ver, ok := v.(string)
	return ver, ok
}

func httpGet[T any](ctx context.Context, c Client, url string, qs ..._shared.Querier) (*T, error) {
	url = buildFullURL(ctx, url, qs)
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

func httpDelete[T any](
	ctx context.Context,
	c Client,
	url string,
	qs ..._shared.Querier,
) (*T, error) {
	url = buildFullURL(ctx, url, qs)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		log.Errorf("fail to create HTTP request: %v", err)
		return nil, err
	}
	return call[T](c, req)
}

// httpPatchJsonPatch sends a PATCH request with Content-Type: application/json-patch+json.
// Required by the Work Item Create/Update APIs.
func httpPatchJsonPatch[T any](
	ctx context.Context,
	c Client,
	url string,
	body any,
	qs ..._shared.Querier,
) (*T, error) {
	bodyReader, err := prepareBody(body)
	if err != nil {
		return nil, err
	}

	url = buildFullURL(ctx, url, qs)
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json-patch+json")
	return call[T](c, req)
}

func httpDelete[T any](
	ctx context.Context,
	c Client,
	url string,
	qs ..._shared.Querier,
) (*T, error) {
	url = buildFullURL(ctx, url, qs)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		log.Errorf("fail to create HTTP request: %v", err)
		return nil, err
	}

	return call[T](c, req)
}

func httpX[T any](
	ctx context.Context,
	c Client,
	method, url string,
	body any,
	qs ..._shared.Querier,
) (*T, error) {
	bodyReader, err := prepareBody(body)
	if err != nil {
		return nil, err
	}

	url = buildFullURL(ctx, url, qs)
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return call[T](c, req)
}

func buildFullURL(ctx context.Context, url string, qs []_shared.Querier) string {
	if ver, ok := ApiVersion(ctx); ok {
		qs = append(qs, _shared.KV[string]{Key: apiVersionQuery, Value: ver})
	} else {
		qs = append(qs, _shared.KV[string]{Key: apiVersionQuery, Value: apiVersion7_1})
	}
	return _shared.AppendQueries(url, qs...)
}

func prepareBody(body any) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return strings.NewReader(string(jsonBody)), nil
}

func call[T any](c Client, req *http.Request) (*T, error) {
	return callAndDecode(c, req, decode[T])
}

func callAndDecode[T any](
	c Client,
	req *http.Request,
	decodeFn func(reader io.Reader) (*T, error),
) (*T, error) {
	start := time.Now()
	req.Header.Set("Authorization", "Bearer "+c.token)
	resp, err := c.http.Do(req)
	if err != nil {
		log.Error("fail to call HTTP request", "url", req.URL, "err", err)
		return nil, err
	}

	log.Debugf(
		"[%.3f ms] %s %s",
		float64(time.Since(start))/float64(time.Millisecond),
		req.Method,
		req.URL.RequestURI(),
	)
	defer func() { _ = resp.Body.Close() }()

	if err = validateResponse(resp); err != nil {
		logErrResponse(resp)
		return nil, err
	}

	return decodeFn(resp.Body)
}

func logErrResponse(resp *http.Response) {
	log.Errorf("HTTP response: %s %s", resp.Status, resp.Request.URL.RequestURI())

	for k, v := range resp.Header {
		log.Errorf("Header %s: %s", k, strings.Join(v, ", "))
	}

	all, _ := io.ReadAll(resp.Body)
	contentType := resp.Header.Get("content-type")
	if !strings.HasPrefix(contentType, "application/json") {
		return // don't log, too long
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

func decode[T any](body io.Reader) (*T, error) {
	t := new(T)

	err := json.NewDecoder(body).Decode(t)
	if err == nil {
		return t, nil
	}

	log.Error("fail to decode response body", "target_type", reflect.TypeFor[T](), "err", err)
	return nil, err
}
