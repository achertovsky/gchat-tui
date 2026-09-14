# Step 4 — Google OAuth 2.0 Authentication

## Objective

Implement Google authentication for the desktop/CLI application.

## Requirements

- Use Google's OAuth 2.0 flow appropriate for an installed/native application.
- Use the official Google authentication libraries where practical.
- Store tokens securely using OS credential storage where available.
- Support token refresh.
- Support logout by deleting locally stored credentials.
- Never print tokens.
- Keep authentication code independent of the TUI.

Create an authentication interface that the rest of the application can use without knowing OAuth implementation details.

Before coding, verify:

- The required Google Cloud project configuration.
- The Google Chat API enablement requirement.
- Current OAuth scopes required for reading spaces, reading messages, sending messages, and any search functionality.

## Acceptance criteria

- A user can authenticate interactively.
- A valid access token can be obtained.
- Expired access tokens are refreshed.
- Logout removes locally stored credentials.
- Authentication errors are returned clearly.
- Tests cover non-network logic.
