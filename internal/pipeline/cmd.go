package pipeline

import "github.com/urfave/cli/v3"

var Cmd = &cli.Command{
	Name:     "pipeline",
	Aliases:  []string{"pp"},
	Usage:    "list, view, run pipeline",
	Commands: []*cli.Command{ppList, ppRun, ppBrowse},
}

var ppList = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "list pull requests in the repo",
}

var ppRun = &cli.Command{
	Name:    "run",
	Aliases: []string{"c"},
	Usage:   "create a pull request",
}

var ppBrowse = &cli.Command{
	Name:    "browse",
	Aliases: []string{"u"},
	Usage:   "browse recent runs of a pipeline on the web",
}
