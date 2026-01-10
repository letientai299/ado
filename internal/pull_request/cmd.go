package pull_request

import "github.com/urfave/cli/v3"

var Cmd = &cli.Command{
	Name:     "pull-request",
	Aliases:  []string{"pr", "pull"},
	Usage:    "list, view, create or manipulate pull requests",
	Commands: []*cli.Command{prList, prCreate, prUpdate, prBrowse},
}

var prList = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "list pull requests in the repo",
}

var prBrowse = &cli.Command{
	Name:    "browse",
	Aliases: []string{"open", "o"},
	Usage:   "browse a pull request in the web",
}

var prCreate = &cli.Command{
	Name:    "create",
	Aliases: []string{"c"},
	Usage:   "create a pull request",
}

var prUpdate = &cli.Command{
	Name:    "update",
	Aliases: []string{"u"},
	Usage:   "update a pull request",
}
