package schema

import (
	"reflect"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/letientai299/ado/internal/config"
)

// PropertySchema represents a schema property
type PropertySchema struct {
	Type        string                    `yaml:"type,omitempty"`
	Description string                    `yaml:"description,omitempty"`
	Default     any                       `yaml:"default,omitempty"`
	Properties  map[string]PropertySchema `yaml:"properties,omitempty"`
}

// Schema represents a YAML schema
type Schema struct {
	Schema      string                    `yaml:"$schema"`
	Title       string                    `yaml:"title"`
	Description string                    `yaml:"description,omitempty"`
	Type        string                    `yaml:"type"`
	Properties  map[string]PropertySchema `yaml:"properties"`
}

// Generate creates a schema from the Config struct and registered command configs
func Generate() *Schema {
	schema := &Schema{
		Schema:      "http://json-schema.org/draft-07/schema#",
		Title:       "ADO CLI Configuration",
		Description: "Configuration schema for the Azure DevOps CLI",
		Type:        "object",
		Properties:  make(map[string]PropertySchema),
	}

	// Add properties from Config struct
	configType := reflect.TypeOf(config.Config{})
	addStructProperties(schema.Properties, configType)

	// Add properties from registered command configs
	for path, cmdCfg := range config.Registry() {
		if cmdCfg.Target == nil {
			continue
		}

		targetType := reflect.TypeOf(cmdCfg.Target)
		if targetType.Kind() == reflect.Ptr {
			targetType = targetType.Elem()
		}

		prop := PropertySchema{
			Type:        "object",
			Description: cmdCfg.Desc,
			Properties:  make(map[string]PropertySchema),
		}
		addStructProperties(prop.Properties, targetType)

		// Handle nested paths like "pull-request.list"
		parts := strings.Split(path, ".")
		addNestedProperty(schema.Properties, parts, prop)
	}

	return schema
}

// GenerateYAML outputs the schema as YAML bytes
func GenerateYAML() ([]byte, error) {
	schema := Generate()
	return yaml.Marshal(schema)
}

// addStructProperties extracts properties from a struct type
func addStructProperties(props map[string]PropertySchema, t reflect.Type) {
	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get YAML tag name
		yamlTag := field.Tag.Get("yaml")
		if yamlTag == "" || yamlTag == "-" {
			continue
		}
		name := strings.Split(yamlTag, ",")[0]

		prop := PropertySchema{
			Type: goTypeToSchemaType(field.Type),
		}

		// Recursively handle nested structs
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if fieldType.Kind() == reflect.Struct && prop.Type == "object" {
			prop.Properties = make(map[string]PropertySchema)
			addStructProperties(prop.Properties, fieldType)
		}

		props[name] = prop
	}
}

// addNestedProperty adds a property at a nested path
func addNestedProperty(props map[string]PropertySchema, path []string, prop PropertySchema) {
	if len(path) == 0 {
		return
	}

	if len(path) == 1 {
		// Merge with existing if present
		if existing, ok := props[path[0]]; ok {
			for k, v := range prop.Properties {
				if existing.Properties == nil {
					existing.Properties = make(map[string]PropertySchema)
				}
				existing.Properties[k] = v
			}
			props[path[0]] = existing
		} else {
			props[path[0]] = prop
		}
		return
	}

	// Create intermediate objects
	if _, ok := props[path[0]]; !ok {
		props[path[0]] = PropertySchema{
			Type:       "object",
			Properties: make(map[string]PropertySchema),
		}
	}

	existing := props[path[0]]
	if existing.Properties == nil {
		existing.Properties = make(map[string]PropertySchema)
	}
	addNestedProperty(existing.Properties, path[1:], prop)
	props[path[0]] = existing
}

// goTypeToSchemaType converts Go types to JSON Schema types
func goTypeToSchemaType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Map, reflect.Struct:
		return "object"
	case reflect.Ptr:
		return goTypeToSchemaType(t.Elem())
	default:
		return "string"
	}
}
