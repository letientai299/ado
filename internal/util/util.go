package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/log"
	"github.com/goccy/go-json"
	"github.com/mattn/go-isatty"
)

var useColor = isatty.IsTerminal(os.Stdout.Fd()) ||
	isatty.IsCygwinTerminal(os.Stdout.Fd()) ||
	os.Getenv("COLOR") == "always"

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

// RunBash executes the given Bash script and return stdout, stderr and exit error
func RunBash(script string) (stdout, stderr string, err error) {
	log.Debugf("executing bash script: %s", script)
	cmd := exec.Command("bash", "-c", script)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	stdout = outBuf.String()
	stderr = errBuf.String()
	if err != nil {
		log.Errorf("error executing bash script: %s", err)
		log.Errorf("stdout:\n%s", stdout)
		log.Errorf("stderr:\n%s", stderr)
	}
	return stdout, stderr, err
}

// DumpJSON prints the object as prettified JSON in stdout.
func DumpJSON(v any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	var options []json.EncodeOptionFunc
	if useColor {
		options = append(options, json.Colorize(json.DefaultColorScheme))
	}

	err := encoder.EncodeWithOption(v, options...)
	if err != nil {
		log.Errorf("fail to dump json: %v, err=%v", v, err)
		return err
	}

	return nil
}
