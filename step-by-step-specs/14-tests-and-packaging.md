# Step 14 — Tests and Packaging

## Objective

Prepare the application for reliable personal use and distribution.

## Requirements

Add tests for:

- Configuration.
- Authentication helpers.
- API client requests and pagination.
- API error handling.
- SQLite migrations and queries.
- Conversation and message state transitions.
- Search behavior.
- Unread-state behavior where deterministic.

Use fake API clients and fake storage in application-level tests.

Add:

- `go test ./...`
- `go vet ./...`
- Formatting checks.
- A build script or Makefile.
- README setup instructions.
- Google Cloud OAuth setup instructions.
- Configuration and data-directory documentation.
- Build instructions for Linux, macOS, and Windows where practical.

Do not include credentials in the repository.

## Acceptance criteria

- Tests pass.
- Static checks pass.
- A clean checkout can be built by following the README.
- The README explains setup, authentication, permissions, limitations, and keyboard controls.
- The resulting executable can be launched from a terminal.
