## ADDED Requirements
### Requirement: Request-level upstream override headers
The system SHALL allow a request to specify an upstream base URL and upstream auth headers via request headers when the feature is enabled.

#### Scenario: Override is applied
- **WHEN** `X-Relay-Upstream-Base-URL` and `X-Relay-Upstream-Headers` are provided, the feature is enabled, and the host is allowlisted
- **THEN** the relay MUST use the provided base URL and auth headers for the upstream request

### Requirement: Override allowlist enforcement
The system SHALL reject upstream overrides that are not in the configured allowlist.

#### Scenario: Host not allowlisted
- **WHEN** the override host is not in the allowlist
- **THEN** the relay MUST return an error response and MUST NOT send the request to the override host

### Requirement: Override disabled behavior
The system SHALL reject upstream overrides when the feature is disabled.

#### Scenario: Override disabled
- **WHEN** override headers are provided and the feature is disabled
- **THEN** the relay MUST return an error response and MUST NOT use the override values

### Requirement: Gemini auth header precedence
The system SHALL preserve override auth headers when forwarding Gemini requests.

#### Scenario: Override auth headers provided
- **WHEN** override headers include `x-goog-api-key` or `Authorization`
- **THEN** the relay MUST NOT overwrite those headers with channel defaults

### Requirement: Sensitive header logging
The system SHALL avoid logging override auth header values.

#### Scenario: Logging with override
- **WHEN** override headers are present
- **THEN** logs MUST include only non-sensitive metadata and MUST exclude auth values
