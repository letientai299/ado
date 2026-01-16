package util

import "os"

// TryGitRoot finds the git repo root, or fallback to current working if fail
func TryGitRoot() string {
	gitRoot, err := GitRoot()
	if err == nil {
		return gitRoot
	}
	wd, _ := os.Getwd()
	return wd
}

func GitRoot() (string, error) {
	return Bash("git rev-parse --show-toplevel")
}
