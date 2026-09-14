package main

import (
	"context"
	"errors"
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
				fmt.Fprintln(os.Stderr, "Opening a browser to authenticate.")
				fmt.Fprintln(os.Stderr, "If no browser opens, copy this URL into one:")
				fmt.Fprintln(os.Stderr, url)
				fmt.Fprintln(os.Stderr, "Waiting for authorization; press Ctrl+C to cancel.")
			})
			if errors.Is(err, os.ErrNotExist) {
				printOAuthSetup(config.OAuthClientConfigPath)
				os.Exit(1)
			}
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

func printOAuthSetup(clientConfigPath string) {
	fmt.Fprintln(os.Stderr, "Google OAuth setup is required before signing in.")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "1. Open https://console.cloud.google.com/ and create or select a project.")
	fmt.Fprintln(os.Stderr, "2. Open APIs & Services > Library, search for \"Google Chat API\", and enable it.")
	fmt.Fprintln(os.Stderr, "3. Open Google Auth platform and configure the Branding, Audience, and Data Access pages.")
	fmt.Fprintln(os.Stderr, "   If the app is External and still in testing, add your Google account as a test user.")
	fmt.Fprintln(os.Stderr, "4. Open Google Auth platform > Clients, create an OAuth client ID, and select \"Desktop app\".")
	fmt.Fprintln(os.Stderr, "5. Download the OAuth client JSON file.")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "The application expected that file at:")
	fmt.Fprintln(os.Stderr, "  "+clientConfigPath)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Either copy the downloaded file there, or set its location for this shell:")
	fmt.Fprintln(os.Stderr, "  export GCHAT_TUI_OAUTH_CLIENT_CONFIG=/path/to/client_secret.json")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Then run `chat login` again. The application will open a browser for authorization.")
}
