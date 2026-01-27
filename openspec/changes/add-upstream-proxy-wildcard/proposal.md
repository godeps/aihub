# Change: Add wildcard and suffix matching for upstream proxy map

## Why
Proxy mapping currently supports exact host matches only. Operators need wildcard and suffix matching to route whole domain families through a proxy.

## What Changes
- Support wildcard entries (e.g., `*.example.com`) and suffix matches (e.g., `.example.com`) in `request_upstream_proxy_map` keys.
- Document matching rules and precedence.

## Impact
- Affected specs: request-upstream-override
- Affected code: middleware/request_upstream_override.go, docs
