## Context
Request-level upstream override is configured via options; operators need environment variables for container startup configuration.

## Goals / Non-Goals
- Goals:
  - Support env overrides for enable flag, allowlist, and proxy map.
  - Keep behavior consistent with existing options.
- Non-Goals:
  - Per-request proxy configuration via headers.

## Decisions
- Decision: Parse env overrides in the general setting getter to apply at runtime.
- Decision: Use JSON for list/map values to avoid ambiguous parsing.

## Risks / Trade-offs
- Risk: Invalid JSON in env leads to silent fallback. Mitigation: log a warning and keep existing config.

## Migration Plan
- No migration; defaults remain.
