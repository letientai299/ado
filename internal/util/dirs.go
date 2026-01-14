package util

func GitRoot() (string, error) {
	return Bash("git rev-parse --show-toplevel")
}
