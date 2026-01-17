package sh

import (
	"runtime"
	"testing"
)

func TestBash(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping bash test on windows")
	}

	out, err := Bash("echo 'hello'")
	if err != nil {
		t.Errorf("Bash failed: %v", err)
	}
	if out != "hello" {
		t.Errorf("Expected 'hello', got %q", out)
	}
}

func TestPwsh(t *testing.T) {
	// Check if pwsh is installed
	_, err := execShell("pwsh", "--version")
	if err != nil {
		t.Skip("pwsh not installed, skipping TestPwsh")
	}

	out, err := Pwsh("Write-Host 'hello'")
	if err != nil {
		t.Errorf("Pwsh failed: %v", err)
	}
	// Pwsh Write-Host often adds a newline or uses different formatting,
	// but execShell trims right spaces.
	if out != "hello" {
		t.Errorf("Expected 'hello', got %q", out)
	}
}
