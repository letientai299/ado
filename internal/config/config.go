package config

import (
	"github.com/spf13/cobra"
)

// AzAdoResource is the fixed Azure DevOps resource ID, won't change.
const AzAdoResource = "499b84ac-1321-427f-aa17-267ca6975798"

type ctxKey string

const (
	ctxKeyGlobal ctxKey = "global"
	ctxKeyToken  ctxKey = "token"
)

const (
	EnvAdoTenantID = "ADO_TENANT_ID"
)

type Global struct {
	Repo
	Debug bool `json:"debug,omitempty"`
}

type Repo struct {
	TenantID string `json:"tenant_id,omitempty"`
	Name     string `json:"name,omitempty"`
	Org      string `json:"org,omitempty"`
	Project  string `json:"project,omitempty"`
}

func AddGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool("debug", false, "enable debug logging")
}

func Resolve(cmd *cobra.Command, args []string) error {
	return nil
}
