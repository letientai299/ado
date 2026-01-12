package pull_request

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/goccy/go-json"
	"github.com/letientai299/ado/internal/models"
	"github.com/letientai299/ado/internal/util"
	"github.com/spf13/cobra"
)

var prList = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List pull requests in the repo",
	RunE: func(cmd *cobra.Command, args []string) error {
		// ctx := cmd.Context()
		// token := ctx.Value("token").(string)
		// return List(ctx, token)
		return nil
	},
}

func List(ctx context.Context, token string) error {
	// calling
	// "https://dev.azure.com/skype/ES/_apis/git/repositories/4b7a7c70-ae1c-4f06-8111-eb100a5afada/pullrequests?api-version=7.1"
	// url :=
	// "https://dev.azure.com/skype/ES/_apis/git/repositories/4b7a7c70-ae1c-4f06-8111-eb100a5afada/pullrequests?api-version=7.1"
	// url :=
	// "https://skype.visualstudio.com/ES/_apis/git/repositories/4b7a7c70-ae1c-4f06-8111-eb100a5afada/pullrequests?api-version=7.1"
	url := "https://skype.visualstudio.com/ES/_apis/git/repositories/infra_regional_resilience/pullrequests?api-version=7.1"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	if token != "" {
		token = strings.TrimSpace(token)
		req.Header.Add("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var data struct {
		Value []models.GitPullRequest `json:"value"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	util.DumpJSON(data)
	return nil
}
