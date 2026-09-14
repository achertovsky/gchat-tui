// Package config loads application paths and diagnostic logging settings.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const appName = "gchat-tui"

// Config contains local application settings. It intentionally excludes OAuth
// credentials and tokens, which are loaded by the authentication layer.
type Config struct {
	OAuthClientConfigPath string
	DataDir               string
	DatabasePath          string
	LogPath               string
	LogLevel              string
}

// Load reads configuration from environment variables and fills in OS-specific
// defaults when they are not set.
func Load() (Config, error) {
	return load(os.LookupEnv, defaultDataDir)
}

func load(lookupEnv func(string) (string, bool), dataDir func() (string, error)) (Config, error) {
	defaultDir, err := dataDir()
	if err != nil {
		return Config{}, fmt.Errorf("determine data directory: %w", err)
	}

	config := Config{
		DataDir:               valueOrDefault(lookupEnv, "GCHAT_TUI_DATA_DIR", defaultDir),
		OAuthClientConfigPath: valueOrDefault(lookupEnv, "GCHAT_TUI_OAUTH_CLIENT_CONFIG", filepath.Join(defaultDir, "client_credentials.json")),
		DatabasePath:          valueOrDefault(lookupEnv, "GCHAT_TUI_DATABASE_PATH", filepath.Join(defaultDir, "gchat-tui.db")),
		LogPath:               valueOrDefault(lookupEnv, "GCHAT_TUI_LOG_PATH", filepath.Join(defaultDir, "gchat-tui.log")),
		LogLevel:              valueOrDefault(lookupEnv, "GCHAT_TUI_LOG_LEVEL", "info"),
	}
	if _, ok := lookupEnv("GCHAT_TUI_DATA_DIR"); ok {
		config.OAuthClientConfigPath = valueOrDefault(lookupEnv, "GCHAT_TUI_OAUTH_CLIENT_CONFIG", filepath.Join(config.DataDir, "client_credentials.json"))
		config.DatabasePath = valueOrDefault(lookupEnv, "GCHAT_TUI_DATABASE_PATH", filepath.Join(config.DataDir, "gchat-tui.db"))
		config.LogPath = valueOrDefault(lookupEnv, "GCHAT_TUI_LOG_PATH", filepath.Join(config.DataDir, "gchat-tui.log"))
	}
	config.LogLevel = strings.ToLower(config.LogLevel)

	return config, config.validate()
}

func defaultDataDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, appName), nil
}

func valueOrDefault(lookupEnv func(string) (string, bool), key, fallback string) string {
	if value, ok := lookupEnv(key); ok {
		return value
	}
	return fallback
}

func (c Config) validate() error {
	if strings.TrimSpace(c.DataDir) == "" {
		return fmt.Errorf("GCHAT_TUI_DATA_DIR must not be empty")
	}
	for _, setting := range []struct {
		name  string
		value string
	}{
		{"GCHAT_TUI_OAUTH_CLIENT_CONFIG", c.OAuthClientConfigPath},
		{"GCHAT_TUI_DATABASE_PATH", c.DatabasePath},
		{"GCHAT_TUI_LOG_PATH", c.LogPath},
	} {
		if strings.TrimSpace(setting.value) == "" {
			return fmt.Errorf("%s must not be empty", setting.name)
		}
	}

	switch c.LogLevel {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("GCHAT_TUI_LOG_LEVEL must be debug, info, warn, or error")
	}
	return nil
}
