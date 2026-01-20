package config_cmd

import (
	_ "embed"
	"errors"
	"fmt"
	"os"

	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/styles"
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

			if err := os.WriteFile(configPath, []byte(initAdoYAML), 0o600); err != nil {
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
				return errors.New("config file not found")
			}
			return editor.Open(cfg.Editor, file)
		},
	}
	return c
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
