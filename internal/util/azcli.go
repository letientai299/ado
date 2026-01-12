package util

import (
	"fmt"

	"github.com/charmbracelet/log"
)

// AzAdoResource is the fixed Azure DevOps resource ID, won't change.
const (
	AzAdoResource       = "499b84ac-1321-427f-aa17-267ca6975798"
	shellGetAccessToken = `
az account get-access-token --query accessToken -o tsv \
  --resource %s \
  --tenant %s`
)

func GetToken(tenant string) (string, error) {
	script := fmt.Sprintf(shellGetAccessToken, AzAdoResource, tenant)
	stdout, err := Bash(script)
	if err != nil {
		log.Errorf("fail to get token: %v", err)
		return "", err
	}

	return stdout, err
}
