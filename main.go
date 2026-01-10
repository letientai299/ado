package main

import (
	"context"
	"log"
	"os"

	"github.com/letientai299/ado/internal/pipeline"
	"github.com/letientai299/ado/internal/pull_request"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name: "ado",
		Commands: []*cli.Command{
			pull_request.Cmd,
			pipeline.Cmd,
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
