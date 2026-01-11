package main

import (
	"context"
	"os"

	"github.com/charmbracelet/log"
	_ "github.com/joho/godotenv/autoload"
	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/pipeline"
	"github.com/letientai299/ado/internal/pull_request"
	"github.com/letientai299/ado/internal/util"
	"github.com/urfave/cli/v3"
)

func main() {
	log.SetLevel(log.DebugLevel)
	tenantId := os.Getenv(config.EnvAdoTenantID)
	token, err := util.GetToken(tenantId)
	if err != nil {
		log.Fatal(err)
	}

	cmd := &cli.Command{
		Name:  "ado",
		Usage: "Azure DevOps CLI",
		Commands: []*cli.Command{
			pull_request.Cmd,
			pipeline.Cmd,
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			return context.WithValue(ctx, "token", token), nil
		},
	}

	if err = cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
