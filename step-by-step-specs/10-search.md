# Step 10 — Search

## Objective

Add keyboard-accessible search for conversations and supported message content.

## Requirements

- Enter search mode with `/`.
- Search spaces/channels and direct messages.
- Display selectable search results.
- Open the selected result.
- Debounce or otherwise avoid excessive requests.
- Show loading, empty, and error states.
- Support clearing or exiting search with `Esc`.
- Use server-side Google Chat search only if the current API and scopes support it.
- Otherwise search locally cached data once the cache exists, or clearly limit the MVP to supported conversation search.

Do not claim to search all messages if the API does not provide that capability.

## Acceptance criteria

- Search is accessible from the keyboard.
- Results are selectable.
- Selecting a result opens the relevant conversation.
- Search failures and empty results are clear.
- API limitations are documented.
