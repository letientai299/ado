package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/muesli/termenv"
	"github.com/spf13/pflag"
)

var (
	out       = termenv.NewOutput(os.Stdout, termenv.WithProfile(termenv.TrueColor))
	lineRegex = regexp.MustCompile(`^(\s*)(\w+):\s*(.*)$`)
)

func main() {
	watch := pflag.BoolP("watch", "w", false, "watch file for changes and re-render")
	pflag.Parse()

	args := pflag.Args()
	if len(args) < 1 {
		log.Fatal("usage: theme_review [-w] <yaml-file>")
	}

	filePath := args[0]
	renderFile(filePath)

	if *watch {
		watchFile(filePath)
	}
}

func renderFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(processLine(line))
	}

	if err = scanner.Err(); err != nil {
		log.Fatalf("error reading file: %v", err)
	}
	_ = file.Close()
}

func watchFile(filePath string) {
	info, err := os.Stat(filePath)
	if err != nil {
		log.Fatalf("error stat file: %v", err)
	}
	lastMod := info.ModTime()

	for {
		time.Sleep(500 * time.Millisecond)

		info, err = os.Stat(filePath)
		if err != nil {
			continue
		}

		if info.ModTime().After(lastMod) {
			lastMod = info.ModTime()
			out.ClearScreen()
			out.MoveCursor(1, 1)
			renderFile(filePath)
		}
	}
}

func processLine(line string) string {
	matches := lineRegex.FindStringSubmatch(line)
	if matches == nil {
		return line
	}

	indent := matches[1]
	key := matches[2]
	value := strings.TrimSpace(matches[3])

	// Handle quoted values
	unquoted := unquote(value)

	// Skip empty values or non-color values
	if unquoted == "" || !isColor(unquoted) {
		return line
	}

	// Apply color based on whether it's a background field
	var colored string
	if strings.Contains(key, "_background") || key == "background" {
		colored = applyBackground(value, unquoted)
	} else {
		colored = applyForeground(value, unquoted)
	}

	return fmt.Sprintf("%s%s: %s", indent, key, colored)
}

func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func isColor(s string) bool {
	if s == "" {
		return false
	}

	// Hex color: #RGB, #RRGGBB
	if s[0] == '#' {
		hex := s[1:]
		if len(hex) == 3 || len(hex) == 6 {
			for _, c := range hex {
				if !isHexDigit(c) {
					return false
				}
			}
			return true
		}
		return false
	}

	// ANSI 16/256 colors: numeric values
	if isNumeric(s) {
		return true
	}

	// Named ANSI colors
	named := map[string]bool{
		"black":          true,
		"red":            true,
		"green":          true,
		"yellow":         true,
		"blue":           true,
		"magenta":        true,
		"cyan":           true,
		"white":          true,
		"bright_black":   true,
		"bright_red":     true,
		"bright_green":   true,
		"bright_yellow":  true,
		"bright_blue":    true,
		"bright_magenta": true,
		"bright_cyan":    true,
		"bright_white":   true,
	}
	return named[strings.ToLower(s)]
}

func isHexDigit(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func applyForeground(display, color string) string {
	return out.String(display).Foreground(out.Color(color)).String()
}

func applyBackground(display, color string) string {
	return out.String(display).Background(out.Color(color)).String()
}
