package config_cmd

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/util"
	"github.com/letientai299/ado/internal/util/editor"
	"github.com/spf13/cobra"
)

//go:embed config.md
var configDoc string

//go:embed config.editor.md
var configEditorDoc string

//go:embed config.theme.md
var configThemeDoc string

const configTemplate = `
# ADO CLI Configuration
# See: etc/schemas/config.json for full schema documentation

# Debug mode enables verbose logging
# debug: false

# Azure tenant ID (auto-detected from az CLI if not set)
# tenant: ""

# Repository settings (auto-detected from git remote if not set)
# repository:
#   org: ""
#   project: ""
#   name: ""

# Theme configuration
# theme:
#   name: tokyo-night
#   true_color: true

# Command-specific configuration (optional)
# pull-request:
#   list:
#     default_output: simple
#     custom_output_templates: {}
`

// Cmd returns the config command group
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg"},
		Short:   "Manage ADO configuration",
		Long:    configDoc,
	}

	cmd.AddCommand(initCmd())
	cmd.AddCommand(editCmd())

	cmd.AddCommand(util.HelpTopic("editor", configEditorDoc))
	cmd.AddCommand(util.HelpTopic("theme", configThemeDoc))
	return cmd
}

func initCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a config file with defaults",
		Long:  "Create a new .ado.yaml config file with default values and documentation.",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := ".ado.yaml"
			if _, err := os.Stat(configPath); err == nil && !force {
				return fmt.Errorf("config file already exists at %s, use --force to overwrite", configPath)
			}

			if err := os.WriteFile(configPath, []byte(configTemplate), 0o600); err != nil {
				return fmt.Errorf("writing config file: %w", err)
			}

			fmt.Printf("Config file created at %s\n", configPath)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite existing config file")
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
			if file == "" {
				return fmt.Errorf("config file not found")
			}
			return editor.Open(cfg.Editor, file)
		},
	}
	return c
}