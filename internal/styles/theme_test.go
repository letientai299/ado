package styles_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/util"
	"github.com/stretchr/testify/require"
)

func Test_createThemeConfigFiles(t *testing.T) {
	all := []styles.Theme{styles.Light, styles.TokyoNight, styles.Dark, styles.NoTTy}
	gitRoot, err := util.GitRoot()
	require.NoError(t, err)

	type ThemeConfig struct {
		Theme styles.Theme `yaml:"theme"`
	}

	themesDir := filepath.Join(gitRoot, "etc/themes")
	for _, theme := range all {
		path := filepath.Join(themesDir, theme.Name+".yml")
		f, err := os.Create(
			filepath.Clean(path),
		) // nolint:gosec // G304: Potential file inclusion via variable. test code.
		require.NoError(t, err)
		enc := yaml.NewEncoder(f)
		_ = enc.Encode(ThemeConfig{Theme: theme})
		_ = f.Close()
	}
}
