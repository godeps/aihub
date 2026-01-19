## Context
Request-level upstream override now supports routing to arbitrary upstreams. Some upstream hosts require an HTTP proxy to be reachable.

## Goals / Non-Goals
- Goals:
  - Allow per-host proxy configuration for request-level upstream overrides.
  - Keep routing logic centralized in middleware.
  - Preserve existing channel proxy behavior when no mapping applies.
- Non-Goals:
  - Dynamic proxy specification from client request.
  - Changes to non-override traffic.

## Decisions
- Decision: Add a `request_upstream_proxy_map` setting keyed by host.
- Decision: Apply mapping after allowlist validation and before request execution.

## Risks / Trade-offs
- Risk: Incorrect mapping could route through unintended proxies. Mitigation: configuration-only, admin-controlled.

## Migration Plan
- Default map is empty; no behavior change until configured.
