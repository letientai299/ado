package main

import (
	"context"
	"os"
	"slices"

	"github.com/charmbracelet/log"
	_ "github.com/joho/godotenv/autoload"
	"github.com/letientai299/ado/internal/commands"
	"github.com/letientai299/ado/internal/config"
	"github.com/letientai299/ado/internal/styles"
	"github.com/muesli/termenv"
)

func main() {
	useColor := styles.UseColor()
	if !useColor {
		log.SetColorProfile(termenv.Ascii)
	}

	log.SetReportCaller(true)
	ctx := config.WithDefault(context.Background(), newConfig(useColor))
	if err := commands.Root().ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

func newConfig(useColor bool) *config.Config {
	return &config.Config{
		Debug: isDebugEnabled(),
		Theme: chooseStyle(useColor),
	}
}

func chooseStyle(useColor bool) styles.Theme {
	if !useColor {
		return styles.NoTTy
	}

	if !termenv.HasDarkBackground() {
		return styles.Light
	}

	return styles.TokyoNight
}

func isDebugEnabled() bool {
	return os.Getenv("ADO_DEBUG") != "" ||
		slices.ContainsFunc(os.Args, func(s string) bool {
			return s == "-d" || s == "--debug"
		})
}
