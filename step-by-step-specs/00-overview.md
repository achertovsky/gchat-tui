# Google Chat CLI — Implementation Plan

## Goal

Build a native Go terminal user interface for Google Chat.

The application will support:

- Listing channels/spaces and direct messages.
- Showing unread conversations.
- Opening and reading conversations.
- Sending text messages.
- Searching conversations and, where supported, messages.
- Local SQLite caching.
- Google OAuth 2.0 authentication.

## Development rules

- Build one step at a time.
- Keep the application compiling after every step.
- Do not implement features from later steps early.
- Verify Google Chat API capabilities and required OAuth scopes before relying on them.
- Keep the TUI independent from the Google Chat API implementation.
- Prefer simple, idiomatic Go over premature abstractions.

## Step order

1. Project skeleton and development baseline.
2. Basic Bubble Tea application shell.
3. Configuration and logging.
4. OAuth authentication.
5. Google Chat API client foundation.
6. Conversation listing.
7. Conversation selection and message history.
8. Sending messages.
9. Unread state.
10. Search.
11. SQLite cache.
12. Synchronization and refresh.
13. Error handling and UX polish.
14. Tests and packaging.

Each step has its own Markdown specification.
