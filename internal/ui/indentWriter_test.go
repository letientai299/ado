package ui

import (
	"bytes"
	"testing"
)

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
