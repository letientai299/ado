package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/letientai299/ado/internal/util"
)

// resolveConfigFile finds the YAML config file and loads it with include support.
func resolveConfigFile(cfg *Config) error {
	filePath, err := findConfigFile()
	if err != nil {
		return err
	}

	if filePath == "" {
		return nil // no config file found
	}

	log.Debugf("found config file %v", filePath)

	data, err := loadYAMLWithIncludes(filePath, make(map[string]struct{}))
	if err != nil {
		return fmt.Errorf("loading config file: %w", err)
	}

	if err = yaml.UnmarshalWithOptions(data, cfg, yaml.Strict()); err != nil {
		return fmt.Errorf("parsing config file: %w", err)
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

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	file, err := parser.ParseBytes(data, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	baseDir := filepath.Dir(absPath)
	for _, doc := range file.Docs {
		if err = processIncludes(doc.Body, baseDir, visited); err != nil {
			return nil, err
		}
	}

	return []byte(file.String()), nil
}

func processIncludes(node ast.Node, baseDir string, visited map[string]struct{}) error {
	if node == nil {
		return nil
	}

	mapping, ok := node.(*ast.MappingNode)
	if !ok {
		return processIncludesChildren(node, baseDir, visited)
	}

	newValues := make([]*ast.MappingValueNode, 0, len(mapping.Values))

	for _, mv := range mapping.Values {
		if mv.Key.String() == "include!" {
			includePath := filepath.Join(baseDir, mv.Value.String())

			// Clone visited map for this branch
			branchVisited := make(map[string]struct{})
			for k, v := range visited {
				branchVisited[k] = v
			}

			includedData, err := loadYAMLWithIncludes(includePath, branchVisited)
			if err != nil {
				return fmt.Errorf("processing include %s: %w", mv.Value.String(), err)
			}

			includedFile, err := parser.ParseBytes(includedData, 0)
			if err != nil {
				return fmt.Errorf("parsing included file %s: %w", includePath, err)
			}

			if len(includedFile.Docs) > 0 {
				if incMapping, ok := includedFile.Docs[0].Body.(*ast.MappingNode); ok {
					newValues = append(newValues, incMapping.Values...)
				}
			}
		} else {
			if err := processIncludes(mv.Value, baseDir, visited); err != nil {
				return err
			}
			newValues = append(newValues, mv)
		}
	}

	mapping.Values = newValues
	return nil
}

func processIncludesChildren(node ast.Node, baseDir string, visited map[string]struct{}) error {
	switch n := node.(type) {
	case *ast.MappingValueNode:
		return processIncludes(n.Value, baseDir, visited)
	case *ast.SequenceNode:
		for _, child := range n.Values {
			if err := processIncludes(child, baseDir, visited); err != nil {
				return err
			}
		}
	}
	return nil
}

// findConfigFile looks for .ado.y(a)ml or `.config/ado.y(a)ml` in the
// working dir, then continue the search up to the git root dir.
func findConfigFile() (string, error) {
	gitRoot, err := util.GitRoot()
	if err != nil {
		log.Warnf("fail to get git root dir: %v", err)
		return "", err
	}

	wd, _ := os.Getwd()

	for {
		for _, f := range configFileNames {
			p := filepath.Join(wd, f)
			if _, err = os.Stat(p); err == nil {
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
