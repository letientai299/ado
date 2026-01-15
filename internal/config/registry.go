package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/goccy/go-yaml"
)

// CommandConfig represents a command's configuration registration.
type CommandConfig struct {
	// Path is the YAML path like "pull-request" or "pull-request.list"
	Path string
	// Desc is description for documentation/schema generation
	Desc string
	// Target is a pointer to the config struct to unmarshal into
	Target any
}

var (
	registryMu sync.RWMutex
	registry   = make(map[string]*CommandConfig)
)

// Register adds a command's config path.
// Panics if the path is already registered (prevents silent conflicts).
func Register(cfg CommandConfig) {
	registryMu.Lock()
	defer registryMu.Unlock()

	if existing, ok := registry[cfg.Path]; ok {
		panic(fmt.Sprintf(
			"duplicate config path %q: already registered (description: %s)",
			cfg.Path, existing.Desc,
		))
	}
	registry[cfg.Path] = &cfg
}

// Registry returns a copy of all registered command configs.
func Registry() map[string]*CommandConfig {
	registryMu.RLock()
	defer registryMu.RUnlock()

	result := make(map[string]*CommandConfig, len(registry))
	for k, v := range registry {
		result[k] = v
	}
	return result
}

// resolveCommandConfigs extracts command-specific configs from raw YAML data.
func resolveCommandConfigs(data map[string]any) error {
	registryMu.RLock()
	defer registryMu.RUnlock()

	for path, cmdCfg := range registry {
		value := nestedValue(data, path)
		if value == nil {
			continue
		}

		// Re-marshal the nested value and unmarshal into target
		bytes, err := yaml.Marshal(value)
		if err != nil {
			return fmt.Errorf("marshaling config for %s: %w", path, err)
		}

		if err := yaml.UnmarshalWithOptions(bytes, cmdCfg.Target, yaml.Strict()); err != nil {
			return fmt.Errorf("parsing config for %s: %w", path, err)
		}
	}

	return nil
}

// nestedValue retrieves a value from a nested map using dot-separated path.
func nestedValue(m map[string]any, path string) any {
	parts := strings.Split(path, ".")
	current := any(m)

	for _, part := range parts {
		if currentMap, ok := current.(map[string]any); ok {
			current = currentMap[part]
		} else {
			return nil
		}
	}
	return current
}
