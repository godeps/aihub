## ADDED Requirements
### Requirement: Environment variable overrides
The system SHALL support environment variables that override request-level upstream override settings at runtime.

#### Scenario: Enable override via env
- **WHEN** the enable env variable is set to true
- **THEN** the request-level upstream override MUST be enabled

### Requirement: Env allowlist parsing
The system SHALL parse allowlist values from environment variables.

#### Scenario: Allowlist env value provided
- **WHEN** the allowlist env variable contains a JSON array of hosts
- **THEN** the allowlist MUST be set to the parsed host list

### Requirement: Env proxy map parsing
The system SHALL parse the upstream proxy map from environment variables.

#### Scenario: Proxy map env value provided
- **WHEN** the proxy map env variable contains a JSON object
- **THEN** the proxy map MUST be set to the parsed host-to-proxy mapping
