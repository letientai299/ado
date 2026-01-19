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

	if err = Open(e.cmd, tmpFile.Name()); err != nil {
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

func Open(cmd, filePath string) error {
	shell := "sh"
	args := []string{"-c", cmd + ` "$1"`, "--", filePath}
	if runtime.GOOS == "windows" {
		shell = "pwsh"
		args = []string{"-NoProfile", "-Command", cmd + " $args[0]", filePath}
	}

	x := exec.Command(shell, args...)
	x.Stdin = os.Stdin
	x.Stdout = os.Stdout
	x.Stderr = os.Stderr
	return x.Run()
}