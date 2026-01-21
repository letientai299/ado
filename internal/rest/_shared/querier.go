package _shared

import (
	"fmt"
	"io"
)

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
		_, _ = w.Write([]byte("&" + string(b) + "=true"))
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
	_, _ = w.Write([]byte("&" + kv.Key + "="))
	_, _ = fmt.Fprint(w, kv.Value)
}

// Map is a map of query parameters.
// Each entry is written as a separate query parameter.
type Map map[string]any

// AppendTo writes all map entries as query parameters.
func (m Map) AppendTo(w io.Writer) {
	for k, v := range m {
		_, _ = w.Write([]byte("&" + k + "="))
		_, _ = fmt.Fprint(w, v)
	}
}
