// Package gitcli provides utilities to query git repo or query data via Git CLI.
package gitcli

import (
	"os"

	"github.com/letientai299/ado/internal/util/sh"
)

// TryRoot finds the git repo root, or fallback to current working if fail
func TryRoot() string {
	gitRoot, err := Root()
	if err == nil {
		return gitRoot
	}
	wd, _ := os.Getwd()
	return wd
}

func Root() (string, error) {
	return sh.Run("git rev-parse --show-toplevel")
}
