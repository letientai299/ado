package gitcli

import (
	"path/filepath"
	"testing"
)

func TestResolveRepoRelativePath(t *testing.T) {
	root := Root()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "absolute path inside repo",
			in:   filepath.Join(root, "internal", "util", "util.go"),
			want: "internal/util/util.go",
		},
		{
			name: "path outside repo returns input as-is",
			in:   filepath.Join(root, "..", "outside_file.txt"),
			want: filepath.Join(root, "..", "outside_file.txt"),
		},
		{
			name: "repo root itself",
			in:   root,
			want: ".",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveRepoRelativePath(tt.in)
			if got != tt.want {
				t.Errorf("ResolveRepoRelativePath(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestOpenCaching(t *testing.T) {
	ClearCache()
	defer ClearCache()

	// First call opens repo
	repo1, err := Open()
	if err != nil {
		t.Fatalf("First Open() error = %v", err)
	}
	if repo1 == nil {
		t.Fatal("First Open() returned nil repo")
	}

	// Second call should return cached repo
	repo2, err := Open()
	if err != nil {
		t.Fatalf("Second Open() error = %v", err)
	}

	// Should be the same pointer
	if repo1 != repo2 {
		t.Error("Second Open() returned different repo, expected cached repo")
	}
}

func TestClearCache(t *testing.T) {
	ClearCache()
	defer ClearCache()

	// Open and cache
	repo1, err := Open()
	if err != nil {
		t.Fatalf("First Open() error = %v", err)
	}

	// Clear cache
	ClearCache()

	// Open again - should be new repo
	repo2, err := Open()
	if err != nil {
		t.Fatalf("Second Open() error = %v", err)
	}

	// Pointers should be different after cache clear
	if repo1 == repo2 {
		t.Error("After ClearCache(), Open() should return a new repo")
	}
}

// BenchmarkOpen benchmarks git repository opening with caching.
// This demonstrates the performance improvement from caching.
func BenchmarkOpen(b *testing.B) {
	b.Run("cached", func(b *testing.B) {
		ClearCache()
		defer ClearCache()

		// Warm up cache
		_, _ = Open()

		b.ResetTimer()
		b.ReportAllocs()
		for b.Loop() {
			_, _ = Open()
		}
	})

	b.Run("uncached", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			ClearCache()
			_, _ = Open()
		}
	})
}

// BenchmarkCurrentBranch benchmarks getting the current branch.
// This benefits from Open() caching.
func BenchmarkCurrentBranch(b *testing.B) {
	ClearCache()
	defer ClearCache()

	// Warm up cache
	_, _ = CurrentBranch()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_, _ = CurrentBranch()
	}
}

// BenchmarkRoot benchmarks getting the git root.
func BenchmarkRoot(b *testing.B) {
	ClearCache()
	defer ClearCache()

	// Warm up cache
	_ = Root()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = Root()
	}
}

// BenchmarkRemoteURL benchmarks getting the remote URL.
func BenchmarkRemoteURL(b *testing.B) {
	ClearCache()
	defer ClearCache()

	// Warm up cache
	_, _ = RemoteURL()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_, _ = RemoteURL()
	}
}
