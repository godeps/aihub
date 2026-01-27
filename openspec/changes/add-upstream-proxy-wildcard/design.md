## Context
Request-level upstream proxy mapping currently uses exact host matching, which is too rigid for domain families.

## Goals / Non-Goals
- Goals:
  - Support wildcard and suffix matching for proxy map keys.
  - Keep exact matches highest priority.
- Non-Goals:
  - Regex matching or complex patterns.

## Decisions
- Decision: Matching order is exact > wildcard > suffix.
- Decision: `*.example.com` matches subdomains only; `.example.com` matches root and subdomains.

## Risks / Trade-offs
- Risk: Ambiguous mappings if multiple wildcard/suffix entries match. Mitigation: deterministic precedence and longest suffix wins.

## Migration Plan
- No migration; exact match behavior unchanged.
