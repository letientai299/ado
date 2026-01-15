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

	checkExec := func(name, displayName string) error {
		if _, err := exec.LookPath(name); err != nil {
			log.Error(displayName + " is not installed or not in PATH")
			return fmt.Errorf("%s not found: %w", name, err)
		}
		log.Info(fmt.Sprintf("✓ %s is installed", displayName))
		return nil
	}

	if err := checkExec("git", "git"); err != nil {
		return err
	}

	if err := checkExec("az", "Azure CLI (az)"); err != nil {
		return err
	}

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
