package rest

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/goccy/go-json"
	"github.com/letientai299/ado/internal/rest/_shared"
	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/util"
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

	body, err := call(c, req)
	if err != nil {
		return nil, err
	}

	return decode[T](body)
}

func call(c Client, req *http.Request) (io.ReadCloser, error) {
	log.Debugf("HTTP request: %s %s", req.Method, req.URL)

	req.Header.Set("Authorization", "Bearer "+c.token)
	resp, err := c.http.Do(req)
	if err != nil {
		log.Error("fail to call HTTP request", "url", req.URL, "err", err)
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

func decode[T any](body io.ReadCloser) (*T, error) {
	t := new(T)
	bs, err := io.ReadAll(body)
	if err != nil {
		log.Errorf("fail to read response body: %v", err)
		return nil, err
	}

	if err = json.NewDecoder(bytes.NewReader(bs)).Decode(t); err != nil {
		log.Error(
			"fail to decode response body",
			"target_type",
			reflect.TypeFor[T](),
			"err",
			err,
			"body",
			string(bs),
		)
		return nil, err
	}

	return t, nil
}

func extractError(code int, body io.ReadCloser) error {
	msg, err := io.ReadAll(body)
	var x any
	_ = json.Unmarshal(msg, &x)
	s := styles.YAML(x)
	if code >= 500 {
		return errors.Join(ErrInternalServer, err, util.StrErr(s))
	}

	return errors.Join(ErrInvalidRequest, err, util.StrErr(s))
}
