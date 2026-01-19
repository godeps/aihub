# Request-Level Upstream Override (Gemini generateContent)

## Goal
Allow callers to override the upstream base URL and auth headers per request so the system can forward to a user-specified proxy backend without using the channel's default auth.

## Scope
- Primary: Gemini `generateContent` (streaming and non-streaming).
- Extensible to other channels later.

## User Experience
Callers send the standard Gemini request plus two optional headers:
- `X-Relay-Upstream-Base-URL`: The real upstream base URL.
- `X-Relay-Upstream-Headers`: JSON string with auth headers to send upstream.

Example headers:
- `X-Relay-Upstream-Base-URL: https://proxy.example.com`
- `X-Relay-Upstream-Headers: {"Authorization":"Bearer <token>","x-goog-api-key":"<key>"}`

Behavior:
- If the override headers are provided and enabled, the request is forwarded to the specified base URL with the specified auth headers.
- Existing request body and path stay unchanged.

## Product Requirements
- Request-level override must be disabled by default.
- Admin can enable it and configure an allowlist of upstream hosts.
- The system must reject overrides if the upstream host is not in the allowlist.
- The override must not bypass quota, rate-limit, or usage accounting.

## Technical Design
### 1) New Settings
Add to operation settings:
- `request_upstream_override_enabled` (bool, default false)
- `request_upstream_override_allowlist` ([]string, default empty)

### 2) Middleware Override
After channel selection, parse headers:
- Read `X-Relay-Upstream-Base-URL` and `X-Relay-Upstream-Headers`.
- Validate feature enabled.
- Validate URL scheme and host are in allowlist.
- Parse JSON header override map.
- If valid, override:
  - `ContextKeyChannelBaseUrl`
  - `ContextKeyChannelHeaderOverride`

### 3) Gemini Auth Handling
In `relay/channel/gemini/adaptor.go`, avoid overwriting auth when override headers already define auth:
- If override has `x-goog-api-key` or `Authorization`, do not set default `x-goog-api-key`.

### 4) OpenAPI Documentation
Update `docs/openapi/relay.json` for `POST /v1beta/models/{model}:generateContent`:
- Add header params for `X-Relay-Upstream-Base-URL` and `X-Relay-Upstream-Headers`.

### 5) Logging and Auditing
Log whether override is used and the effective upstream base URL (mask sensitive data). Avoid logging auth values.

## Security Considerations
- Strict allowlist enforcement to prevent SSRF.
- Optional host validation to block localhost and private IPs.
- Header JSON size limit to prevent abuse.

## Implementation Tasks (OpenSpec)
1. Add request-level override settings.
2. Implement header parsing and override in middleware.
3. Adjust Gemini adaptor auth behavior.
4. Update OpenAPI docs with new headers.
5. Add allowlist validation and logging.
6. Add tests for override enabled/disabled/invalid cases.

## Acceptance Criteria
- When enabled and allowlisted, requests using override headers are routed to the specified upstream with provided auth headers.
- When disabled or not allowlisted, the override is rejected and default channel behavior remains.
- No auth credentials are logged.
