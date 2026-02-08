---
id: 12
status: done
priority: Normal
blocked_by: []
tags: [docs]
---

# Document the `--passthrough` Flag

The code supports a `--passthrough` flag that echoes stdin to stdout while saving to the vault, but it is not mentioned in `README.md` or `config.example`.

## Fix

Add `--passthrough` to the README options list and include a brief usage example. Decide whether it belongs in `config.example` and document accordingly.

## Acceptance Criteria
- README lists `--passthrough` and explains its behavior
- README includes a minimal example showing passthrough use
- Any config example updates remain accurate for users who don't want passthrough

## Narrative
- 2026-02-08: Created from repo review. The flag exists in code but is undocumented. (by @assistant)
- 2026-02-08: Documented `--passthrough` in README and added it to the example config. (by @assistant)
