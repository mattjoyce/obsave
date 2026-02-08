---
id: 10
status: done
priority: High
blocked_by: []
tags: [bug, crash]
---

# Fix Panic When CLI Properties Set With Nil Map

When no config file is loaded (or the config omits `properties`), `config.Properties` is `nil`. If the user passes `-p/--properties`, the default `properties_handling` is `merge`, which calls `mergeProperties` and writes into the nil map, causing a panic.

## Fix

Ensure `config.Properties` is initialized before `addProperties` or `mergeProperties` are called. Options:
- Initialize `config.Properties` in the default config struct, or
- Initialize inside `addProperties` and `mergeProperties` when nil.

## Acceptance Criteria
- `obsave -p "k=v"` does not panic when no config file is present
- `properties_handling=add` and `properties_handling=merge` both work with nil `config.Properties`
- Existing behavior for configs that already define `properties` is unchanged

## Narrative
- 2026-02-08: Created from repo review. This is a user-facing crash when CLI properties are used without config defaults. (by @assistant)
- 2026-02-08: Initialized `config.Properties` in `addProperties` and `mergeProperties` to prevent nil map panics. (by @assistant)
