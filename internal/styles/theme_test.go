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

	themesDir := filepath.Join(gitRoot, "etc/themes")
	for _, theme := range all {
		path := filepath.Join(themesDir, theme.Name+".yml")
		f, err := os.Create(path)
		require.NoError(t, err)
		enc := yaml.NewEncoder(f)
		_ = enc.Encode(theme)
		_ = f.Close()
	}
}
