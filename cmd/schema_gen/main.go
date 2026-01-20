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
	"sort"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/invopop/jsonschema"
	"github.com/letientai299/ado/internal/commands"
	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/util/gitcli"
)

func main() {
	gitRoot := gitcli.Root()
	_ = os.Chdir(gitRoot)

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
	_ = commands.Root()

	registry := config.Registry()
	log.Infof("Found %d registered targets in registry", len(registry))

	reflector := jsonschema.Reflector{
		RequiredFromJSONSchemaTags: true,
		DoNotReference:             true,
		ExpandedStruct:             true,
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

		setDefaultsFromInstance(targetSchema, cmdCfg.Target)
		addNestedProperty(s, path, targetSchema)
	}

	// Post-process schema to add include! and markdownDescription
	processSchema(s)

	// Reorder properties: global config first, then command configs, alphabetically
	reorderSchemaProperties(s)

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

	// Write output files
	writeFile(filepath.Join(gitRoot, "etc", "schemas", "config.json"), mustMarshalJSON(s))
	writeFile(
		filepath.Join(gitRoot, "internal/commands/config_cmd/init.ado.yml"),
		[]byte(generateExampleYAML(s)),
	)
}

func mustMarshalJSON(v any) []byte {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func writeFile(path string, data []byte) {
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		log.Fatalf("failed to create directory: %v", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		log.Fatalf("failed to write %s: %v", path, err)
	}
	log.Infof("Generated %s", path)
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

// setDefaultsFromInstance extracts default values from a config instance and sets them in the
// schema. It marshals the instance to JSON to get a map, then walks the map to set defaults.
func setDefaultsFromInstance(schema *jsonschema.Schema, instance any) {
	// Marshal to JSON then unmarshal to map to get field values with proper names
	data, err := json.Marshal(instance)
	if err != nil {
		return
	}

	var values map[string]any
	if err := json.Unmarshal(data, &values); err != nil {
		return
	}

	setDefaultsFromMap(schema, values)
}

func setDefaultsFromMap(schema *jsonschema.Schema, values map[string]any) {
	if schema == nil || schema.Properties == nil {
		return
	}
	for key, value := range values {
		prop, _ := schema.Properties.Get(key)
		if prop == nil {
			continue
		}
		if nested, ok := value.(map[string]any); ok {
			setDefaultsFromMap(prop, nested)
		} else if str, ok := value.(string); ok && str != "" {
			prop.Default = str
		} else if value != nil {
			prop.Default = value
		}
	}
}

func addNestedProperty(base *jsonschema.Schema, path string, prop *jsonschema.Schema) {
	parts := strings.Split(path, ".")
	current := base

	// Navigate/create intermediate objects
	for _, part := range parts[:len(parts)-1] {
		if current.Properties == nil {
			current.Properties = jsonschema.NewProperties()
		}
		if next, ok := current.Properties.Get(part); ok {
			current = next
		} else {
			next := &jsonschema.Schema{Type: "object", Properties: jsonschema.NewProperties()}
			current.Properties.Set(part, next)
			current = next
		}
	}

	// Set the final property
	if current.Properties == nil {
		current.Properties = jsonschema.NewProperties()
	}
	current.Properties.Set(parts[len(parts)-1], prop)
}

func processSchema(s *jsonschema.Schema) {
	if s == nil {
		return
	}

	// Move description to markdownDescription for IDE support
	if s.Description != "" {
		if s.Extras == nil {
			s.Extras = make(map[string]any)
		}
		s.Extras["markdownDescription"] = s.Description
		s.Description = ""
	}

	// Add include! directive for objects
	if s.Type == "object" {
		if s.Properties == nil {
			s.Properties = jsonschema.NewProperties()
		}
		s.Properties.Set("include!", &jsonschema.Schema{
			Type:   "string",
			Extras: map[string]any{"markdownDescription": "Load config from another file"},
		})
	}

	// Recurse into nested schemas
	for pair := s.Properties.Oldest(); pair != nil; pair = pair.Next() {
		processSchema(pair.Value)
	}
	processSchema(s.Items)
	for _, sub := range s.AnyOf {
		processSchema(sub)
	}
	for _, sub := range s.AllOf {
		processSchema(sub)
	}
	for _, sub := range s.OneOf {
		processSchema(sub)
	}
}

// cachedGlobalKeys holds the keys derived from config.Config struct.
var cachedGlobalKeys map[string]bool

// getGlobalConfigKeys returns the top-level keys from the Config struct.
func getGlobalConfigKeys() map[string]bool {
	if cachedGlobalKeys != nil {
		return cachedGlobalKeys
	}
	cachedGlobalKeys = make(map[string]bool)
	t := reflect.TypeOf(config.Config{})
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("yaml")
		if tag == "" {
			tag = t.Field(i).Tag.Get("json")
		}
		if name := strings.Split(tag, ",")[0]; name != "" {
			cachedGlobalKeys[name] = true
		}
	}
	return cachedGlobalKeys
}

// sortKeys returns keys sorted with global config keys first, then command keys.
// Within each group, keys are sorted alphabetically.
func sortKeys(keys []string) []string {
	globalKeys := getGlobalConfigKeys()
	var global, cs []string
	for _, k := range keys {
		if globalKeys[k] {
			global = append(global, k)
		} else {
			cs = append(cs, k)
		}
	}
	sort.Strings(global)
	sort.Strings(cs)
	return append(global, cs...)
}

// reorderSchemaProperties reorders properties in the schema to have stable key ordering.
// Global config keys come first, then command-level keys, all alphabetically within groups.
func reorderSchemaProperties(s *jsonschema.Schema) {
	if s == nil {
		return
	}

	if s.Properties != nil {
		// Collect and sort keys
		var keys []string
		for pair := s.Properties.Oldest(); pair != nil; pair = pair.Next() {
			keys = append(keys, pair.Key)
		}

		// Rebuild properties in sorted order
		oldProps := s.Properties
		s.Properties = jsonschema.NewProperties()
		for _, key := range sortKeys(keys) {
			if prop, ok := oldProps.Get(key); ok {
				s.Properties.Set(key, prop)
				reorderSchemaProperties(prop)
			}
		}
	}

	// Recurse into nested schemas
	reorderSchemaProperties(s.Items)
	for _, sub := range s.AnyOf {
		reorderSchemaProperties(sub)
	}
	for _, sub := range s.AllOf {
		reorderSchemaProperties(sub)
	}
	for _, sub := range s.OneOf {
		reorderSchemaProperties(sub)
	}
}

// staticExamples maps field keys to static example values.
// Use this for fields that need special formatting or commented-out examples.
var staticExamples = map[string]string{
	"theme": `## Use include! to load theme from external files.
## See https://github.com/letientai299/ado/tree/main/etc/themes for available themes.
## See also 'ado help config theme'
#theme:
#  include!: path/to/some_theme.yml`,
}

const yamlHeader = "# ADO CLI Configuration file. See `ado help config`.\n"

func generateExampleYAML(s *jsonschema.Schema) string {
	var sb strings.Builder
	sb.WriteString(yamlHeader)
	writeSchemaAsYAML(&sb, s, 0, "")
	return sb.String()
}

func writeSchemaAsYAML(sb *strings.Builder, s *jsonschema.Schema, indent int, path string) {
	if s == nil || s.Properties == nil {
		return
	}

	indentStr := strings.Repeat("  ", indent)

	// Collect property keys (skip include! directive)
	var keys []string
	for pair := s.Properties.Oldest(); pair != nil; pair = pair.Next() {
		key := pair.Key
		if key == "include!" {
			continue // Skip the `include!` directive in the properties listing
		}
		keys = append(keys, key)
	}
	// Sort: global config keys first, then command keys, alphabetically within each group
	keys = sortKeys(keys)

	for i, key := range keys {
		prop, _ := s.Properties.Get(key)
		if prop == nil {
			continue
		}

		currentPath := key
		if path != "" {
			currentPath = path + "." + key
		}

		// Check for static example override
		if staticVal, ok := staticExamples[currentPath]; ok {
			_, _ = fmt.Fprintf(sb, "%s\n", staticVal)
			if indent == 0 && i < len(keys)-1 {
				sb.WriteString("\n")
			}
			continue
		}

		// Write description as double-comment (## for docs)
		// Removing the first # gives valid YAML with proper comments
		if desc := getDescription(prop); desc != "" {
			hasContent := false
			for _, line := range strings.Split(desc, "\n") {
				if line = strings.TrimSpace(line); line != "" {
					_, _ = fmt.Fprintf(sb, "#%s# %s\n", indentStr, line)
					hasContent = true
				} else if hasContent {
					_, _ = fmt.Fprintf(sb, "#%s#\n", indentStr)
				}
			}
		}

		// Write the key and value (single # comment)
		switch {
		case isMapType(prop):
			// Map types get a multi-line example
			_, _ = fmt.Fprintf(sb, "#%s%s:\n", indentStr, key)
			_, _ = fmt.Fprintf(sb, "#%s  example_key: example_value\n", indentStr)
		case prop.Type == "object" && prop.Properties != nil && prop.Properties.Len() > 1:
			_, _ = fmt.Fprintf(sb, "#%s%s:\n", indentStr, key)
			writeSchemaAsYAML(sb, prop, indent+1, currentPath)
		default:
			_, _ = fmt.Fprintf(sb, "#%s%s: %s\n", indentStr, key, getDefaultValue(prop))
		}

		// Add a blank line between top-level keys only
		if indent == 0 && i < len(keys)-1 {
			sb.WriteString("\n")
		}
	}
}

// isMapType returns true if the schema represents a map type (object with additionalProperties
// schema). Regular structs have additionalProperties set to false (FalseSchema), while maps have a
// type schema.
func isMapType(s *jsonschema.Schema) bool {
	if s.Type != "object" || s.AdditionalProperties == nil {
		return false
	}
	// FalseSchema is used for additionalProperties: false (regular structs)
	// A real map has a schema with a Type defined
	return s.AdditionalProperties != jsonschema.FalseSchema && s.AdditionalProperties.Type != ""
}

func getDescription(s *jsonschema.Schema) string {
	if md, _ := s.Extras["markdownDescription"].(string); md != "" {
		return md
	}
	return s.Description
}

// getDefaultValue returns the default value from the schema, or an empty placeholder by type.
func getDefaultValue(s *jsonschema.Schema) string {
	if s.Default != nil {
		if str, ok := s.Default.(string); ok {
			return fmt.Sprintf("%q", str)
		}
		return fmt.Sprintf("%v", s.Default)
	}

	// Return an empty placeholder by type for fields without defaults
	switch s.Type {
	case "string":
		return `""`
	case "array":
		return "[]"
	case "object":
		return "{}"
	default:
		return ""
	}
}
