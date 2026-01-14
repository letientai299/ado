package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/util"
)

func autoDetect(cfg *Config) error {
	if err := detectTenant(cfg, util.Bash); err != nil {
		return err
	}

	return detectRepo(cfg, util.Bash)
}

func detectTenant(cfg *Config, bash util.BashFunc) error {
	if cfg.Tenant != "" && cfg.Username != "" {
		return nil // skip detecting since tenant info is already set
	}

	raw, err := bash(`az account show --query "{tenantId:tenantId,username:user.name}" -o tsv`)
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

func detectRepo(cfg *Config, bash util.BashFunc) error {
	if cfg.Repository.Name != "" && cfg.Repository.Org != "" && cfg.Repository.Project != "" {
		return nil // skip detecting since repo info is already set
	}

	gitOrigin, err := bash(`git remote get-url origin`)
	if err != nil {
		log.Errorf("fail to get git origin url: %v", err)
		return err
	}

	org, project, repo, err := parseRepoInfo(gitOrigin)
	if err != nil {
		log.Errorf("fail to parse git origin url for ADO repo info: %v", err)
		return err
	}

	if cfg.Repository.Org == "" {
		cfg.Repository.Org = org
	}
	if cfg.Repository.Project == "" {
		cfg.Repository.Project = project
	}
	if cfg.Repository.Name == "" {
		cfg.Repository.Name = repo
	}

	return nil
}

// parseRepoInfo parses the origin URL to get the organization, project, and repo name.
// It recognizes these URL formats:
//
//   - General format: https://dev.azure.com/{org}/{project}/_git/{repo}
//   - Per instance: https://{org}.{host}/{project}/_git/{repo}
//   - SSH format: git@ssh.dev.azure.com:v3/{org}/{project}/{repo}
func parseRepoInfo(origin string) (string, string, string, error) {
	if strings.HasPrefix(origin, "git@") {
		return parseRepoInfoSSH(origin)
	}

	u, err := url.Parse(origin)
	if err != nil {
		return "", "", "", err
	}

	path := strings.TrimPrefix(u.Path, "/")
	parts := strings.Split(path, "/")

	var org, project, repo string

	// Find _git index
	gitIdx := -1
	for i, part := range parts {
		if part == "_git" {
			gitIdx = i
			break
		}
	}

	if gitIdx == -1 {
		return "", "", "", fmt.Errorf("invalid Azure DevOps url: %s", origin)
	}

	if gitIdx+1 >= len(parts) {
		return "", "", "", fmt.Errorf("invalid Azure DevOps url (missing repo): %s", origin)
	}
	repo = parts[gitIdx+1]

	if gitIdx-1 < 0 {
		return "", "", "", fmt.Errorf("invalid Azure DevOps url (missing project): %s", origin)
	}
	project = parts[gitIdx-1]

	if u.Hostname() == "dev.azure.com" {
		if len(parts) < 1 {
			return "", "", "", fmt.Errorf("invalid Azure DevOps url (missing org): %s", origin)
		}
		org = parts[0]
	} else {
		hostParts := strings.Split(u.Hostname(), ".")
		if len(hostParts) < 2 {
			return "", "", "", fmt.Errorf("invalid Azure DevOps host: %s", origin)
		}
		org = hostParts[0]
	}

	return org, project, repo, nil
}

func parseRepoInfoSSH(origin string) (string, string, string, error) {
	// SSH format: git@ssh.dev.azure.com:v3/{org}/{project}/{repo}
	parts := strings.SplitN(origin, ":", 2)
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("invalid ssh url: %s", origin)
	}
	path := parts[1]
	pathParts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(pathParts) < 4 {
		return "", "", "", fmt.Errorf("invalid ssh url path: %s", origin)
	}
	// pathParts should be ["v3", "{org}", "{project}", "{repo}"]
	return pathParts[1], pathParts[2], pathParts[3], nil
}
