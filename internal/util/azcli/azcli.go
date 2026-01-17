// Package azcli provides utilities to extract or query data via Azure CLI.
package azcli

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/util/cache"
	"github.com/letientai299/ado/internal/util/sh"
)

// AzAdoResource is the fixed Azure DevOps resource ID, won't change.
const (
	AzAdoResource      = "499b84ac-1321-427f-aa17-267ca6975798"
	getAccessTokenBash = `
az account get-access-token --query accessToken -o tsv \
  --resource %s \
  --tenant %s`
	tokenCacheKey = "az_token"
	// Azure tokens are valid for ~1 hour, refresh 5 min before expiry
	tokenRefreshBuffer = 5 * time.Minute
)

type cachedToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Tenant    string    `json:"tenant"`
}

func (c cachedToken) isValid(tenant string) bool {
	if c.Token == "" {
		return false
	}
	if tenant != "" && c.Tenant != tenant {
		return false
	}
	return time.Now().Add(tokenRefreshBuffer).Before(c.ExpiresAt)
}

func GetToken(tenant string) (string, error) {
	// Check cache first
	var cached cachedToken
	if cache.Get(tokenCacheKey, &cached) && cached.isValid(tenant) {
		log.Debug("using cached Azure token")
		return cached.Token, nil
	}

	// Fetch new token
	token, err := fetchToken(tenant)
	if err != nil {
		return "", err
	}

	// Cache for ~1 hour (Azure default token lifetime)
	cached = cachedToken{
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Tenant:    tenant,
	}
	if err := cache.Set(tokenCacheKey, cached); err != nil {
		log.Warnf("failed to cache token: %v", err)
	}

	return token, nil
}

func fetchToken(tenant string) (string, error) {
	script := fmt.Sprintf(getAccessTokenBash, AzAdoResource, tenant)
	run := sh.Bash
	if runtime.GOOS == "windows" {
		script = strings.ReplaceAll(script, "\\", "`")
		run = sh.Pwsh
	}

	stdout, err := run(script)
	if err != nil {
		log.Errorf("fail to get token: %v", err)
		return "", err
	}

	return stdout, err
}
