package _shared

import (
	"bytes"
	"testing"
)

func TestBool_AppendTo(t *testing.T) {
	tests := []struct {
		name     string
		b        Bool
		expected string
	}{
		{
			name:     "empty string does not append",
			b:        Bool(""),
			expected: "",
		},
		{
			name:     "non-empty string appends true",
			b:        Bool("includeCommits"),
			expected: "&includeCommits=true",
		},
		{
			name:     "another key",
			b:        Bool("searchCriteria.includeLinks"),
			expected: "&searchCriteria.includeLinks=true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			tt.b.AppendTo(buf)
			if buf.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, buf.String())
			}
		})
	}
}

func TestQueriers_AppendTo(t *testing.T) {
	tests := []struct {
		name     string
		qs       Queriers
		expected string
	}{
		{
			name:     "empty slice",
			qs:       Queriers{},
			expected: "",
		},
		{
			name:     "nil slice",
			qs:       nil,
			expected: "",
		},
		{
			name:     "single querier",
			qs:       Queriers{Bool("enabled")},
			expected: "&enabled=true",
		},
		{
			name: "multiple queriers",
			qs: Queriers{
				Bool("enabled"),
				KV[string]{Key: "name", Value: "test"},
			},
			expected: "&enabled=true&name=test",
		},
		{
			name: "nested queriers",
			qs: Queriers{
				Bool("first"),
				Queriers{
					Bool("nested"),
					KV[int]{Key: "count", Value: 5},
				},
				Bool("last"),
			},
			expected: "&first=true&nested=true&count=5&last=true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			tt.qs.AppendTo(buf)
			if buf.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, buf.String())
			}
		})
	}
}

func TestKV_AppendTo(t *testing.T) {
	tests := []struct {
		name     string
		kv       Querier
		expected string
	}{
		{
			name:     "string value",
			kv:       KV[string]{Key: "name", Value: "test"},
			expected: "&name=test",
		},
		{
			name:     "empty string value",
			kv:       KV[string]{Key: "name", Value: ""},
			expected: "&name=",
		},
		{
			name:     "int value",
			kv:       KV[int]{Key: "count", Value: 42},
			expected: "&count=42",
		},
		{
			name:     "zero int value",
			kv:       KV[int]{Key: "count", Value: 0},
			expected: "&count=0",
		},
		{
			name:     "negative int value",
			kv:       KV[int]{Key: "offset", Value: -10},
			expected: "&offset=-10",
		},
		{
			name:     "bool true value",
			kv:       KV[bool]{Key: "active", Value: true},
			expected: "&active=true",
		},
		{
			name:     "bool false value",
			kv:       KV[bool]{Key: "active", Value: false},
			expected: "&active=false",
		},
		{
			name:     "float value",
			kv:       KV[float64]{Key: "ratio", Value: 3.14},
			expected: "&ratio=3.14",
		},
		{
			name:     "empty key",
			kv:       KV[string]{Key: "", Value: "value"},
			expected: "&=value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			tt.kv.AppendTo(buf)
			if buf.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, buf.String())
			}
		})
	}
}

func TestMap_AppendTo(t *testing.T) {
	tests := []struct {
		name     string
		m        Map
		expected []string // multiple expected values due to map iteration order
	}{
		{
			name:     "empty map",
			m:        Map{},
			expected: []string{""},
		},
		{
			name:     "nil map",
			m:        nil,
			expected: []string{""},
		},
		{
			name:     "single entry",
			m:        Map{"key": "value"},
			expected: []string{"&key=value"},
		},
		{
			name:     "int value",
			m:        Map{"count": 42},
			expected: []string{"&count=42"},
		},
		{
			name: "multiple entries",
			m:    Map{"a": 1, "b": 2},
			expected: []string{
				"&a=1&b=2",
				"&b=2&a=1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			tt.m.AppendTo(buf)
			result := buf.String()
			found := false
			for _, exp := range tt.expected {
				if result == exp {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected one of %v, got %q", tt.expected, result)
			}
		})
	}
}

func TestQuerier_InterfaceCompliance(t *testing.T) {
	// These tests verify that all types implement the Querier interface
	// at runtime as well as the compile-time checks in querier.go
	var _ Querier = Bool("")
	var _ Querier = Queriers{}
	var _ Querier = KV[string]{}
	var _ Querier = KV[int]{}
	var _ Querier = KV[bool]{}
	var _ Querier = Map{}
}

func TestAppendQueries(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		queries  []Querier
		expected string
	}{
		{
			name:     "no queries",
			url:      "https://dev.azure.com/org/project/_apis/git/repositories",
			queries:  nil,
			expected: "https://dev.azure.com/org/project/_apis/git/repositories",
		},
		{
			name:     "empty queries slice",
			url:      "https://dev.azure.com/org/project/_apis/git/repositories",
			queries:  []Querier{},
			expected: "https://dev.azure.com/org/project/_apis/git/repositories",
		},
		{
			name:     "single query",
			url:      "https://dev.azure.com/org/project/_apis/git/repositories",
			queries:  []Querier{KV[string]{Key: "api-version", Value: "7.0"}},
			expected: "https://dev.azure.com/org/project/_apis/git/repositories?api-version=7.0",
		},
		{
			name: "multiple queries",
			url:  "https://dev.azure.com/org/project/_apis/git/repositories",
			queries: []Querier{
				KV[string]{Key: "api-version", Value: "7.0"},
				Bool("includeLinks"),
			},
			expected: "https://dev.azure.com/org/project/_apis/git/repositories?api-version=7.0&includeLinks=true",
		},
		{
			name: "with int value",
			url:  "https://dev.azure.com/org/project/_apis/build/builds",
			queries: []Querier{
				KV[string]{Key: "api-version", Value: "7.0"},
				KV[int]{Key: "$top", Value: 100},
			},
			expected: "https://dev.azure.com/org/project/_apis/build/builds?api-version=7.0&$top=100",
		},
		{
			name: "empty bool does not add parameter",
			url:  "https://example.com/api",
			queries: []Querier{
				KV[string]{Key: "key", Value: "value"},
				Bool(""),
			},
			expected: "https://example.com/api?key=value",
		},
		{
			name:     "empty url with queries",
			url:      "",
			queries:  []Querier{KV[string]{Key: "key", Value: "value"}},
			expected: "?key=value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AppendQueries(tt.url, tt.queries...)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func BenchmarkAppendQueries(b *testing.B) {
	url := "https://dev.azure.com/org/project/_apis/git/repositories"
	queries := []Querier{
		KV[string]{Key: "api-version", Value: "7.0"},
		Bool("includeLinks"),
		KV[int]{Key: "$top", Value: 100},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AppendQueries(url, queries...)
	}
}

func BenchmarkAppendQueries_NoQueries(b *testing.B) {
	url := "https://dev.azure.com/org/project/_apis/git/repositories"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AppendQueries(url)
	}
}

func BenchmarkAppendQueries_SingleQuery(b *testing.B) {
	url := "https://dev.azure.com/org/project/_apis/git/repositories"
	query := KV[string]{Key: "api-version", Value: "7.0"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AppendQueries(url, query)
	}
}

func BenchmarkBool_AppendTo(b *testing.B) {
	buf := &bytes.Buffer{}
	bq := Bool("includeLinks")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		bq.AppendTo(buf)
	}
}

func BenchmarkKV_AppendTo_String(b *testing.B) {
	buf := &bytes.Buffer{}
	kv := KV[string]{Key: "api-version", Value: "7.0"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		kv.AppendTo(buf)
	}
}

func BenchmarkKV_AppendTo_Int(b *testing.B) {
	buf := &bytes.Buffer{}
	kv := KV[int]{Key: "$top", Value: 100}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		kv.AppendTo(buf)
	}
}
