---
id: 9
status: done
priority: Normal
blocked_by: []
tags: [improvement, consistency]
---

# Use Consistent Home Directory Method

Two different methods are used to get the home directory:
- `loadConfig` uses `user.Current()` (`obsave.go:92`)
- `expandAndCleanPath` uses `os.UserHomeDir()` (`obsave.go:51`)

Should pick one and use it consistently. `os.UserHomeDir()` is simpler and preferred.

## Acceptance Criteria
- Single method used throughout for home directory resolution
- Behavior unchanged

## Narrative
- 2026-02-08: Created from repo review. `os.UserHomeDir()` is the simpler, more modern approach. (by @assistant)
- 2026-02-08: Switched `loadConfig` to use `os.UserHomeDir()` for consistency. (by @assistant)
