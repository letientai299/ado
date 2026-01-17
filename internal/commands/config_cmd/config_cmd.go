package config_cmd

import (
	"fmt"
	"os"

	"github.com/letientai299/ado/internal/config/schema"
	"github.com/letientai299/ado/internal/styles"
	"github.com/spf13/cobra"
)

const configTemplate = `
# ADO CLI Configuration
# See: ado config schema --output - for full schema documentation

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
#     output: simple   # Output format: simple, json, yaml
#     mine: false      # Show only PRs created by you
#     draft: false     # Include draft PRs
`

// Cmd returns the config command group
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg"},
		Short:   "Manage ADO configuration",
		Long:    "Commands for managing ADO CLI configuration files and schemas.",
	}
	cmd.AddCommand(schemaCmd(), initCmd())
	return cmd
}

func schemaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Generate configuration schema",
		// TODO (tai): add long usage, also enhance the schema to contains comments and validation.
		RunE: func(cmd *cobra.Command, args []string) error {
			data := schema.Generate()
			return styles.DumpYAML(data)
		},
	}

	return cmd
}

func initCmd() *cobra.Command {
	var force bool

	// TODO (tai): this fails if git root isn't exist
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a config file with defaults",
		Long:  "Create a new .ado.yaml config file with default values and documentation.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO (tai): `config init` should ask for `.ado.yml` or `.config/ado.yml`, and other info
			//  e.g. gen ADO PAT, or use tennant ID.
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
