# Step 7 — Read Message History

## Objective

Open a selected conversation and display its messages.

## Requirements

- Load messages for the selected conversation.
- Display sender, timestamp, and message text.
- Display messages in chronological order.
- Support scrolling.
- Support loading older messages through pagination.
- Handle empty conversations.
- Handle deleted, unavailable, or malformed message data gracefully.
- Show loading and error states.

Keep message rendering separate from API response parsing.

## Acceptance criteria

- Selecting a conversation opens its message view.
- Messages are readable and ordered correctly.
- Scrolling works.
- Older messages can be loaded.
- Empty and failed requests are handled visibly.
