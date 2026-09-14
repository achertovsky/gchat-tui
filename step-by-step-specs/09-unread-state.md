# Step 9 — Unread Conversations

## Objective

Make unread channels and direct messages easy to identify and navigate to.

## Requirements

- Determine unread state using Google Chat API capabilities available to the application.
- Do not infer authoritative unread state solely from local timestamps unless clearly labeled as a fallback.
- Display unread indicators in the sidebar.
- Sort or group unread conversations in a useful way.
- Mark a conversation read when the supported API behavior allows it.
- Preserve unread state across refreshes where possible.
- Handle APIs that do not expose the exact unread behavior required.

Document any limitations imposed by Google Chat API permissions or capabilities.

## Acceptance criteria

- Unread conversations are visually distinct.
- The user can navigate between unread conversations.
- Read/unread state updates after opening or explicit refresh where supported.
- Unsupported behavior is documented rather than faked.
