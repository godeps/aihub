# Change: Add upstream proxy mapping for request-level override

## Why
Some upstream domains must be reached through a corporate HTTP proxy, and this needs to be controlled per domain when request-level upstream override is used.

## What Changes
- Add a configuration map that routes specific upstream hosts through a configured proxy.
- Apply the proxy mapping during request-level upstream override handling.
- Document the new configuration field and behavior.

## Impact
- Affected specs: request-upstream-override
- Affected code: middleware/request_upstream_override.go, setting/operation_setting/general_setting.go, docs
