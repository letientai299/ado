package util_test

import (
	"testing"

	"github.com/letientai299/ado/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestParseRepoInfo(t *testing.T) {
	tests := []struct {
		origin      string
		wantOrg     string
		wantProject string
		wantRepo    string
		wantErr     bool
	}{
		{
			origin:      "https://dev.azure.com/org/project/_git/repo",
			wantOrg:     "org",
			wantProject: "project",
			wantRepo:    "repo",
			wantErr:     false,
		},
		{
			origin:      "https://org.visualstudio.com/project/_git/repo",
			wantOrg:     "org",
			wantProject: "project",
			wantRepo:    "repo",
			wantErr:     false,
		},
		{
			origin:      "https://skype.visualstudio.com/DefaultCollection/ES/_git/dummy",
			wantOrg:     "skype",
			wantProject: "ES",
			wantRepo:    "dummy",
			wantErr:     false,
		},
		{
			origin:      "https://dev.azure.com/msazure/One/_git/other",
			wantOrg:     "msazure",
			wantProject: "One",
			wantRepo:    "other",
			wantErr:     false,
		},
		{
			origin:      "git@ssh.dev.azure.com:v3/org/project/repo",
			wantOrg:     "org",
			wantProject: "project",
			wantRepo:    "repo",
			wantErr:     false,
		},
		{
			origin:      "https://org.visualstudio.com/DefaultCollection/project/_git/repo",
			wantOrg:     "org",
			wantProject: "project",
			wantRepo:    "repo",
			wantErr:     false,
		},
		{
			origin:  "invalid-url",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.origin, func(t *testing.T) {
			t.Parallel()

			gotOrg, gotProject, gotRepo, err := util.ParseRepoInfo(tt.origin)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantOrg, gotOrg)
				assert.Equal(t, tt.wantProject, gotProject)
				assert.Equal(t, tt.wantRepo, gotRepo)
			}
		})
	}
}
