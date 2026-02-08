---
id: 1
status: done
priority: High
blocked_by: []
tags: [bug, crash]
---

# Fix Panic on Short Paths

`expandAndCleanPath` at `obsave.go:49` uses `path[:2] == "~/"` which will panic with an index out of range if `path` is empty or a single character.

## Fix

Replace `path[:2] == "~/"` with `strings.HasPrefix(path, "~/")`. Consider handling the exact `"~"` case (expand to home dir) for completeness.

## Acceptance Criteria
- Empty string path does not panic
- Single character path does not panic
- `~/foo` still expands correctly
- `~` alone expands to the user home directory (if implemented)

## Narrative
- 2026-02-08: Created from repo review. This is the highest priority bug — any empty or short vault_path triggers a runtime panic. (by @assistant)
- 2026-02-08: Added note to support bare `~` expansion to align with common shell expectations. (by @assistant)
- 2026-02-08: Replaced slice check with `strings.HasPrefix` and added explicit `~` handling to avoid panics and expand home paths. (by @assistant)
