package main

import (
	"context"
	"os"
	"slices"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/x/term"
	_ "github.com/joho/godotenv/autoload"
	"github.com/letientai299/ado/internal/commands"
	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/styles"
	"github.com/muesli/termenv"
)

func main() {
	ctx := config.WithDefault(context.Background(), newConfig())
	if err := commands.Root().ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}
}

func newConfig() *config.Config {
	useColor := os.Getenv("COLOR") == "always" ||
		(term.IsTerminal(os.Stdout.Fd()) && term.IsTerminal(os.Stderr.Fd()))

	return &config.Config{
		Debug: isDebugEnabled(),
		Theme: chooseBestStyle(useColor),
	}
}

func chooseBestStyle(useColor bool) styles.Theme {
	if !useColor {
		return styles.ThemeNoTTy
	}

	if !termenv.HasDarkBackground() {
		return styles.ThemeLight
	}

	return styles.ThemeTokyoNight
}

func isDebugEnabled() bool {
	return os.Getenv("ADO_DEBUG") != "" ||
		slices.ContainsFunc(os.Args, func(s string) bool {
			return s == "-d" || s == "--debug"
		})
}
