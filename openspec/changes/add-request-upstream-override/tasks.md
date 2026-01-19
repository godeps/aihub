## 1. Implementation
- [x] 1.1 Add request-level override settings (enable flag + host allowlist)
- [x] 1.2 Parse override headers and enforce allowlist in middleware
- [x] 1.3 Apply context overrides for base URL and header override
- [x] 1.4 Respect override auth headers in Gemini adaptor
- [x] 1.5 Update OpenAPI docs with new headers
- [x] 1.6 Add logging and validation guardrails

## 2. Tests
- [x] 2.1 Override enabled/allowlisted routes to custom upstream
- [x] 2.2 Override disabled or blocked returns error
- [x] 2.3 Override auth headers take precedence over default key
