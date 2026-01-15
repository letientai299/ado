package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/styles"
	"github.com/spf13/cobra"
)

type ctxKey string

const (
	ctxKeyGlobal ctxKey = "global"
)

const (
	envAdoTenant = "ADO_TENANT"
	envAdoDebug  = "ADO_DEBUG"
	envAdoPat    = "ADO_PAT"
)

var configFileNames = []string{
	".ado.yml",
	".ado.yaml",
	".config/ado.yml",
	".config/ado.yaml",
}

const (
	flagDebug  = "debug"
	flagTenant = "tenant"
)

func From(ctx context.Context) *Config {
	cfg := ctx.Value(ctxKeyGlobal).(*Config)
	return cfg
}

func WithDefault(ctx context.Context, cfg *Config) context.Context {
	return context.WithValue(ctx, ctxKeyGlobal, cfg)
}

type Config struct {
	Repository Repository   `yaml:"repository,omitempty" json:"repository,omitempty"`
	Debug      bool         `yaml:"debug,omitempty"      json:"debug,omitempty"`
	Theme      styles.Theme `yaml:"theme"                json:"theme"`

	// Tenant is used to generate Microsoft Entra token, could be auto-detected via az CLI.
	// If the default tenant is not the one usable for ADO queries, users can set this value
	// in the config file, or via envAdoTenant.
	//
	// If Token is already set via envAdoPat, the Tenant value is unused.
	Tenant string `yaml:"tenant,omitempty" json:"tenant,omitempty"`

	// Token is used to authenticate to ADO, must not be logged.
	// If envAdoPat is available, this will be set to its value.
	// Otherwise, this will try to generate a Microsoft Entra token via az CLI.
	Token string `yaml:"-" json:"-"`
}

func (c Config) SetLogLevel() {
	if c.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

type Repository struct {
	Org     string `json:"org,omitempty"     yaml:"org,omitempty"`
	Project string `json:"project,omitempty" yaml:"project,omitempty"`
	Name    string `json:"name,omitempty"    yaml:"name,omitempty"`
}

func (r Repository) WebURL() string {
	return fmt.Sprintf("https://dev.azure.com/%s/%s/_git/%s", r.Org, r.Project, r.Name)
}

func AddGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(flagDebug, "d", false, "enable debug logging")
	cmd.PersistentFlags().StringP(flagTenant, "t", "", "tenant to get access token")
}

// Resolve load Config configs from these sources in this priority order:
//
//   - Built-in defaults
//   - YAML file
//   - Environment variables
//   - Command line flags
//   - Auto detect (heavy, need shell-out) for those missing values
func Resolve(cmd *cobra.Command, _ []string) error {
	// this should be the builtin config, as nothing is loaded yet.
	cfg := From(cmd.Context())

	// enable debug log as soon as possible
	cfg.SetLogLevel()

	resolvers := []func(*Config) error{
		resolveConfigFile,
		resolveEnv,
		flagsResolver(cmd),
		autoDetect,
	}

	for _, resolve := range resolvers {
		if err := resolve(cfg); err != nil {
			return err
		}
		cfg.SetLogLevel()
	}

	styles.Init(cfg.Theme)
	log.Debugf("resolved config:\n%v", styles.YAML(cfg))
	return nil
}

func flagsResolver(cmd *cobra.Command) func(cfg *Config) error {
	return func(cfg *Config) error {
		flags := cmd.Flags()
		var err error
		var allErr error

		if flags.Changed(flagDebug) {
			cfg.Debug, err = flags.GetBool(flagDebug)
			allErr = errors.Join(allErr, err)
		}

		if flags.Changed(flagTenant) {
			cfg.Tenant, err = flags.GetString(flagTenant)
			allErr = errors.Join(allErr, err)
		}

		if allErr != nil {
			log.Warnf("fail to bind flags value to config: %v", styles.YAML(err))
		}

		return allErr
	}
}

// resolveEnv binds env var with the prefix ADO_ to the config.
func resolveEnv(cfg *Config) error {
	if v, ok := os.LookupEnv(envAdoTenant); ok {
		cfg.Tenant = v
	}

	if v, ok := os.LookupEnv(envAdoDebug); ok {
		v = strings.ToLower(v)
		cfg.Debug = v != "false" && v != "0"
	}

	if v, ok := os.LookupEnv(envAdoPat); ok {
		cfg.Token = v
	}

	return nil
}
