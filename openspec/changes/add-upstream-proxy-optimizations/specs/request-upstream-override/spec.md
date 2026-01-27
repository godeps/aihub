## ADDED Requirements
### Requirement: No-proxy override
The system SHALL allow a proxy map entry to explicitly disable proxy usage for a matched upstream host.

#### Scenario: No-proxy entry
- **WHEN** the proxy map entry for a matched host is set to "none"
- **THEN** the relay MUST clear any existing channel proxy for that request

### Requirement: Proxy map caching
The system SHALL cache proxy map resolution results per host within a process.

#### Scenario: Cached resolution
- **WHEN** the same upstream host is resolved multiple times
- **THEN** the relay MUST return the cached result without re-scanning the map
