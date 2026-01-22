package styles

import (
	"io"
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name     string
		tpl      string
		data     any
		expected string
		wantErr  bool
	}{
		{
			name:     "simple string",
			tpl:      "Hello {{.}}",
			data:     "World",
			expected: "Hello World",
		},
		{
			name:     "with replaceAll",
			tpl:      `{{.Name | replaceAll "/" "-"}}`,
			data:     struct{ Name string }{Name: "feature/foo"},
			expected: "feature-foo",
		},
		{
			name:     "with trimPrefix",
			tpl:      `{{.Name | trimPrefix "refs/heads/"}}`,
			data:     struct{ Name string }{Name: "refs/heads/main"},
			expected: "main",
		},
		{
			name:     "range over slice",
			tpl:      `{{range .Items}}- {{.}}{{"\n"}}{{end}}`,
			data:     struct{ Items []string }{Items: []string{"a", "b", "c"}},
			expected: "- a\n- b\n- c\n",
		},
		{
			name:    "invalid template",
			tpl:     "{{.Invalid",
			data:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sb strings.Builder
			err := Render(&sb, tt.tpl, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && sb.String() != tt.expected {
				t.Errorf("Render() = %q, want %q", sb.String(), tt.expected)
			}
		})
	}
}

func TestRenderS(t *testing.T) {
	result, err := RenderS("Hello {{.}}", "World")
	if err != nil {
		t.Errorf("RenderS() error = %v", err)
	}
	if result != "Hello World" {
		t.Errorf("RenderS() = %q, want %q", result, "Hello World")
	}
}

func TestTemplateCache(t *testing.T) {
	ClearTemplateCache()

	tpl := "Test {{.}}"

	// First call should parse and cache
	result1, err := RenderS(tpl, "1")
	if err != nil {
		t.Fatalf("First RenderS() error = %v", err)
	}
	if result1 != "Test 1" {
		t.Errorf("First RenderS() = %q, want %q", result1, "Test 1")
	}

	// Second call should use cache
	result2, err := RenderS(tpl, "2")
	if err != nil {
		t.Fatalf("Second RenderS() error = %v", err)
	}
	if result2 != "Test 2" {
		t.Errorf("Second RenderS() = %q, want %q", result2, "Test 2")
	}
}

// BenchmarkRenderS benchmarks template rendering with caching.
func BenchmarkRenderS(b *testing.B) {
	ClearTemplateCache()

	tests := []struct {
		name string
		tpl  string
		data any
	}{
		{
			name: "simple",
			tpl:  "Hello {{.}}",
			data: "World",
		},
		{
			name: "with_funcs",
			tpl:  `{{.Name | replaceAll "/" "-" | trimSpace}}`,
			data: struct{ Name string }{Name: "feature/foo-bar"},
		},
		{
			name: "range_5_items",
			tpl:  `{{range .Items}}- {{.}}{{"\n"}}{{end}}`,
			data: struct{ Items []string }{Items: []string{"a", "b", "c", "d", "e"}},
		},
		{
			name: "pr_title_template",
			tpl:  `{{.BranchName | replaceAll "/" "-"}}`,
			data: struct{ BranchName string }{BranchName: "feature/JIRA-123/add-new-feature"},
		},
		{
			name: "pr_desc_template",
			tpl: `{{range .Commits}}- {{.Subject}}
{{end}}`,
			data: struct {
				Commits []struct{ Subject string }
			}{
				Commits: []struct{ Subject string }{
					{Subject: "feat: add foo"},
					{Subject: "fix: bar bug"},
					{Subject: "docs: update readme"},
				},
			},
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			// Warm up cache
			_, _ = RenderS(tt.tpl, tt.data)

			b.ResetTimer()
			b.ReportAllocs()
			for b.Loop() {
				_, _ = RenderS(tt.tpl, tt.data)
			}
		})
	}
}

// BenchmarkRenderSUncached measures the cost without caching (first parse).
func BenchmarkRenderSUncached(b *testing.B) {
	tpl := `{{.BranchName | replaceAll "/" "-"}}`
	data := struct{ BranchName string }{BranchName: "feature/foo"}

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		ClearTemplateCache()
		_, _ = RenderS(tpl, data)
	}
}

// BenchmarkRenderWithCustomFuncMap benchmarks with custom func maps (slow path).
func BenchmarkRenderWithCustomFuncMap(b *testing.B) {
	tpl := "Hello {{upper .}}"
	data := "world"
	customFuncs := map[string]any{
		"upper": strings.ToUpper,
	}

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		var sb strings.Builder
		_ = Render(&sb, tpl, data, customFuncs)
	}
}

// BenchmarkHighlightMatch benchmarks the highlight matching function.
func BenchmarkHighlightMatch(b *testing.B) {
	tests := []struct {
		name    string
		s       string
		matches []int
	}{
		{"no_matches", "hello world", nil},
		{"few_matches", "hello world", []int{0, 6}},
		{"many_matches", "hello world", []int{0, 1, 2, 3, 4, 6, 7, 8, 9, 10}},
		{"long_string", strings.Repeat("hello ", 100), []int{0, 5, 10, 15, 20}},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = HighlightMatch(tt.s, tt.matches)
			}
		})
	}
}

// BenchmarkRenderOut benchmarks output rendering (to discard).
func BenchmarkRenderOut(b *testing.B) {
	tpl := `{{range .Items}}- {{.}}{{"\n"}}{{end}}`
	data := struct{ Items []string }{Items: []string{"a", "b", "c", "d", "e"}}

	// Warm up cache
	_ = Render(io.Discard, tpl, data)

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = Render(io.Discard, tpl, data)
	}
}
