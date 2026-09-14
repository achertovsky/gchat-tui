package main

import (
	"context"
	"fmt"
	"os"

	"github.com/achertovsky/gchat-tui/internal/auth"
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

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "login":
			provider, err := auth.New(config.OAuthClientConfigPath, config.DataDir, func(url string) {
				fmt.Fprintln(os.Stderr, "Open this URL in a browser to authenticate:")
				fmt.Fprintln(os.Stderr, url)
			})
			if err == nil {
				err = provider.Login(context.Background())
			}
			if err != nil {
				logger.Error("OAuth login failed")
				fmt.Fprintln(os.Stderr, "authentication error:", err)
				os.Exit(1)
			}
			logger.Info("OAuth login completed")
			return
		case "logout":
			if err := auth.Logout(config.DataDir); err != nil {
				logger.Error("OAuth logout failed")
				fmt.Fprintln(os.Stderr, "logout error:", err)
				os.Exit(1)
			}
			logger.Info("OAuth credentials removed")
			return
		default:
			fmt.Fprintln(os.Stderr, "usage: chat [login|logout]")
			os.Exit(2)
		}
	}

	if err := ui.Run(); err != nil {
		logger.Error("terminal UI exited with an error")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
