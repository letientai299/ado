package editor

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/log"
)

func New(tmpFilePattern, cmd string) Editor {
	return Editor{
		tmpFilePattern: tmpFilePattern,
		cmd:            cmd,
	}
}

type Editor struct {
	tmpFilePattern string
	cmd            string
}

func (e Editor) Edit(original string) (string, error) {
	tmpFile, err := os.CreateTemp(os.TempDir(), e.tmpFilePattern)
	if err != nil {
		return "", err
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	log.Debugf("editing %s", tmpFile.Name())
	if _, err = tmpFile.WriteString(original); err != nil {
		log.Errorf("fail to write to temp file: %v", err)
		_ = tmpFile.Close()
		return "", err
	}
	if err = tmpFile.Close(); err != nil {
		log.Errorf("fail to close temp file: %v", err)
		return "", err
	}

	shell := "sh"
	args := []string{"-c", e.cmd + ` "$1"`, "--", tmpFile.Name()}
	if runtime.GOOS == "windows" {
		shell = "pwsh"
		args = []string{"-NoProfile", "-Command", e.cmd + " $args[0]", tmpFile.Name()}
	}

	if err := e.run(shell, args...); err != nil {
		log.Errorf("fail to run editor: %v", err)
		return "", err
	}

	updated, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		log.Errorf("fail to read temp file: %v", err)
		return "", err
	}

	return string(updated), nil
}

func (e Editor) run(shell string, args ...string) error {
	cmd := exec.Command(shell, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
