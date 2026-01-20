package config_cmd

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/ui"
	"github.com/letientai299/ado/internal/util"
	"github.com/letientai299/ado/internal/util/editor"
	"github.com/letientai299/ado/internal/util/gitcli"
	"github.com/spf13/cobra"
)

//go:embed config.md
var configDoc string

//go:embed config.editor.md
var configEditorDoc string

//go:embed config.theme.md
var configThemeDoc string

//go:generate go run ../../../cmd/schema_gen
//go:embed init.ado.yml
var initAdoYAML string

// Cmd returns the config command group
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg"},
		Short:   "Manage ADO configuration",
		Long:    configDoc,
	}
	cmd.AddCommand(dumpCmd())
	cmd.AddCommand(initCmd())
	cmd.AddCommand(editCmd())

	cmd.AddCommand(util.HelpTopic("editor", configEditorDoc))
	cmd.AddCommand(util.HelpTopic("theme", configThemeDoc))
	return cmd
}

func initCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a config file with defaults",
		Long:  "Create a new .ado.yaml config file with default values and documentation.",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := createConfigFile()
			return err
		},
	}
	return cmd
}

func editCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "edit",
		Aliases: []string{"e"},
		Short:   "Edit the config file in the configured editor",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.From(cmd.Context())
			file := cfg.FilePath()
			if file != "" {
				return editor.Open(cfg.Editor, file)
			}

			if !ui.Confirm("Config file not found. Create one?", true) {
				return nil
			}

			var err error
			file, err = createConfigFile()
			if err != nil {
				return err
			}

			return editor.Open(cfg.Editor, file)
		},
	}
	return c
}

func createConfigFile() (string, error) {
	file, err := config.FindConfigFile()
	hasConfig := err == nil && file != ""
	if hasConfig {
		log.Info("Config file exist.", "file", file)
		log.Info("Use `ado config edit` to open it")
		return file, nil
	}

	gitRoot := gitcli.Root()
	dotConfigDir := filepath.Join(gitRoot, ".config")

	if _, err = os.Stat(dotConfigDir); err != nil {
		file = filepath.Join(gitRoot, ".ado.yaml")
	} else {
		file, err = chooseLocation(dotConfigDir, gitRoot)
		if err != nil {
			return "", err
		}
	}

	if err = os.WriteFile(file, []byte(initAdoYAML), 0o600); err != nil {
		return file, fmt.Errorf("writing config file: %w", err)
	}

	fmt.Printf("Created %s\n", file)
	return file, nil
}

func chooseLocation(dotConfigDir, gitRoot string) (string, error) {
	options := []string{
		filepath.Join(dotConfigDir, "ado.yaml"),
		filepath.Join(gitRoot, ".ado.yaml"),
	}

	choice := ui.Pick(options, ui.PickConfig[string]{
		Title: "Where to put the config file?",
		Render:      func(w io.Writer, it string, matches []int) { _, _ = fmt.Fprintf(w, it) },
		FilterValue: func(item string) string { return item },
	})

	if choice.IsNil() {
		return "", errors.New("cancelled")
	}

	return choice.Get(), nil
}

func dumpCmd() *cobra.Command {
	type complete struct {
		*config.Config `yaml:",inline"`
		// CommandConfigs map[string]any `yaml:",inline" json:","`
		CommandConfigs map[string]any `yaml:",inline"`
	}

	c := &cobra.Command{
		Use:     "dump",
		Aliases: []string{"d"},
		Short:   "Dump the resolved config",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.From(cmd.Context())
			registry := config.Registry()
			m := make(map[string]any, len(registry))
			for k, v := range registry {
				m[k] = v.Target
			}

			all := complete{
				Config:         cfg,
				CommandConfigs: m,
			}

			return styles.DumpYAML(all)
		},
	}

	return c
}
