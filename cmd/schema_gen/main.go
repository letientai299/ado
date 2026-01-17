package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/invopop/jsonschema"
	"github.com/letientai299/ado/internal/commands/pipeline"
	"github.com/letientai299/ado/internal/commands/pull_request"
	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/util/gitcli"
)

func main() {
	// 1. Discover command-specific configs via AST
	// This helps us know WHICH paths are registered and what their struct names are,
	// but we still rely on reflection via Registered targets to get the full schema.
	targets, err := discoverConfigs("internal/commands")
	if err != nil {
		log.Fatalf("failed to discover configs: %v", err)
	}

	// We still need to call Cmd() to actually populate the config.Registry()
	// because we want to use the reflect-based approach easily, and we need
	// actual instances to get reflect.Type.
	// NOTE: The config structs are now private (unexported), but because we
	// are calling the Cmd() functions from their respective packages,
	// they correctly register themselves with config.Register.
	_ = pull_request.Cmd()
	_ = pipeline.Cmd()

	registry := config.Registry()
	log.Infof("Found %d registered targets in registry", len(registry))

	reflector := jsonschema.Reflector{
		DoNotReference:             true,
		ExpandedStruct:             true,
		RequiredFromJSONSchemaTags: true,
	}
	if err = reflector.AddGoComments("github.com/letientai299/ado", "./"); err != nil {
		log.Infof("Warning: failed to add Go comments: %v", err)
	}

	s := reflector.Reflect(&config.Config{})

	for path, cmdCfg := range registry {
		if cmdCfg.Target == nil {
			continue
		}

		targetType := reflect.TypeOf(cmdCfg.Target)
		if targetType.Kind() == reflect.Ptr {
			targetType = targetType.Elem()
		}

		targetSchema := reflector.ReflectFromType(targetType)
		if targetSchema.Description == "" {
			targetSchema.Description = cmdCfg.Desc
		}

		addNestedProperty(s, path, targetSchema)
	}

	// Post-process schema to add include! and markdownDescription
	processSchema(s)

	// Check if AST discovery found something that is NOT in the registry
	for path, typeName := range targets {
		if _, ok := registry[path]; !ok {
			log.Infof(
				"Warning: AST found config at %q (type %s) but it's not in registry. Did you forget to call its Cmd() in schema_gen/main.go?",
				path,
				typeName,
			)
		}
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	gitRoot, err := gitcli.Root()
	if err != nil {
		log.Fatalf("failed to find git root: %v", err)
	}

	outputPath := filepath.Join(gitRoot, "etc", "schemas", "config.json")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o750); err != nil {
		log.Fatalf("failed to create directory: %v", err)
	}

	if err = os.WriteFile(outputPath, data, 0o600); err != nil {
		log.Fatalf("failed to write schema file: %v", err)
	}
	log.Infof("Schema generated at %s", outputPath)
}

func discoverConfigs(root string) (map[string]string, error) {
	configs := make(map[string]string)
	fileSet := token.NewFileSet()

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		f, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			// Look for config.Register(...)
			if x, ok := sel.X.(*ast.Ident); !ok || x.Name != "config" || sel.Sel.Name != "Register" {
				return true
			}

			if len(call.Args) != 1 {
				return true
			}

			// Expecting config.CommandConfig literal
			comp, ok := call.Args[0].(*ast.CompositeLit)
			if !ok {
				return true
			}

			var configPath, targetType string
			for _, elt := range comp.Elts {
				kv, ok := elt.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				key, ok := kv.Key.(*ast.Ident)
				if !ok {
					continue
				}

				switch key.Name {
				case "Path":
					if lit, ok := kv.Value.(*ast.BasicLit); ok {
						configPath = strings.Trim(lit.Value, "\"")
					}
				case "Target":
					// Usually &SomeStruct{} or opts (where opts is *SomeStruct)
					targetType = exprToString(kv.Value)
				}
			}

			if configPath != "" {
				configs[configPath] = targetType
			}

			return true
		})

		return nil
	})

	return configs, err
}

func exprToString(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.UnaryExpr:
		return exprToString(v.X)
	case *ast.CompositeLit:
		return exprToString(v.Type)
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", exprToString(v.X), v.Sel.Name)
	default:
		return fmt.Sprintf("%T", v)
	}
}

func addNestedProperty(base *jsonschema.Schema, path string, prop *jsonschema.Schema) {
	partsList := strings.Split(path, ".")

	current := base
	for i, part := range partsList {
		if current.Properties == nil {
			current.Properties = jsonschema.NewProperties()
		}

		if i == len(partsList)-1 {
			current.Properties.Set(part, prop)
		} else {
			existing, ok := current.Properties.Get(part)
			if !ok {
				next := &jsonschema.Schema{
					Type:       "object",
					Properties: jsonschema.NewProperties(),
				}
				current.Properties.Set(part, next)
				current = next
			} else {
				current = existing
			}
		}
	}
}

func processSchema(s *jsonschema.Schema) {
	if s == nil {
		return
	}

	if s.Description != "" {
		if s.Extras == nil {
			s.Extras = make(map[string]any)
		}
		s.Extras["markdownDescription"] = s.Description
		s.Description = ""
	}

	if s.Type == "object" {
		if s.Properties == nil {
			s.Properties = jsonschema.NewProperties()
		}
		s.Properties.Set("include!", &jsonschema.Schema{
			Type: "string",
			Extras: map[string]any{
				"markdownDescription": "Load config from another file",
			},
		})
	}

	if s.Properties != nil {
		for pair := s.Properties.Oldest(); pair != nil; pair = pair.Next() {
			processSchema(pair.Value)
		}
	}

	if s.Items != nil {
		processSchema(s.Items)
	}

	if s.AnyOf != nil {
		for _, sub := range s.AnyOf {
			processSchema(sub)
		}
	}

	if s.AllOf != nil {
		for _, sub := range s.AllOf {
			processSchema(sub)
		}
	}

	if s.OneOf != nil {
		for _, sub := range s.OneOf {
			processSchema(sub)
		}
	}
}
