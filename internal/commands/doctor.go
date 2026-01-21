package commands

import (
	"fmt"
	"os/exec"

	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/ui"
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
		log.Info(fmt.Sprintf("%s %s is installed", ui.IconSuccess, displayName))
		return nil
	}

	if err := checkExec("az", "Azure CLI (az)"); err != nil {
		return err
	}

	log.Info(ui.IconSuccess + " Azure CLI is available")
	log.Info("All checks passed!")
	return nil
}
