---
id: 11
status: done
priority: Normal
blocked_by: []
tags: [bug, error-handling]
---

# Fix Panic When Specified Config Has YAML Parse Errors

If a user passes `-c/--config` and the YAML is invalid, `loadConfig` returns `(nil, nil, error)`. The caller then does `log.Fatalf(..., *specifiedConfigPath)` which dereferences a nil pointer and panics, hiding the real error.

## Fix

Guarantee a non-nil config path in error cases, or guard the dereference in `main()` before printing the path.

## Acceptance Criteria
- Invalid YAML in a specified config file yields a clean error message without a panic
- Error output still includes the config path when available

## Narrative
- 2026-02-08: Created from repo review. Current error handling can panic and mask YAML parse errors. (by @assistant)
- 2026-02-08: Ensured `loadConfig` returns a non-nil config path on YAML parse errors to avoid nil deref in error reporting. (by @assistant)
