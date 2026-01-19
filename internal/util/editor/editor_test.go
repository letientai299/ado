package editor

import (
	"testing"
)

func TestEditor_Edit(t *testing.T) {
	// Use a simple shell command as an editor that replaces content
	editor := New("test-edit-*", "echo 'modified content' >")

	original := "original content"
	updated, err := editor.Edit(original)
	if err != nil {
		t.Fatalf("Edit failed: %v", err)
	}

	expected := "modified content\n"
	if updated != expected {
		t.Errorf("expected %q, got %q", expected, updated)
	}
}

func TestEditor_Edit_NoChange(t *testing.T) {
	// Use 'true' as an editor (does nothing)
	editor := New("test-edit-*", "true")

	original := "original content"
	updated, err := editor.Edit(original)
	if err != nil {
		t.Fatalf("Edit failed: %v", err)
	}

	if updated != original {
		t.Errorf("expected %q, got %q", original, updated)
	}
}
