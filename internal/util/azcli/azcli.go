// Package azcli provides utilities to extract or query data via Azure CLI.
package azcli

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/goccy/go-json"
	"github.com/letientai299/ado/internal/util/cache"
	"github.com/letientai299/ado/internal/util/sh"
)

// AzAdoResource is the fixed Azure DevOps resource ID, won't change.
const (
	AzAdoResource      = "499b84ac-1321-427f-aa17-267ca6975798"
	getAccessTokenBash = `
az account get-access-token --query '{token:accessToken,expiry_unix:expires_on}' -o json \
  --resource %s \
  --tenant %s`
	tokenCacheKey = "az_token"
)

type accessToken struct {
	Token      string    `json:"token"`
	Expiry     time.Time `json:"expiry"`
	ExpiryUnix int64     `json:"expiry_unix"`
	Tenant     string    `json:"tenant"`
}

func (c accessToken) isValid(tenant string) bool {
	if c.Token == "" {
		return false
	}
	if tenant != "" && c.Tenant != tenant {
		return false
	}
	return time.Now().Before(c.Expiry)
}

func GetToken(tenant string) (string, error) {
	// Check cache first
	if cached, ok := cache.Get[accessToken](tokenCacheKey); ok && cached.isValid(tenant) {
		log.Debug("using cached Azure token")
		return cached.Token, nil
	}

	// Fetch new token
	token, err := fetchToken(tenant)
	if err != nil {
		return "", err
	}

	// Cache for ~1 hour (Azure default token lifetime)
	if err = cache.Set(tokenCacheKey, token); err != nil {
		log.Warnf("failed to cache token: %v", err)
	}

	return token.Token, nil
}

func fetchToken(tenant string) (*accessToken, error) {
	script := fmt.Sprintf(getAccessTokenBash, AzAdoResource, tenant)
	run := sh.Bash
	if runtime.GOOS == "windows" {
		script = strings.ReplaceAll(script, "\\", "`")
		run = sh.Pwsh
	}

	stdout, err := run(script)
	if err != nil {
		log.Errorf("fail to get token: %v", err)
		return nil, err
	}

	token := &accessToken{}
	if err = json.Unmarshal([]byte(stdout), token); err != nil {
		return nil, err
	}

	token.Tenant = tenant
	token.Expiry = time.Unix(token.ExpiryUnix, 0)
	return token, err
}
