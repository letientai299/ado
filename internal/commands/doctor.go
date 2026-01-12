package commands

import (
	"fmt"
	"os/exec"

	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/util"
	"github.com/spf13/cobra"
)

func Doctor() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Run prerequisite checks for ado",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doctorCheck()
		},
	}
}

// doctorCheck checks if git, az CLI are available and whether az authenticated.
func doctorCheck() error {
	// TODO (tai): refactor this to use lipgloss for styling success and failure check.
	log.Info("Checking prerequisites...")

	if _, err := exec.LookPath("git"); err != nil {
		log.Error("git is not installed or not in PATH")
		return fmt.Errorf("git not found: %w", err)
	}
	log.Info("✓ git is installed")

	if _, err := exec.LookPath("az"); err != nil {
		log.Error("Azure CLI (az) is not installed or not in PATH")
		return fmt.Errorf("az CLI not found: %w", err)
	}
	log.Info("✓ Azure CLI (az) is installed")

	log.Info("Checking Azure CLI authentication...")
	_, err := util.Bash("az account show")
	if err != nil {
		log.Error("Azure CLI is not authenticated. Please run 'az login'.")
		return fmt.Errorf("az not authenticated: %w", err)
	}
	log.Info("✓ Azure CLI is authenticated")

	log.Info("All checks passed!")
	return nil
}
