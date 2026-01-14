package config

import (
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadYAMLWithIncludes(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		expected string
		wantErr  string
	}{
		{
			name:     "basic inclusion",
			filePath: "testdata/basic/main.yaml",
			expected: "tenant: mytenant\nusername: testuser\ndebug: true\n",
		},
		{
			name:     "quoted path inclusion",
			filePath: "testdata/quoted/main.yaml",
			expected: "username: testuser\n",
		},
		{
			name:     "nested inclusion",
			filePath: "testdata/nested/f1.yaml",
			expected: "key1: val1\nkey2: val2\nkey3: val3\n",
		},
		{
			name:     "circular dependency",
			filePath: "testdata/circular/f1.yaml",
			wantErr:  "circular include detected",
		},
		{
			name:     "multiple inclusions at different levels",
			filePath: "testdata/multiple/main.yaml",
			expected: "a: 1\nb: 2\nsection1:\n  key1: val1\nsection2:\n  key2: val2\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := loadYAMLWithIncludes(tt.filePath, make(map[string]struct{}))
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}

			require.NoError(t, err)

			var actualMap, expectedMap map[string]any
			require.NoError(t, yaml.Unmarshal(data, &actualMap))
			require.NoError(t, yaml.Unmarshal([]byte(tt.expected), &expectedMap))
			assert.Equal(t, expectedMap, actualMap)
		})
	}
}
