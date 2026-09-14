package config

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	config, err := load(func(string) (string, bool) { return "", false }, func() (string, error) {
		return "/tmp/gchat-tui", nil
	})
	if err != nil {
		t.Fatalf("load defaults: %v", err)
	}

	if config.DataDir != "/tmp/gchat-tui" {
		t.Errorf("data directory = %q, want %q", config.DataDir, "/tmp/gchat-tui")
	}
	if config.OAuthClientConfigPath != filepath.Join(config.DataDir, "client_credentials.json") {
		t.Errorf("unexpected OAuth client config path: %q", config.OAuthClientConfigPath)
	}
	if config.DatabasePath != filepath.Join(config.DataDir, "gchat-tui.db") {
		t.Errorf("unexpected database path: %q", config.DatabasePath)
	}
	if config.LogPath != filepath.Join(config.DataDir, "gchat-tui.log") {
		t.Errorf("unexpected log path: %q", config.LogPath)
	}
	if config.LogLevel != "info" {
		t.Errorf("log level = %q, want info", config.LogLevel)
	}
}

func TestLoadRejectsInvalidLogLevel(t *testing.T) {
	_, err := load(func(key string) (string, bool) {
		if key == "GCHAT_TUI_LOG_LEVEL" {
			return "verbose", true
		}
		return "", false
	}, func() (string, error) {
		return "/tmp/gchat-tui", nil
	})
	if err == nil || !strings.Contains(err.Error(), "GCHAT_TUI_LOG_LEVEL") {
		t.Fatalf("expected a clear log-level error, got %v", err)
	}
}
