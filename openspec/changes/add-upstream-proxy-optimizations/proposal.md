# Change: Add proxy map optimizations (no-proxy override + cache)

## Why
Operators need to explicitly disable proxy for certain upstreams even if a channel proxy exists, and resolve proxy map lookups efficiently.

## What Changes
- Support explicit "no proxy" override in `request_upstream_proxy_map` values.
- Add a small in-memory cache for host-to-proxy resolution.
- Document the new behavior.

## Impact
- Affected specs: request-upstream-override
- Affected code: middleware/request_upstream_override.go, docs
