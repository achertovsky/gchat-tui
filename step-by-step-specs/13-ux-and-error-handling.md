# Step 13 — UX and Error Handling

## Objective

Make the application usable for daily terminal work.

## Requirements

Improve:

- Keyboard navigation.
- Focus management.
- Sidebar and message-panel layout.
- Terminal resizing.
- Help view with keybindings.
- Loading indicators.
- Error messages.
- Empty states.
- Long message wrapping.
- Unicode and non-ASCII text handling.
- Small terminal behavior.
- Graceful shutdown.

Errors should be actionable and should not expose tokens or sensitive request details.

Add confirmation only for destructive actions, if any are introduced.

## Acceptance criteria

- The interface remains usable in a small terminal.
- Every major loading and failure state has visible feedback.
- Keybindings are discoverable.
- Long messages do not break the layout.
- The app exits without leaving terminal settings corrupted.
