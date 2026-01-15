package commands

import (
	"fmt"
	"runtime/debug"

	"github.com/letientai299/ado/internal/styles"
	"github.com/spf13/cobra"
)

func Version() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			info, ok := debug.ReadBuildInfo()
			if !ok {
				fmt.Println("No build info available")
				return
			}

			printInfo := func(key, value string) {
				fmt.Printf(
					"%s %s\n",
					styles.HeadingStyle(fmt.Sprintf("%-12s", key)),
					value,
				)
			}

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
