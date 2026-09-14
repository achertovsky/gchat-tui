package main

import (
	"fmt"
	"os"

	"github.com/achertovsky/gchat-tui/internal/config"
	"github.com/achertovsky/gchat-tui/internal/logging"
	"github.com/achertovsky/gchat-tui/internal/ui"
)

func main() {
	config, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "configuration error:", err)
		os.Exit(1)
	}
	logger, err := logging.New(config.LogPath, config.LogLevel)
	if err != nil {
		fmt.Fprintln(os.Stderr, "logging error:", err)
		os.Exit(1)
	}
	logger.Info("starting terminal UI")

	if err := ui.Run(); err != nil {
		logger.Error("terminal UI exited with an error")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
