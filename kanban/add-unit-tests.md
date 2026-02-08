---
id: 4
status: done
priority: Normal
blocked_by: []
tags: [improvement, testing]
---

# Add Unit Tests

No test files exist in the repository. Good candidates for unit tests:

- `replaceTags`, `addTags`, `mergeTags`
- `replaceProperties`, `addProperties`, `mergeProperties`
- `createFrontmatter`
- `expandAndCleanPath`
- `loadConfig`
- `saveToObsidian` (with temp directories)
- property handling when `config.Properties` is nil
- config parse error handling (should not panic)

## Acceptance Criteria
- `obsave_test.go` exists with tests for core functions
- `go test ./...` passes

## Narrative
- 2026-02-08: Created from repo review. Zero test coverage currently. (by @assistant)
- 2026-02-08: Added coverage targets for nil properties handling and config parse error behavior. (by @assistant)
- 2026-02-08: Added `obsave_test.go` covering tags/properties, frontmatter, config loading, path expansion, and save behavior; `go test ./...` passes. (by @assistant)
