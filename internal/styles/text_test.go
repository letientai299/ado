package styles

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestIndentTo(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		input    string
		expected string
	}{
		{
			name:     "single line",
			n:        2,
			input:    "hello",
			expected: "  hello",
		},
		{
			name:     "multiple lines",
			n:        2,
			input:    "hello\nworld",
			expected: "  hello\n  world",
		},
		{
			name:     "empty string",
			n:        2,
			input:    "",
			expected: "",
		},
		{
			name:     "trailing newline",
			n:        2,
			input:    "hello\n",
			expected: "  hello\n  ",
		},
		{
			name:     "multiple newlines",
			n:        2,
			input:    "hello\n\nworld",
			expected: "  hello\n\n  world",
		},
		{
			name:     "pre/post paragraph newlines",
			n:        2,
			input:    "\nhello\n\nworld\n",
			expected: "\n  hello\n\n  world\n  ",
		},
		{
			name:     "only newlines",
			n:        2,
			input:    "\n\n\n",
			expected: "\n\n\n  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			IndentTo(buf, tt.n, tt.input)
			if buf.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, buf.String())
			}
		})
	}
}

func BenchmarkIndentTo(b *testing.B) {
	s := strings.Repeat("hello world\n", 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var sb strings.Builder
		IndentTo(&sb, 2, s)
	}
}

func BenchmarkIndentWriter(b *testing.B) {
	s := []byte(strings.Repeat("hello world\n", 100))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iw := NewIndentWriter(io.Discard, "  ")
		_, _ = iw.Write(s)
	}
}

func TestIndentWriter(t *testing.T) {
	tests := []struct {
		name     string
		indent   string
		input    string
		input2   string
		expected string
	}{
		{
			name:     "single line",
			indent:   "  ",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "multiple lines",
			indent:   "  ",
			input:    "hello\nworld",
			expected: "hello\n  world",
		},
		{
			name:     "trailing newline",
			indent:   "  ",
			input:    "hello\nworld\n",
			expected: "hello\n  world\n",
		},
		{
			name:     "multiple writes",
			indent:   "  ",
			input:    "hello\n",
			input2:   "world",
			expected: "hello\n  world",
		},
		{
			name:     "multiple lines in second write",
			indent:   "  ",
			input:    "hello\n",
			input2:   "world\nagain",
			expected: "hello\n  world\n  again",
		},
		{
			name:     "multiple newlines",
			indent:   "  ",
			input:    "hello\n\nworld",
			expected: "hello\n\n  world",
		},
		{
			name:     "pre/post paragraph newlines",
			indent:   "  ",
			input:    "\nhello\n\nworld\n",
			expected: "\n  hello\n\n  world\n",
		},
		{
			name:     "only newlines",
			indent:   "  ",
			input:    "\n\n\n",
			expected: "\n\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			iw := NewIndentWriter(buf, tt.indent)
			_, err := iw.Write([]byte(tt.input))
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.input2 != "" {
				_, err = iw.Write([]byte(tt.input2))
				if err != nil {
					t.Errorf("unexpected error on second write: %v", err)
				}
			}
			if buf.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, buf.String())
			}
		})
	}
}
