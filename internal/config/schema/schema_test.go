package schema

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	schema := Generate()

	assert.Equal(t, "http://json-schema.org/draft-07/schema#", schema.Schema)
	assert.Equal(t, "ADO CLI Configuration", schema.Title)
	assert.Equal(t, "object", schema.Type)

	// Check that basic properties from Config struct are present
	assert.Contains(t, schema.Properties, "debug")
	assert.Contains(t, schema.Properties, "tenant")
	assert.Contains(t, schema.Properties, "repository")
	assert.Contains(t, schema.Properties, "theme")

	// Check nested properties
	repoProps := schema.Properties["repository"]
	assert.Equal(t, "object", repoProps.Type)
	assert.Contains(t, repoProps.Properties, "org")
	assert.Contains(t, repoProps.Properties, "project")
	assert.Contains(t, repoProps.Properties, "name")
}

func TestGenerateYAML(t *testing.T) {
	data, err := GenerateYAML()
	require.NoError(t, err)
	assert.Contains(t, string(data), "$schema")
	assert.Contains(t, string(data), "ADO CLI Configuration")
}

func TestGoTypeToSchemaType(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"string", "", "string"},
		{"bool", false, "boolean"},
		{"int", 0, "integer"},
		{"float", 0.0, "number"},
		{"slice", []string{}, "array"},
		{"map", map[string]any{}, "object"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := goTypeToSchemaType(reflect.TypeOf(tt.input))
			assert.Equal(t, tt.expected, result)
		})
	}
}
