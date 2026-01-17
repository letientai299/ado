package util

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

type StrErr string

func (s StrErr) Error() string { return string(s) }

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

type BashFunc func(script string) (string, error)

// Bash executes the given script and returns stdout, and exit error.
// In case of error, it logs the full content of stdout and stderr.
func Bash(script string) (stdout string, err error) {
	script = strings.TrimSpace(script)
	log.Debugf("executing bash script:\n%s", styles.Indent(2, script))
	cmd := exec.Command("bash", "-c", script)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	stdout = outBuf.String()
	stderr := errBuf.String()
	if err != nil {
		log.Errorf("error executing bash script: %s", err)
		log.Errorf("stdout:\n%s", stdout)
		log.Errorf("stderr:\n%s", stderr)
	}
	return strings.TrimRightFunc(stdout, unicode.IsSpace), err
}

func Ptr[T any](v T) *T { return &v }
