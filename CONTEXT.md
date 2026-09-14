# Project Context

`gchat-tui` is a Go 1.19 Bubble Tea client for Google Chat. Follow
`step-by-step-specs/00-overview.md` and implement only one numbered spec step
at a time. Do not write unit tests; validate manually.

## Current state

- Steps 1-8 are implemented.
- Step 9 is partially implemented: `chat.users.readstate` support is present,
  unread conversations are marked with `*` and sorted first.
- Required next action: `./chat login --reauthorize` and accept the
  `chat.users.readstate` scope.
- Read state uses Google Chat's authoritative `lastReadTime`, compared with
  the newest message timestamp.

## Commands

```bash
go build ./cmd/chat
./chat
./chat login
./chat login --reauthorize
./chat logout
```

## OAuth and API

- OAuth client JSON: `GCHAT_TUI_OAUTH_CLIENT_CONFIG`, defaulting to
  `~/.config/gchat-tui/client_credentials.json`.
- Configure the Google Chat app and OAuth consent screen in the same Cloud
  project. Required scopes are defined in `internal/auth/auth.go`.
- Pending direct-message invitations cannot be accepted through the Google
  Chat API. The app should prevent sending and direct the user to
  `https://chat.google.com/`.

## Structure

- `cmd/chat`: command wiring.
- `internal/auth`: OAuth, token storage, identity.
- `internal/googlechat`: Google Chat API wrapper.
- `internal/ui`: Bubble Tea UI; keep it API-independent.
- `internal/storage`: reserved for the later SQLite-cache step.
