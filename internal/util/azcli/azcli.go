// Package azcli provides utilities to extract or query data  Azure CLI.
package azcli

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/util/sh"
)

// AzAdoResource is the fixed Azure DevOps resource ID, won't change.
const (
	AzAdoResource      = "499b84ac-1321-427f-aa17-267ca6975798"
	getAccessTokenBash = `
az account get-access-token --query accessToken -o tsv \
  --resource %s \
  --tenant %s`
)

func GetToken(tenant string) (string, error) {
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
