## Context
Proxy resolution currently scans the proxy map on each request and cannot explicitly disable a preconfigured channel proxy.

## Goals / Non-Goals
- Goals:
  - Allow a proxy map entry to disable proxy for a host.
  - Reduce repeated map scans with a small cache.
- Non-Goals:
  - Global cross-process cache.

## Decisions
- Decision: Use a reserved value ("none") to indicate no-proxy override.
- Decision: Cache host-to-proxy resolution in-memory and reset cache when settings are reloaded (or with a TTL).

## Risks / Trade-offs
- Risk: Cache staleness when config updates. Mitigation: reset cache when settings are updated on reload.

## Migration Plan
- No migration; new values are optional.
