package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/achertovsky/gchat-tui/internal/auth"
	"github.com/achertovsky/gchat-tui/internal/config"
	"github.com/achertovsky/gchat-tui/internal/googlechat"
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
			reauthorize := len(os.Args) == 3 && os.Args[2] == "--reauthorize"
			if len(os.Args) > 2 && !reauthorize {
				fmt.Fprintln(os.Stderr, "usage: chat login [--reauthorize]")
				os.Exit(2)
			}
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
			if err != nil {
				logger.Error("OAuth client configuration is invalid")
				fmt.Fprintln(os.Stderr, "authentication error:", err)
				os.Exit(1)
			}
			fmt.Fprintln(os.Stderr, "Checking existing credentials...")
			if !reauthorize {
				if _, err := provider.Token(context.Background()); err == nil {
					fmt.Fprintln(os.Stderr, "You are already authenticated.")
					logger.Info("OAuth login skipped because valid credentials already exist")
					return
				} else if !errors.Is(err, auth.ErrNotAuthenticated) {
					logger.Warn("Stored OAuth credentials could not be used; requesting authorization again")
				}
			}
			err = provider.Login(context.Background())
			if err != nil {
				logger.Error("OAuth login failed")
				fmt.Fprintln(os.Stderr, "authentication error:", err)
				os.Exit(1)
			}
			logger.Info("OAuth login completed")
			fmt.Fprintln(os.Stderr, "Authentication successful. OAuth credentials have been saved.")
			return
		case "logout":
			fmt.Fprintln(os.Stderr, "Removing stored credentials...")
			if err := auth.Logout(config.DataDir); err != nil {
				logger.Error("OAuth logout failed")
				fmt.Fprintln(os.Stderr, "logout error:", err)
				os.Exit(1)
			}
			logger.Info("OAuth credentials removed")
			return
		default:
			fmt.Fprintln(os.Stderr, "usage: chat [login [--reauthorize]|logout]")
			os.Exit(2)
		}
	}

	provider, authErr := auth.New(config.OAuthClientConfigPath, config.DataDir, nil)
	var chatClient *googlechat.Client
	if authErr == nil {
		chatClient, authErr = googlechat.New(provider.HTTPClient(context.Background()))
	}
	loadConversations := func(ctx context.Context) ([]ui.Conversation, error) {
		if authErr != nil {
			return nil, authErr
		}
		if _, err := provider.Token(ctx); err != nil {
			return nil, err
		}
		spaces, err := chatClient.ListSpaces(ctx)
		if err != nil {
			return nil, err
		}
		conversations := make([]ui.Conversation, len(spaces))
		for index, space := range spaces {
			conversations[index] = ui.Conversation{
				Name:        space.Name,
				DisplayName: space.DisplayName,
				Type:        space.Type,
			}
		}
		return conversations, nil
	}
	loadMessages := func(ctx context.Context, conversationName, pageToken string) (ui.MessagePage, error) {
		if authErr != nil {
			logger.Warn("message loading failed: authentication initialization failed")
			return ui.MessagePage{}, authErr
		}
		if _, err := provider.Token(ctx); err != nil {
			logger.Warn("message loading failed: OAuth credentials are unavailable or invalid")
			return ui.MessagePage{}, err
		}
		messages, nextPageToken, err := chatClient.ListMessagesPage(ctx, conversationName, pageToken)
		if err != nil {
			logger.Warn(fmt.Sprintf("message loading failed: %v", err))
			return ui.MessagePage{}, err
		}
		page := ui.MessagePage{
			Messages:      make([]ui.Message, len(messages)),
			NextPageToken: nextPageToken,
		}
		for index, message := range messages {
			senderName := message.Sender.DisplayName
			if senderName == "" {
				senderName = message.Sender.Name
			}
			page.Messages[index] = ui.Message{
				SenderName: senderName,
				Text:       message.Text,
				Timestamp:  message.CreateTime,
			}
		}
		return page, nil
	}
	sendMessage := func(ctx context.Context, conversation ui.Conversation, text string) (ui.Message, error) {
		if authErr != nil {
			logger.Warn("message sending failed: authentication initialization failed")
			return ui.Message{}, authErr
		}
		if _, err := provider.Token(ctx); err != nil {
			logger.Warn("message sending failed: OAuth credentials are unavailable or invalid")
			return ui.Message{}, err
		}
		if conversation.Type == "DIRECT_MESSAGE" {
			pending, err := chatClient.HasPendingInvitation(ctx, conversation.Name)
			if err != nil {
				logger.Warn(fmt.Sprintf("membership check failed: %v", err))
				return ui.Message{}, err
			}
			if pending {
				return ui.Message{}, googlechat.ErrInvitationPending
			}
		}
		message, err := chatClient.CreateMessage(ctx, conversation.Name, text)
		if err != nil {
			logger.Warn(fmt.Sprintf("message sending failed: %v", err))
			return ui.Message{}, err
		}
		senderName := message.Sender.DisplayName
		if senderName == "" {
			senderName = "You"
		}
		if message.Text == "" {
			message.Text = text
		}
		return ui.Message{
			SenderName: senderName,
			Text:       message.Text,
			Timestamp:  message.CreateTime,
		}, nil
	}
	if err := ui.Run(loadConversations, loadMessages, sendMessage); err != nil {
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
	fmt.Fprintln(os.Stderr, "3. Open Google Chat API > Configuration and configure the Chat app's name and availability.")
	fmt.Fprintln(os.Stderr, "4. Open Google Auth platform and configure the Branding, Audience, and Data Access pages.")
	fmt.Fprintln(os.Stderr, "   If the app is External and still in testing, add your Google account as a test user.")
	fmt.Fprintln(os.Stderr, "5. Open Google Auth platform > Clients, create an OAuth client ID, and select \"Desktop app\".")
	fmt.Fprintln(os.Stderr, "6. Download the OAuth client JSON file.")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "The application expected that file at:")
	fmt.Fprintln(os.Stderr, "  "+clientConfigPath)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Either copy the downloaded file there, or set its location for this shell:")
	fmt.Fprintln(os.Stderr, "  export GCHAT_TUI_OAUTH_CLIENT_CONFIG=/path/to/client_secret.json")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Then run `chat login` again. The application will open a browser for authorization.")
}
