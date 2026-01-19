# Change: Add request-level upstream override for Gemini generateContent

## Why
Clients need to route Gemini image generation through user-specified proxy backends with their own auth, without relying on channel default credentials.

## What Changes
- Add request-level upstream base URL override and auth header override for Gemini generateContent.
- Add settings to enable the feature and constrain upstream hosts via allowlist.
- Update Gemini request header behavior to respect override auth headers.
- Document new headers in OpenAPI.

## Impact
- Affected specs: request-upstream-override
- Affected code: middleware, relay/channel/gemini, openapi docs, settings
