package main

import (
	"context"

	"github.com/charmbracelet/log"
	_ "github.com/joho/godotenv/autoload"
	"github.com/letientai299/ado/internal/commands"
	"github.com/letientai299/ado/internal/config"
)

func main() {
	ctx := config.WithDefault(context.Background())
	if err := commands.Root().ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}
}
