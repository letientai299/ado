package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/letientai299/ado/internal/styles"
	"github.com/letientai299/ado/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type ctxKey string

const (
	ctxKeyGlobal ctxKey = "global"
)

var configFileNames = []string{
	".ado.yml",
	".ado.yaml",
	".config/ado.yml",
	".config/ado.yaml",
}

var koanfUnmarshalConf = koanf.UnmarshalConf{
	Tag: "yaml",
	DecoderConfig: &mapstructure.DecoderConfig{
		TagName:          "yaml",
		Squash:           true,
		WeaklyTypedInput: true,
	},
}

const (
	flagDebug  = "debug"
	flagTenant = "tenant"
)

func WithDefault(ctx context.Context, cfg *Config) context.Context {
	return context.WithValue(ctx, ctxKeyGlobal, cfg)
}

type Config struct {
	Repository Repository   `yaml:"repository,omitempty" json:"repository,omitempty"`
	Tenant     string       `yaml:"tenant,omitempty"     json:"tenant,omitempty"`
	Username   string       `yaml:"username,omitempty"   json:"username,omitempty"`
	Debug      bool         `yaml:"debug,omitempty"      json:"debug,omitempty"`
	Theme      styles.Theme `yaml:"theme"                json:"theme"`
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
	log.Debugf("resolved config: %v", styles.YAML(cfg))
	addTemplateHelpers()
	return nil
}

func addTemplateHelpers() {
	cobra.AddTemplateFunc("headingStyle", styles.HeadingStyle)
	cobra.AddTemplateFunc("flagStyle", styles.FlagStyle)
	cobra.AddTemplateFunc("cmdStyle", styles.CmdStyle)
	cobra.AddTemplateFunc("flags", flagSlice)
	cobra.AddTemplateFunc("flagName", flagName)
	cobra.AddTemplateFunc("indent", util.Indent)
	cobra.AddTemplateFunc("wrap", styles.Wrap)
}

func flagSlice(fs *pflag.FlagSet) []*pflag.Flag {
	var list []*pflag.Flag
	fs.VisitAll(func(f *pflag.Flag) {
		if !f.Hidden {
			list = append(list, f)
		}
	})
	return list
}

func flagName(f *pflag.Flag) string {
	var s string
	if f.Shorthand != "" {
		s = "-" + f.Shorthand + ", --" + f.Name
	} else {
		s = "    --" + f.Name
	}
	if f.Value.Type() != "bool" {
		s += " " + f.Value.Type()
	}
	return s
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

func From(ctx context.Context) *Config {
	cfg := ctx.Value(ctxKeyGlobal).(*Config)
	return cfg
}

// resolveConfigFile finds the YAML config file and loads it using koanf YAML
// parsers to load the file.
func resolveConfigFile(cfg *Config) error {
	filePath, err := findConfigFile()
	if err != nil {
		return err
	}

	if filePath == "" {
		return nil // no config file found
	}

	log.Debugf("found config file %v", filePath)
	k := koanf.New(".")
	parser := yaml.Parser()

	if err = k.Load(file.Provider(filePath), parser); err != nil {
		return err
	}

	if err = k.UnmarshalWithConf("", &cfg, koanfUnmarshalConf); err != nil {
		log.Fatalf("fail to parse config file: %v", err)
		return err
	}

	return nil
}

// findConfigFile looks for .ado.y(a)ml or `.config/ado.y(a)ml` in the
// working dir, then continue the search up to the git root dir.
func findConfigFile() (string, error) {
	gitRoot, err := util.GitRoot()
	if err != nil {
		log.Warnf("fail to get git root dir: %v", err)
		return "", err
	}

	wd, _ := os.Getwd()

	for {
		for _, f := range configFileNames {
			p := filepath.Join(wd, f)
			if _, err = os.Stat(p); err == nil {
				return p, nil
			}
		}

		if wd == gitRoot || wd == filepath.Dir(wd) {
			break
		}
		wd = filepath.Dir(wd)
	}

	return "", nil
}

// resolveEnv binds env var with the prefix ADO_ to the config.
func resolveEnv(cfg *Config) error {
	k := koanf.New(".")
	envProvider := env.Provider(".", env.Opt{
		Prefix: "ADO_",
		TransformFunc: func(k, v string) (string, any) {
			key := strings.ToLower(strings.TrimPrefix(k, "ADO_"))
			return strings.ToLower(key), v
		},
	})
	err := k.Load(envProvider, nil)
	if err != nil {
		log.Warnf("failed to load environment variables: %v", err)
		return err
	}

	if err = k.UnmarshalWithConf("", cfg, koanfUnmarshalConf); err != nil {
		log.Warnf("failed to unmarshal environment variables: %v", err)
		return err
	}

	return nil
}
