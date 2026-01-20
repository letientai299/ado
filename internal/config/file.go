package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/goccy/go-yaml"
	"github.com/letientai299/ado/internal/util/gitcli"
)

const includeDirective = "include!"

// resolveConfigFile finds the YAML config file and loads it with include support.
func resolveConfigFile(cfg *Config) error {
	filePath, err := FindConfigFile()
	if err != nil {
		return err
	}

	if filePath == "" {
		return nil // no config file found
	}

	log.Debugf("found config file %v", filePath)
	cfg.filePath = filePath

	data, err := loadYAMLWithIncludes(filePath, make(map[string]struct{}))
	if err != nil {
		return fmt.Errorf("loading config file: %w", err)
	}

	if err = yaml.UnmarshalWithOptions(data, cfg); err != nil {
		return fmt.Errorf("parsing config file: %w", err)
	}

	// Resolve command-specific configs from the raw data
	var rawData map[string]any
	if err = yaml.Unmarshal(data, &rawData); err != nil {
		return fmt.Errorf("parsing config for command configs: %w", err)
	}
	if err = resolveCommandConfigs(cfg.cmd, rawData); err != nil {
		return fmt.Errorf("resolving command configs: %w", err)
	}

	return nil
}

// loadYAMLWithIncludes loads a YAML file, processes any `include!` directives,
// and returns the final YAML bytes.
func loadYAMLWithIncludes(path string, visited map[string]struct{}) ([]byte, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolving path %s: %w", path, err)
	}

	if _, ok := visited[absPath]; ok {
		return nil, fmt.Errorf("circular include detected: %s", absPath)
	}
	visited[absPath] = struct{}{}

	rawData, err := os.ReadFile(
		filepath.Clean(absPath),
	) // nolint:gosec // G304: Potential file inclusion via variable. we trust the path here.
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	// Parse to map for include processing
	var data map[string]any
	if err = yaml.Unmarshal(rawData, &data); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	baseDir := filepath.Dir(absPath)
	if err = processIncludes(data, baseDir, visited); err != nil {
		return nil, err
	}

	return yaml.Marshal(data)
}

// processIncludes recursively processes maps looking for `include!` directives.
// It merges the included content into the current mapping.
func processIncludes(data map[string]any, baseDir string, visited map[string]struct{}) error {
	for key, value := range data {
		// Handle `include!` directive
		if key == includeDirective {
			includePath, ok := value.(string)
			if !ok {
				continue
			}

			// Clone visited the map for this branch
			branchVisited := make(map[string]struct{})
			for k, v := range visited {
				branchVisited[k] = v
			}

			absIncludePath := filepath.Join(baseDir, includePath)
			includedData, err := loadYAMLWithIncludes(absIncludePath, branchVisited)
			if err != nil {
				return fmt.Errorf("processing include %s: %w", includePath, err)
			}

			var includedMap map[string]any
			if err = yaml.Unmarshal(includedData, &includedMap); err != nil {
				return fmt.Errorf("parsing included file %s: %w", includePath, err)
			}

			// Merge included map into data (existing keys take precedence)
			for k, v := range includedMap {
				if _, exists := data[k]; !exists {
					data[k] = v
				}
			}

			// Remove the `include!` key after processing
			delete(data, includeDirective)
			continue
		}

		// Recursively process nested maps
		if nestedMap, ok := value.(map[string]any); ok {
			if err := processIncludes(nestedMap, baseDir, visited); err != nil {
				return err
			}
		}

		// Recursively process arrays that may contain maps
		if arr, ok := value.([]any); ok {
			for _, item := range arr {
				if nestedMap, ok := item.(map[string]any); ok {
					if err := processIncludes(nestedMap, baseDir, visited); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

// FindConfigFile looks for .ado.y(a)ml or `.config/ado.y(a)ml` in the
// working dir, then continue the search up to the git root dir.
func FindConfigFile() (string, error) {
	wd, _ := os.Getwd()
	gitRoot := gitcli.Root()

	for {
		for _, f := range configFileNames {
			p := filepath.Join(wd, f)
			if _, err := os.Stat(p); err == nil {
				return p, nil
			}
		}

		if wd == gitRoot || wd == filepath.Dir(wd) {
			break
		}
		wd = filepath.Dir(wd)
	}

	return "", nil
}
