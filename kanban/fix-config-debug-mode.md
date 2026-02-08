---
id: 3
status: done
priority: Normal
blocked_by: []
tags: [bug]
---

# Fix Debug Logging From Config

Setting `debug: true` in the config file doesn't enable the `debugMode` global variable. Only the CLI `--debug` flag sets it. This means config-based debug mode is partially broken — `config.Debug` is set but `debugMode` (used by `debugLog`) remains false.

## Fix

After loading config, sync `debugMode` from `config.Debug`. Also set the same log flags used by the CLI path so debug output format is consistent regardless of how debug is enabled.

## Acceptance Criteria
- `debug: true` in config file enables debug log output
- CLI `--debug` still works as before
- Debug log format (date/time/file) is the same for config and CLI

## Narrative
- 2026-02-08: Created from repo review. The global `debugMode` and `config.Debug` are out of sync when debug is set via config. (by @assistant)
- 2026-02-08: Clarified that enabling debug via config should also set log formatting flags to match CLI behavior. (by @assistant)
- 2026-02-08: Synced `debugMode` and log flags from `config.Debug` after config load. (by @assistant)
