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

## Google OAuth setup

Before logging in, create a Google Cloud project, enable the Google Chat API,
configure the OAuth consent screen, and create a **Desktop app** OAuth client.
Download its client configuration JSON and set
`GCHAT_TUI_OAUTH_CLIENT_CONFIG` to its location.

The application requests these user-authorized Google Chat scopes:

- `chat.spaces.readonly` to list conversations
- `chat.messages.readonly` to read messages
- `chat.messages.create` to send messages
- `chat.memberships.readonly` to detect pending direct-message invitations
- `userinfo.email` to verify that the signed-in account joined a direct message

Run `go run ./cmd/chat login` to authorize the app in a browser. Use
`go run ./cmd/chat login --reauthorize` to force a new consent flow after
changing the scopes. OAuth
credentials are stored in the OS credential store when available, with a
local, owner-only fallback file when it is not. Run
`go run ./cmd/chat logout` to delete them.
