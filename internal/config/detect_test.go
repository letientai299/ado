package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectTenant(t *testing.T) {
	tests := []struct {
		name     string
		bash     func(script string) (string, error)
		initial  Config
		expected Config
		wantErr  bool
	}{
		{
			name: "tenant and username missing",
			bash: func(script string) (string, error) {
				return "tenant-id\tuser@example.com", nil
			},
			expected: Config{
				Tenant:   "tenant-id",
				Username: "user@example.com",
			},
		},
		{
			name: "tenant already set",
			bash: func(script string) (string, error) {
				return "tenant-id\tuser@example.com", nil
			},
			initial: Config{
				Tenant: "existing-tenant",
			},
			expected: Config{
				Tenant:   "existing-tenant",
				Username: "user@example.com",
			},
		},
		{
			name: "username already set",
			bash: func(script string) (string, error) {
				return "tenant-id\tuser@example.com", nil
			},
			initial: Config{
				Username: "existing-user",
			},
			expected: Config{
				Tenant:   "tenant-id",
				Username: "existing-user",
			},
		},
		{
			name: "both already set - should skip bash",
			bash: func(script string) (string, error) {
				return "", errors.New("should not be called")
			},
			initial: Config{
				Tenant:   "existing-tenant",
				Username: "existing-user",
			},
			expected: Config{
				Tenant:   "existing-tenant",
				Username: "existing-user",
			},
		},
		{
			name: "bash error",
			bash: func(script string) (string, error) {
				return "", errors.New("bash error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := tt.initial
			err := detectTenant(&cfg, tt.bash)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, cfg)
			}
		})
	}
}

func TestDetectRepo(t *testing.T) {
	tests := []struct {
		name     string
		bash     func(script string) (string, error)
		initial  Config
		expected Config
		wantErr  bool
	}{
		{
			name: "repo info missing",
			bash: func(script string) (string, error) {
				return "https://dev.azure.com/org/project/_git/repo", nil
			},
			expected: Config{
				Repository: Repository{
					Org:     "org",
					Project: "project",
					Name:    "repo",
				},
			},
		},
		{
			name: "repo info already set",
			bash: func(script string) (string, error) {
				return "", errors.New("should not be called")
			},
			initial: Config{
				Repository: Repository{
					Name:    "existing-repo",
					Org:     "existing-org",
					Project: "existing-project",
				},
			},
			expected: Config{
				Repository: Repository{
					Name:    "existing-repo",
					Org:     "existing-org",
					Project: "existing-project",
				},
			},
		},
		{
			name: "bash error",
			bash: func(script string) (string, error) {
				return "", errors.New("bash error")
			},
			wantErr: true,
		},
		{
			name: "parse error",
			bash: func(script string) (string, error) {
				return "invalid-url", nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := tt.initial
			err := detectRepo(&cfg, tt.bash)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, cfg)
			}
		})
	}
}
