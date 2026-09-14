// Package logging provides file-based diagnostic logging.
package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type level int

const (
	debugLevel level = iota
	infoLevel
	warnLevel
	errorLevel
)

// Logger writes diagnostic messages to a file.
type Logger struct {
	logger    *log.Logger
	threshold level
}

// New creates a logger that writes to path. Callers must not include secrets
// such as OAuth tokens in log messages.
func New(path, configuredLevel string) (*Logger, error) {
	threshold, err := parseLevel(configuredLevel)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("create log directory: %w", err)
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}
	return &Logger{
		logger:    log.New(file, "", log.LstdFlags|log.LUTC),
		threshold: threshold,
	}, nil
}

func (l *Logger) Debug(message string) { l.log(debugLevel, "DEBUG", message) }
func (l *Logger) Info(message string)  { l.log(infoLevel, "INFO", message) }
func (l *Logger) Warn(message string)  { l.log(warnLevel, "WARN", message) }
func (l *Logger) Error(message string) { l.log(errorLevel, "ERROR", message) }

func (l *Logger) log(messageLevel level, label, message string) {
	if messageLevel >= l.threshold {
		l.logger.Printf("%s %s", label, message)
	}
}

func parseLevel(value string) (level, error) {
	switch strings.ToLower(value) {
	case "debug":
		return debugLevel, nil
	case "info":
		return infoLevel, nil
	case "warn":
		return warnLevel, nil
	case "error":
		return errorLevel, nil
	default:
		return 0, fmt.Errorf("invalid log level %q", value)
	}
}
