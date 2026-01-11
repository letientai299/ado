package pull_request

import (
	"context"

	"github.com/urfave/cli/v3"
)

var Cmd = &cli.Command{
	Name:     "pull-request",
	Aliases:  []string{"pr", "pull"},
	Usage:    "list, view, create or manipulate pull requests",
	Commands: []*cli.Command{prList, prCreate, prUpdate, prBrowse},
}

var prBrowse = &cli.Command{
	Name:    "browse",
	Aliases: []string{"open", "o"},
	Usage:   "browse a pull request in the web",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	},
}

var prCreate = &cli.Command{
	Name:    "create",
	Aliases: []string{"c"},
	Usage:   "create a pull request",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	},
}

var prUpdate = &cli.Command{
	Name:    "update",
	Aliases: []string{"u"},
	Usage:   "update a pull request",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		return nil
	},
}
