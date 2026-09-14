# Step 1 — Project Skeleton

## Objective

Create a minimal, buildable Go project for the application.

## Requirements

- Initialize a Go module.
- Use a clear module path.
- Add `cmd/chat/main.go` as the executable entry point.
- Add the following package directories:

```text
internal/googlechat
internal/chat
internal/storage
internal/ui
```

- The application must compile and run.
- Running the binary should print a short placeholder message and exit successfully.
- Add a basic `README.md`.
- Add `.gitignore` for Go build artifacts, local configuration, credentials, and SQLite files.

## Suggested commands

```bash
go mod init <module-path>
go run ./cmd/chat
go test ./...
go build ./cmd/chat
```

## Acceptance criteria

- `go test ./...` succeeds.
- `go build ./cmd/chat` succeeds.
- The repository has the expected directory structure.
- No Google API or TUI functionality is implemented yet.
