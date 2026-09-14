# Step 3 — Configuration and Logging

## Objective

Add configuration and diagnostic logging without cluttering the TUI.

## Requirements

Implement configuration for:

- Google OAuth client configuration location.
- Local data directory.
- SQLite database path.
- Log file path.
- Optional log level.

Use sensible OS-specific directories where practical. Do not hardcode personal paths.

Implement logging that:

- Writes to a file rather than the terminal UI.
- Supports useful levels such as debug, info, warn, and error.
- Does not log access tokens, refresh tokens, or other secrets.

Add a configuration package if useful, while keeping responsibilities clear.

## Acceptance criteria

- The application can load defaults without requiring configuration.
- Invalid configuration produces a clear error.
- Logs are written outside the TUI.
- Secrets are not written to logs.
- Tests cover default configuration and invalid configuration.
