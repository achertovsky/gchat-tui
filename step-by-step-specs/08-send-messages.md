# Step 8 — Send Messages

## Objective

Allow the user to send plain-text messages in a selected conversation.

## Requirements

- Add a message input field.
- Provide a clear way to focus the input.
- Send the message using the Google Chat API.
- Prevent sending empty or whitespace-only messages.
- Show sending, success, and failure states.
- Avoid losing the input text when sending fails.
- Add the sent message to the view after confirmed success.
- Prevent accidental duplicate sends where practical.

Initially support plain text only. Do not implement attachments, reactions, threads, or rich formatting yet.

## Acceptance criteria

- A user can type and send a message.
- Empty messages are rejected locally.
- Successful sends appear in the conversation.
- Failed sends preserve the draft and show an error.
- The TUI remains responsive during sending.
