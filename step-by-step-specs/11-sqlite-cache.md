# Step 11 — SQLite Cache

## Objective

Add persistent local storage for conversations, messages, and synchronization metadata.

## Requirements

Use SQLite with a migration mechanism.

Store, at minimum:

- Conversations/spaces.
- Conversation type and display metadata.
- Messages.
- Sender metadata needed for rendering.
- Message timestamps and remote IDs.
- Pagination/synchronization metadata.
- Locally observed unread information only where appropriate.

Define indexes for common operations such as:

- Listing conversations.
- Finding messages by conversation.
- Searching cached text.
- Looking up remote IDs.

Use parameterized SQL. Do not store access or refresh tokens in SQLite unless explicitly justified and secured.

## Acceptance criteria

- Database initialization is repeatable.
- Migrations work on a fresh database and an existing database.
- CRUD operations have tests.
- The app can display cached data without network access.
- SQLite errors are handled without corrupting the TUI state.
