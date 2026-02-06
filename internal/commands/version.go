package commands

import (
	_ "embed"
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

//go:embed version.md
var versionDoc string

func Version() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  versionDoc,
		Run: func(cmd *cobra.Command, args []string) {
			info, ok := debug.ReadBuildInfo()
			if !ok {
				fmt.Println("No build info available")
				return
			}

			printInfo := func(key, value string) {
				fmt.Printf("%-12s %s\n", key, value)
			}

			fmt.Println("ado - Azure DevOps CLI")
			fmt.Println("https://github.com/letientai299/ado")
			printInfo("Version:", info.Main.Version)
			for _, setting := range info.Settings {
				switch setting.Key {
				case "vcs.revision":
					printInfo("Revision:", setting.Value)
				case "vcs.time":
					printInfo("Build Time:", setting.Value)
				}
			}
		},
	}
}
