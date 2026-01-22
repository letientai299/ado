// Package cache provides simple file-based JSON caching.
package cache

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/charmbracelet/log"
	"github.com/goccy/go-json"
)

const appName = "ado"

type entry[T any] struct {
	Data      T         `json:"data"`
	ExpiresAt time.Time `json:"expires_at"`
}

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

// Get retrieves a cached value by key. Returns false if not found, expired or error.
func Get[T any](key string) (*T, bool) {
	e := new(entry[T])
	path, err := cachePath(key)
	if err != nil {
		log.Debug("cache miss: failed to get cache path", "key", key, "err", err)
		return nil, false
	}

	data, err := os.ReadFile(path) //nolint:gosec
	if err != nil {
		log.Debug("cache miss", "key", key, "path", path)
		return nil, false
	}

	if err = json.Unmarshal(data, e); err != nil {
		log.Debug("cache error: failed to unmarshal", "key", key, "path", path, "err", err)
		return nil, false
	}

	if time.Now().After(e.ExpiresAt) {
		log.Debug("cache miss: expired", "key", key, "path", path, "expiry", e.ExpiresAt)
		_ = os.Remove(path)
		return nil, false
	}

	log.Debug("cache hit", "key", key, "path", path)
	return &e.Data, true
}

// Set stores a value in the cache by key with a TTL.
func Set(key string, v any, ttl time.Duration) error {
	path, err := cachePath(key)
	if err != nil {
		return err
	}

	log.Debug("cache set", "key", key, "path", path, "ttl", ttl)

	if err = os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}

	e := entry[any]{
		Data:      v,
		ExpiresAt: time.Now().Add(ttl),
	}

	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o600)
}
