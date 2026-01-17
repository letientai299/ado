// Package sh provides functions for executing shell scripts and system commands.
package sh

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"unicode"

	"github.com/charmbracelet/log"
	"github.com/letientai299/ado/internal/styles"
)

// ScriptRunner is a function type that executes a shell script and returns its output.
type ScriptRunner func(script string) (string, error)

// Run uses Pwsh or Bash to execute the script depending on the current platform.
// For complex scripts, use the specific shell function instead.
func Run(script string) (string, error) {
	if runtime.GOOS == "windows" {
		return Pwsh(script)
	}
	return Bash(script)
}

// Bash executes the given script using bash and returns stdout and any error.
// Uses --norc and --noprofile for faster startup by skipping rc/profile files.
// If the execution fails, it logs stdout and stderr for debugging.
func Bash(script string) (stdout string, err error) {
	return execShell("bash", "--norc", "--noprofile", "-c", script)
}

// Pwsh executes the given script using PowerShell (pwsh) and returns stdout and any error.
// Uses -NoProfile, -NoLogo, -NonInteractive for faster startup and non-interactive execution.
// If the execution fails, it logs stdout and stderr for debugging.
func Pwsh(script string) (stdout string, err error) {
	return execShell("pwsh", "-NoProfile", "-NoLogo", "-NonInteractive", "-Command", script)
}

func execShell(shell string, args ...string) (stdout string, err error) {
	script := strings.TrimSpace(args[len(args)-1])
	log.Debugf("executing %s script:\n%s", shell, styles.Indent(2, script))
	cmd := exec.Command(shell, args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	stdout = outBuf.String()
	stderr := errBuf.String()
	if err != nil {
		log.Errorf("error executing %s script: %s", shell, err)
		log.Errorf("stdout:\n%s", stdout)
		log.Errorf("stderr:\n%s", stderr)
	}
	return strings.TrimRightFunc(stdout, unicode.IsSpace), err
}

// Browse uses system tools to open a URL in the default browser.
func Browse(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("cmd", "/c", "start", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	if err != nil {
		log.Fatal(err)
	}
}
