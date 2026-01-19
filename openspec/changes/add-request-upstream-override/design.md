## Context
Gemini generateContent requests currently route to a channel-selected upstream and use channel auth. Some clients need to route to a custom proxy backend per request with their own auth headers.

## Goals / Non-Goals
- Goals:
  - Allow request-level upstream override with explicit enablement and allowlist.
  - Preserve quota, accounting, and routing logic outside the override.
  - Avoid leaking auth headers in logs.
- Non-Goals:
  - Generalize to all channels in the first iteration.
  - Support arbitrary HTTP methods beyond existing Gemini relay paths.

## Decisions
- Decision: Add two request headers for override values and apply after channel selection.
- Decision: Enforce allowlist and disable by default to prevent SSRF.
- Decision: Gemini adaptor must not overwrite auth headers when override is supplied.

## Risks / Trade-offs
- Risk: Misconfiguration could route to untrusted hosts. Mitigation: strict allowlist and disabled by default.
- Risk: Override errors may affect clients expecting fallback. Mitigation: clear error response and docs.

## Migration Plan
- Deploy with feature disabled by default.
- Admins enable feature and configure allowlist as needed.

## Open Questions
- Final error codes and response body for blocked overrides.
