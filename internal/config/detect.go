package config

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/util"
)

func autoDetect(cfg *Config) error {
	if err := detectTenant(cfg); err != nil {
		return err
	}

	return detectRepo(cfg)
}

func detectTenant(cfg *Config) error {
	if cfg.Tenant != "" && cfg.Username != "" {
		return nil // skip detecting since tenant info is already set
	}

	raw, err := util.Bash(`az account show --query "{tenantId:tenantId,username:user.name}" -o tsv`)
	if err != nil {
		log.Errorf("fail to detect tenant: %v", err)
		return err
	}

	parts := strings.Split(raw, "\t")
	if cfg.Tenant == "" {
		cfg.Tenant = parts[0]
	}

	if cfg.Username == "" {
		cfg.Username = parts[1]
	}
	return nil
}

func detectRepo(cfg *Config) error {
	if cfg.Repo != "" {
		return nil // skip detecting since repo info is already set
	}

	gitOrigin, err := util.Bash(`git remote get-url origin`)
	if err != nil {
		log.Errorf("fail to get git origin url: %v", err)
		return err
	}

	cfg.Org, cfg.Project, cfg.Repo, err = util.ParseRepoInfo(gitOrigin)
	if err != nil {
		log.Errorf("fail to parse git origin url for ADO repo info: %v", err)
		return err
	}

	return nil
}
