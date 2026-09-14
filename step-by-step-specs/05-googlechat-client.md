# Step 5 — Google Chat API Client Foundation

## Objective

Create a typed, testable client for the Google Chat REST API.

## Requirements

Implement a package under `internal/googlechat` that:

- Accepts an authenticated HTTP client.
- Wraps the official Google Chat API.
- Defines application-level types for spaces/conversations, users, and messages.
- Converts Google API response types into application-level types where useful.
- Handles pagination.
- Handles API errors without exposing secrets.
- Supports context cancellation and request timeouts.

Do not put Bubble Tea code in this package.

Initially implement only the API methods needed for:

- Listing spaces/conversations.
- Listing messages in a space.
- Creating a message.

Verify exact endpoint names, request formats, permissions, and scopes against current Google documentation.

## Acceptance criteria

- The package can be tested with a fake HTTP server.
- Pagination is tested.
- API errors are tested.
- No real API call is required for unit tests.
- The package has no dependency on the TUI.
