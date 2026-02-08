---
id: 6
status: done
priority: Normal
blocked_by: []
tags: [improvement, code-quality]
---

# Clean Up Duplicate scanner.Text() Call

At `obsave.go:531-533`, `scanner.Text()` is called twice per loop iteration — once assigned to `line` and once appended to `contentBuilder`. Not a bug (returns same value) but unnecessary.

## Fix

Use `line` variable for both the passthrough print and the contentBuilder append.

## Acceptance Criteria
- `scanner.Text()` called only once per iteration
- Behavior unchanged

## Narrative
- 2026-02-08: Created from repo review. Minor inefficiency. (by @assistant)
- 2026-02-08: Reused the `line` variable to avoid duplicate `scanner.Text()` calls. (by @assistant)
