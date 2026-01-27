# Change: Add env-based configuration for request-level upstream override

## Why
Deployments often rely on environment variables for configuration; request-level upstream override should be configurable via env to take effect at startup without API or DB changes.

## What Changes
- Add environment variable overrides for upstream override enablement, allowlist, and proxy map.
- Document the env variables and add a sample in start.sh.

## Impact
- Affected specs: request-upstream-override
- Affected code: setting/operation_setting/general_setting.go, start.sh, docs
