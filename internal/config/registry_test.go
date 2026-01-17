package config

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveCommandConfigs(t *testing.T) {
	// Clear registry for this test
	registryMu.Lock()
	originalRegistry := registry
	registry = make(map[string]*CommandConfig)
	registryMu.Unlock()
	defer func() {
		registryMu.Lock()
		registry = originalRegistry
		registryMu.Unlock()
	}()

	var testConfig TestConfig
	Register(CommandConfig{
		Path:   "pull-request.list",
		Desc:   "test config",
		Target: &testConfig,
	})

	data := map[string]any{
		"pull-request": map[string]any{
			"list": map[string]any{
				"output": "json",
				"mine":   true,
			},
		},
	}

	err := resolveCommandConfigs(nil, data)
	require.NoError(t, err)

	assert.Equal(t, "json", testConfig.Output)
	assert.True(t, testConfig.Mine)
}

func TestGetNestedValue(t *testing.T) {
	data := map[string]any{
		"level1": map[string]any{
			"level2": map[string]any{
				"value": "found",
			},
		},
	}

	assert.Equal(t, "found", nestedValue(data, "level1.level2.value"))
	assert.Equal(t, map[string]any{"value": "found"}, nestedValue(data, "level1.level2"))
	assert.Nil(t, nestedValue(data, "nonexistent"))
	assert.Nil(t, nestedValue(data, "level1.nonexistent"))
}

func TestRegisterCommandConfigPanicsOnDuplicate(t *testing.T) {
	// Clear registry for this test
	registryMu.Lock()
	originalRegistry := registry
	registry = make(map[string]*CommandConfig)
	registryMu.Unlock()
	defer func() {
		registryMu.Lock()
		registry = originalRegistry
		registryMu.Unlock()
	}()

	var cfg1, cfg2 TestConfig

	Register(CommandConfig{
		Path:   "test.path",
		Desc:   "first",
		Target: &cfg1,
	})

	assert.Panics(t, func() {
		Register(CommandConfig{
			Path:   "test.path",
			Desc:   "second",
			Target: &cfg2,
		})
	})
}

type TestConfig struct {
	Output string `yaml:"output"`
	Mine   bool   `yaml:"mine"`
}

func (t TestConfig) OnResolved(c *cobra.Command) error {
	return nil
}
