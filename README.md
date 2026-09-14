# gchat-tui

`gchat-tui` is a native Go terminal interface for Google Chat.

The project is currently a buildable skeleton. Google Chat integration and the
terminal interface will be introduced in later implementation steps.

## Notice

That gonna be completely llm-generated whatever with having only that line
written by hand.

## Run

```bash
go run ./cmd/chat
```

## Build and test

```bash
go test ./...
go build ./cmd/chat
```

## Configuration

All settings are optional. The default data directory uses the operating
system's user configuration location. Override individual settings with:

| Variable | Purpose |
| --- | --- |
| `GCHAT_TUI_DATA_DIR` | Local application data directory |
| `GCHAT_TUI_OAUTH_CLIENT_CONFIG` | OAuth client configuration file |
| `GCHAT_TUI_DATABASE_PATH` | SQLite database file |
| `GCHAT_TUI_LOG_PATH` | Diagnostic log file |
| `GCHAT_TUI_LOG_LEVEL` | `debug`, `info`, `warn`, or `error` |
