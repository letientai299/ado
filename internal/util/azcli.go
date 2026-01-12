package util

import (
	"fmt"

	"github.com/charmbracelet/log"
)

// AzAdoResource is the fixed Azure DevOps resource ID, won't change.
const AzAdoResource = "499b84ac-1321-427f-aa17-267ca6975798"

func GetToken(tenantID string) (string, error) {
	const getAccessToken = `
    az account get-access-token --query accessToken -o tsv \
      --resource %s \
      --tenant %s`
	script := fmt.Sprintf(getAccessToken, AzAdoResource, tenantID)
	stdout, err := Bash(script)
	if err != nil {
		log.Errorf("fail to get token: %v", err)
		return "", err
	}

	return stdout, err
}
