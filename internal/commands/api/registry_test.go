package api

import "testing"

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"PRs", "prs"},
		{"ByID", "by_id"},
		{"ForProject", "for_project"},
		{"RepoInfo", "repo_info"},
		{"LogContent", "log_content"},
		{"List", "list"},
		{"ID", "id"},
		{"", ""},
		{"lowercase", "lowercase"},
		{"ALLCAPS", "allcaps"},
		{"Already_Snake", "already_snake"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toSnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("toSnakeCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRegistry_Complete(t *testing.T) {
	r := NewRegistry()

	// Manually register some endpoints for testing
	r.register("git.prs.list", &Endpoint{Path: "git.prs.list"})
	r.register("git.prs.by_id", &Endpoint{Path: "git.prs.by_id"})
	r.register("git.repo_info", &Endpoint{Path: "git.repo_info"})
	r.register("builds.for_project.list", &Endpoint{Path: "builds.for_project.list"})

	tests := []struct {
		prefix   string
		expected []string
	}{
		{"", []string{"builds", "git"}},
		{"g", []string{"git"}},
		{"git", []string{"git"}},
		{"git.", []string{"git.prs", "git.repo_info"}},
		{"git.prs.", []string{"git.prs.by_id", "git.prs.list"}},
		{"nonexistent", nil},
	}

	for _, tt := range tests {
		t.Run(tt.prefix, func(t *testing.T) {
			results := r.Complete(tt.prefix)
			if len(results) != len(tt.expected) {
				t.Errorf("Complete(%q) returned %d results, want %d: %v",
					tt.prefix, len(results), len(tt.expected), results)
				return
			}
			for i, exp := range tt.expected {
				if results[i] != exp {
					t.Errorf("Complete(%q)[%d] = %q, want %q",
						tt.prefix, i, results[i], exp)
				}
			}
		})
	}
}

func TestRegistry_Get(t *testing.T) {
	r := NewRegistry()

	endpoint := &Endpoint{Path: "test.endpoint"}
	r.register("test.endpoint", endpoint)

	// Should find a registered endpoint
	if got := r.Get("test.endpoint"); got != endpoint {
		t.Errorf("Get(test.endpoint) did not return expected endpoint")
	}

	// Should return nil for non-existent
	if got := r.Get("nonexistent"); got != nil {
		t.Errorf("Get(nonexistent) = %v, want nil", got)
	}
}

func TestRegistry_Paths(t *testing.T) {
	r := NewRegistry()

	r.register("c.endpoint", &Endpoint{})
	r.register("a.endpoint", &Endpoint{})
	r.register("b.endpoint", &Endpoint{})

	paths := r.Paths()

	expected := []string{"a.endpoint", "b.endpoint", "c.endpoint"}
	if len(paths) != len(expected) {
		t.Fatalf("Paths() returned %d paths, want %d", len(paths), len(expected))
	}

	for i, exp := range expected {
		if paths[i] != exp {
			t.Errorf("Paths()[%d] = %q, want %q", i, paths[i], exp)
		}
	}
}
