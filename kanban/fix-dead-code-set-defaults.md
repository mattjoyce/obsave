---
id: 2
status: done
priority: Normal
blocked_by: []
tags: [bug, dead-code]
---

# Fix Dead Code - setDefaultsAndOverrides Never Called

The function `setDefaultsAndOverrides` at `obsave.go:67` validates the overwrite mode value, but it is never called from `main()`. The overwrite mode is set directly without validation, so invalid values from config files are silently accepted.

## Fix Options
- Call `setDefaultsAndOverrides` after config loading and CLI overrides in `main()`, or
- Inline the validation logic where overwrite mode is set (and delete dead code)

## Acceptance Criteria
- Invalid `overwrite_mode` values from config or CLI are caught and defaulted to "fail"
- Debug output (if enabled) records the final overwrite mode
- No dead code remains

## Narrative
- 2026-02-08: Created from repo review. The validation function exists but was never wired in. (by @assistant)
- 2026-02-08: Clarified that both config and CLI values should be validated and that debug output should reflect final mode. (by @assistant)
- 2026-02-08: Wired `setDefaultsAndOverrides` into `main()` after CLI overrides to validate and normalize overwrite mode. (by @assistant)
