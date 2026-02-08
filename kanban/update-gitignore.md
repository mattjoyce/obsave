---
id: 5
status: done
priority: Normal
blocked_by: []
tags: [improvement, hygiene]
---

# Add Binary and Test Vaults to .gitignore

The compiled `obsave` binary and test vault directories (`vault/`, `vault1/`, `vault2/`) are untracked but not gitignored. On Windows, `obsave.exe` may also appear.

## Acceptance Criteria
- `.gitignore` includes `obsave`, `obsave.exe`, `vault/`, `vault1/`, `vault2/`
- Untracked files no longer show in `git status`

## Narrative
- 2026-02-08: Created from repo review. (by @assistant)
- 2026-02-08: Added Windows binary (`obsave.exe`) to ignore list. (by @assistant)
- 2026-02-08: Added `obsave` binary, vault directories, and Zone.Identifier patterns to `.gitignore`. (by @assistant)
