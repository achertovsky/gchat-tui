# Step 6 — List Conversations

## Objective

Replace mock sidebar data with real Google Chat conversations.

## Requirements

- Load spaces/conversations after successful authentication.
- Display direct messages and channels/spaces in the sidebar.
- Show a useful display name for each item.
- Handle unnamed or missing metadata gracefully.
- Show loading and error states.
- Allow manual refresh.
- Keep the TUI responsive while network requests run.

Use Bubble Tea commands or another appropriate asynchronous pattern. Do not block the update loop with network requests.

## Acceptance criteria

- Authenticated users see conversations returned by Google Chat.
- Loading and failure states are visible.
- Selecting a conversation is possible.
- Refresh works.
- The app remains responsive during API calls.
