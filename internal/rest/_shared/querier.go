package _shared

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// AppendQueries appends query parameters to a URL.
// The first parameter uses ? to start the query string, subsequent use &.
// If the URL already has query parameters (contains ?), all new params use &.
func AppendQueries(url string, queries ...Querier) string {
	if len(queries) == 0 {
		return url
	}

	var sb strings.Builder
	sb.Grow(len(url) + len(queries)*32) // estimate ~32 bytes per query param
	sb.WriteString(url)

	// Check if URL already has query params
	hasQuery := strings.Contains(url, "?")
	w := &queryWriter{sb: &sb, first: !hasQuery}
	for _, q := range queries {
		q.AppendTo(w)
	}
	return sb.String()
}

// queryWriter wraps strings.Builder to handle ?/& prefix.
type queryWriter struct {
	sb    *strings.Builder
	first bool
}

func (w *queryWriter) Write(p []byte) (int, error) {
	if w.first && len(p) > 0 && p[0] == '&' {
		w.sb.WriteByte('?')
		w.sb.Write(p[1:])
		w.first = false
		return len(p), nil
	}
	return w.sb.Write(p)
}

func (w *queryWriter) WriteString(s string) (int, error) {
	if w.first && len(s) > 0 && s[0] == '&' {
		w.sb.WriteByte('?')
		w.sb.WriteString(s[1:])
		w.first = false
		return len(s), nil
	}
	return w.sb.WriteString(s)
}

// Querier is an interface for types that can append query parameters to a URL.
// Implementations write query parameters in the format "&key=value" to the writer.
type Querier interface {
	AppendTo(io.Writer)
}

// Compile-time interface compliance checks.
var (
	_ Querier = Bool("")
	_ Querier = Queriers{}
	_ Querier = KV[any]{}
	_ Querier = Map{}
)

// Bool is a boolean query parameter that is set to "true" when the key is non-empty.
// Example: Bool("includeCommits") appends "&includeCommits=true".
type Bool string

// AppendTo writes the boolean query parameter if the key is non-empty.
func (b Bool) AppendTo(w io.Writer) {
	if b != "" {
		_, _ = io.WriteString(w, "&")
		_, _ = io.WriteString(w, string(b))
		_, _ = io.WriteString(w, "=true")
	}
}

// Queriers is a slice of Querier implementations that writes all contained parameters.
type Queriers []Querier

// AppendTo writes all contained query parameters.
func (qs Queriers) AppendTo(w io.Writer) {
	for _, q := range qs {
		q.AppendTo(w)
	}
}

// KV is a key-value query parameter.
// The value is converted to string using fmt.Fprint.
type KV[T any] struct {
	Key   string
	Value T
}

// AppendTo writes the key-value pair as a query parameter.
func (kv KV[T]) AppendTo(w io.Writer) {
	_, _ = io.WriteString(w, "&")
	_, _ = io.WriteString(w, kv.Key)
	_, _ = io.WriteString(w, "=")
	writeValue(w, kv.Value)
}

// writeValue writes a value to w using fast paths for common types.
func writeValue(w io.Writer, v any) {
	switch val := v.(type) {
	case string:
		_, _ = io.WriteString(w, val)
	case int:
		_, _ = io.WriteString(w, strconv.Itoa(val))
	case int64:
		_, _ = io.WriteString(w, strconv.FormatInt(val, 10))
	case int32:
		_, _ = io.WriteString(w, strconv.FormatInt(int64(val), 10))
	case bool:
		_, _ = io.WriteString(w, strconv.FormatBool(val))
	default:
		_, _ = fmt.Fprint(w, val)
	}
}

// Map is a map of query parameters.
// Each entry is written as a separate query parameter.
type Map map[string]any

// AppendTo writes all map entries as query parameters.
func (m Map) AppendTo(w io.Writer) {
	for k, v := range m {
		_, _ = io.WriteString(w, "&")
		_, _ = io.WriteString(w, k)
		_, _ = io.WriteString(w, "=")
		writeValue(w, v)
	}
}
