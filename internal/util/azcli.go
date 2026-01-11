package util

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/config"
)

func GetToken(tenantID string) (string, error) {
	const getAccessToken = `
    az account get-access-token --query accessToken -o tsv \
      --resource %s \
      --tenant %s`
	script := fmt.Sprintf(getAccessToken, config.AzAdoResource, tenantID)
	stdout, _, err := RunBash(script)
	if err != nil {
		log.Error("fail to get token: %v", err)
		return "", err
	}

	return stdout, err
}
