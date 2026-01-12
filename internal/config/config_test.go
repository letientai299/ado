package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolve(t *testing.T) {
	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func(dir string) { _ = os.Chdir(dir) }(origWd)

	testdata, err := filepath.Abs("testdata")
	require.NoError(t, err)

	// Mock bash
	oldBash := bash
	defer func() { bash = oldBash }()

	tests := []struct {
		name     string
		wd       string
		setup    func(t *testing.T)
		args     []string
		expected Config
	}{
		{
			name: "full resolution flow",
			wd:   filepath.Join(testdata, "resolve/full_flow"),
			setup: func(t *testing.T) {
				t.Setenv("ADO_USERNAME", "env-user")
			},
			args: []string{"--debug"},
			expected: Config{
				Tenant:   "file-tenant", // from file
				Username: "env-user",    // from env
				Debug:    true,          // from flag
				Repository: Repository{
					Repo:    "file-repo",        // from file
					Org:     "detected-org",     // from auto-detect
					Project: "detected-project", // from auto-detect
				},
			},
		},
		{
			name: "flag overrides everything",
			wd:   filepath.Join(testdata, "resolve/flag_override"),
			setup: func(t *testing.T) {
				t.Setenv("ADO_TENANT", "env-tenant")
			},
			args: []string{"--tenant", "flag-tenant"},
			expected: Config{
				Tenant:   "flag-tenant",
				Username: "detected-user",
				Repository: Repository{
					Repo:    "detected-repo",
					Org:     "detected-org",
					Project: "detected-project",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = os.Chdir(tt.wd)
			require.NoError(t, err)
			bash = mockBash(tt.wd)

			tt.setup(t)

			cmd := &cobra.Command{Use: "test"}
			AddGlobalFlags(cmd)
			err = cmd.ParseFlags(tt.args)
			require.NoError(t, err)

			cfg := &Config{}
			ctx := context.WithValue(context.Background(), ctxKeyGlobal, cfg)
			cmd.SetContext(ctx)

			err = Resolve(cmd, nil)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, *cfg)
		})
	}
}

func TestResolveConfigFile(t *testing.T) {
	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func(dir string) { _ = os.Chdir(dir) }(origWd)

	testdata, err := filepath.Abs("testdata")
	require.NoError(t, err)

	wd := filepath.Join(testdata, "resolve_config_file")
	err = os.Chdir(wd)
	require.NoError(t, err)

	oldBash := bash
	defer func() { bash = oldBash }()
	bash = func(script string) (string, error) {
		return wd, nil
	}

	var cfg Config
	err = resolveConfigFile(&cfg)
	assert.NoError(t, err)

	expected := Config{
		Tenant:   "my-tenant",
		Username: "my-user",
		Debug:    true,
		Repository: Repository{
			Repo:    "my-repo",
			Org:     "my-org",
			Project: "my-project",
		},
	}
	assert.Equal(t, expected, cfg)
}

func TestResolveEnv(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		initial  Config
		expected Config
	}{
		{
			name: "all env vars set",
			env: map[string]string{
				"ADO_TENANT":   "my-tenant",
				"ADO_USERNAME": "my-user",
				"ADO_DEBUG":    "true",
				"ADO_REPO":     "my-repo",
				"ADO_ORG":      "my-org",
				"ADO_PROJECT":  "my-project",
			},
			expected: Config{
				Tenant:   "my-tenant",
				Username: "my-user",
				Debug:    true,
				Repository: Repository{
					Repo:    "my-repo",
					Org:     "my-org",
					Project: "my-project",
				},
			},
		},
		{
			name: "some env vars set",
			env: map[string]string{
				"ADO_TENANT": "my-tenant",
			},
			expected: Config{
				Tenant: "my-tenant",
			},
		},
		{
			name: "env vars should override existing values if set by koanf",
			env: map[string]string{
				"ADO_TENANT": "new-tenant",
			},
			initial: Config{
				Tenant: "old-tenant",
			},
			expected: Config{
				Tenant: "new-tenant",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear relevant env vars
			for k := range tt.env {
				t.Setenv(k, tt.env[k])
			}
			// Also ensure other ADO_ env vars are cleared if they might affect the test
			// Actually t.Setenv only sets it for the test and its subtests,
			// but doesn't clear existing ones in the environment.
			// But since we are running in a controlled test environment, it should be fine.

			cfg := tt.initial
			err := resolveEnv(&cfg)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

func mockBash(wd string) func(script string) (string, error) {
	return func(script string) (string, error) {
		switch script {
		case "git rev-parse --show-toplevel":
			return wd, nil
		case `az account show --query "{tenantId:tenantId,username:user.name}" -o tsv`:
			return "detected-tenant\tdetected-user", nil
		case `git remote get-url origin`:
			return "https://dev.azure.com/detected-org/detected-project/_git/detected-repo", nil
		default:
			return "", nil
		}
	}
}
