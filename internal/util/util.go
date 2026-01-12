package util

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"unicode"

	"github.com/charmbracelet/log"
	"github.com/goccy/go-json"
	"github.com/letientai299/ado/internal/styles"
)

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

// Bash executes the given script and returns stdout, and exit error.
// In case of error, it logs the full content of stdout and stderr.
func Bash(script string) (stdout string, err error) {
	script = strings.TrimSpace(script)
	log.Debugf("executing bash script:\n%s", Indent(2, script))
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

// Indent add indentation of n spaces to every line in the string
func Indent(n int, s string) string {
	padding := strings.Repeat(" ", n)
	return padding + strings.ReplaceAll(s, "\n", "\n"+padding)
}

func JSON(v any) string {
	var buf bytes.Buffer
	encodeJSON(v, json.NewEncoder(&buf))
	return buf.String()
}

// DumpJSON prints the object as prettified JSON in stdout.
func DumpJSON(v any) {
	encodeJSON(v, json.NewEncoder(os.Stdout))
}

func encodeJSON(v any, encoder *json.Encoder) {
	encoder.SetIndent("", "  ")
	var options []json.EncodeOptionFunc
	if styles.UseColor {
		options = append(options, json.Colorize(json.DefaultColorScheme))
	}

	err := encoder.EncodeWithOption(v, options...)
	if err != nil {
		log.Fatal("fail to dump json: %v, err=%v", v, err)
	}
}

// ParseRepoInfo parses the origin URL to get the organization, project, and repo name.
// It recognizes these URL formats:
//
//   - General format: https://dev.azure.com/{org}/{project}/_git/{repo}
//   - Per instance: https://{org}.{host}/{project}/_git/{repo}
//   - SSH format: git@ssh.dev.azure.com:v3/{org}/{project}/{repo}
func ParseRepoInfo(origin string) (string, string, string, error) {
	if strings.HasPrefix(origin, "git@") {
		return parseRepoInfoSSH(origin)
	}

	u, err := url.Parse(origin)
	if err != nil {
		return "", "", "", err
	}

	path := strings.TrimPrefix(u.Path, "/")
	parts := strings.Split(path, "/")

	var org, project, repo string

	// Find _git index
	gitIdx := -1
	for i, part := range parts {
		if part == "_git" {
			gitIdx = i
			break
		}
	}

	if gitIdx == -1 {
		return "", "", "", fmt.Errorf("invalid Azure DevOps url: %s", origin)
	}

	if gitIdx+1 >= len(parts) {
		return "", "", "", fmt.Errorf("invalid Azure DevOps url (missing repo): %s", origin)
	}
	repo = parts[gitIdx+1]

	if gitIdx-1 < 0 {
		return "", "", "", fmt.Errorf("invalid Azure DevOps url (missing project): %s", origin)
	}
	project = parts[gitIdx-1]

	if u.Hostname() == "dev.azure.com" {
		if len(parts) < 1 {
			return "", "", "", fmt.Errorf("invalid Azure DevOps url (missing org): %s", origin)
		}
		org = parts[0]
	} else {
		hostParts := strings.Split(u.Hostname(), ".")
		if len(hostParts) < 2 {
			return "", "", "", fmt.Errorf("invalid Azure DevOps host: %s", origin)
		}
		org = hostParts[0]
	}

	return org, project, repo, nil
}

func parseRepoInfoSSH(origin string) (string, string, string, error) {
	// SSH format: git@ssh.dev.azure.com:v3/{org}/{project}/{repo}
	parts := strings.SplitN(origin, ":", 2)
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("invalid ssh url: %s", origin)
	}
	path := parts[1]
	pathParts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(pathParts) < 4 {
		return "", "", "", fmt.Errorf("invalid ssh url path: %s", origin)
	}
	// pathParts should be ["v3", "{org}", "{project}", "{repo}"]
	return pathParts[1], pathParts[2], pathParts[3], nil
}
