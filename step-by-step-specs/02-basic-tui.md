# Step 2 — Basic Bubble Tea TUI

## Objective

Replace the placeholder executable with a full-screen terminal application using Bubble Tea and Lip Gloss.

## Requirements

Create a basic TUI with:

- A sidebar on the left.
- A main content panel on the right.
- A footer showing basic keybindings.
- A placeholder conversation list.
- A placeholder message view.
- Terminal resize handling.
- `q` and `Ctrl+C` to quit.
- Arrow keys or `j`/`k` to navigate the sidebar.
- `Enter` to select a placeholder conversation.

Use Bubble Tea's model/update/view architecture.

## Constraints

- Use mock data only.
- Do not implement authentication or API calls.
- Keep layout code inside `internal/ui`.
- Keep the application entry point small.

## Acceptance criteria

- The app launches in a full-screen terminal UI.
- The layout adapts when the terminal is resized.
- Navigation works without crashing.
- The app exits cleanly.
- `go test ./...` and `go build ./cmd/chat` succeed.
