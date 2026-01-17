// Package cache provides simple file-based JSON caching.
package cache

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/charmbracelet/log"
	"github.com/goccy/go-json"
)

const appName = "ado"

// cacheDir returns the cache directory for the app.
// Uses XDG_CACHE_HOME on Linux/Mac, LocalAppData on Windows.
func cacheDir() (string, error) {
	var base string

	if v := os.Getenv("XDG_CACHE_HOME"); v != "" {
		base = v
	} else if runtime.GOOS == "windows" {
		base = os.Getenv("LOCALAPPDATA")
		if base == "" {
			base = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
		}
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".cache")
	}

	return filepath.Join(base, appName), nil
}

func cachePath(key string) (string, error) {
	dir, err := cacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, key+".json"), nil
}

// Get retrieves a cached value by key. Returns false if not found or error.
func Get[T any](key string, v *T) bool {
	path, err := cachePath(key)
	if err != nil {
		log.Debug("cache miss: failed to get cache path", "key", key, "err", err)
		return false
	}

	data, err := os.ReadFile(path) //nolint:gosec
	if err != nil {
		log.Debug("cache miss", "key", key, "path", path)
		return false
	}

	if err = json.Unmarshal(data, v); err != nil {
		log.Debug("cache error: failed to unmarshal", "key", key, "path", path, "err", err)
		return false
	}

	log.Debug("cache hit", "key", key, "path", path)
	return true
}

// Set stores a value in the cache by key.
func Set(key string, v any) error {
	path, err := cachePath(key)
	if err != nil {
		return err
	}

	log.Debug("cache set", "key", key, "path", path)

	if err = os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o600)
}
