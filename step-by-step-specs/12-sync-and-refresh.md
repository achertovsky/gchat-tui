# Step 12 — Synchronization and Refresh

## Objective

Combine the remote Google Chat API with the local cache.

## Requirements

- Load cached conversations and messages quickly at startup.
- Synchronize remote data after startup.
- Provide manual refresh.
- Avoid unnecessary duplicate requests.
- Respect API pagination and rate limits.
- Update changed records without creating duplicates.
- Preserve remote IDs as stable identifiers.
- Handle deleted or inaccessible remote records.
- Make synchronization cancellable.
- Clearly distinguish cached data from fresh data when useful.

Do not add background polling unless it is implemented with bounded resource usage and an explicit user-facing setting.

## Acceptance criteria

- The app starts with cached data when available.
- Remote synchronization updates the cache.
- Duplicate records are not created.
- Manual refresh works.
- Network failures do not destroy usable cached data.
