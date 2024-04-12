package main

import (
	"log/slog"

	"github.com/infezek/app-chat/pkg/config"
	"github.com/spf13/cobra"
)

func main() {
	slog.Info("Starting app-chat")
	cfg, err := config.New("chat")
	if err != nil {
		panic(err)
	}
	cmd := &cobra.Command{}
	cmd.AddCommand(
		http(cfg),
		development(cfg),
		login(),
		Seed(),
	)
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
