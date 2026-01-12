package main

import (
	"os"

	"github.com/charmbracelet/log"
	_ "github.com/joho/godotenv/autoload"
	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/pipeline"
	"github.com/letientai299/ado/internal/pull_request"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   os.Args[0],
		Short: "Azure DevOps CLI",
		PersistentPreRunE: config.Resolve,
	}

	config.AddGlobalFlags(rootCmd)

	rootCmd.AddCommand(
		pull_request.Cmd,
		pipeline.Cmd,
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
