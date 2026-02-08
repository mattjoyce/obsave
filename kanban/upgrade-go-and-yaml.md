---
id: 7
status: done
priority: Normal
blocked_by: []
tags: [improvement, dependencies]
---

# Upgrade go.mod to Newer Go Version and yaml.v3

`go.mod` targets Go 1.19 (Aug 2022) and uses `gopkg.in/yaml.v2`. Both are outdated.

## Acceptance Criteria
- `go.mod` updated to a recent Go version
- Dependency changed from `yaml.v2` to `yaml.v3`
- Code updated for any yaml.v3 API differences
- `go build` succeeds
- `go.sum` updated accordingly

## Narrative
- 2026-02-08: Created from repo review. yaml.v3 provides better error messages and YAML spec compliance. (by @assistant)
- 2026-02-08: Added explicit `go.sum` update requirement. (by @assistant)
- 2026-02-08: Upgraded to Go 1.22 and `gopkg.in/yaml.v3`, updated imports, and ran `go mod tidy`. (by @assistant)
