## ADDED Requirements
### Requirement: Upstream proxy mapping
The system SHALL allow configuration of a per-host proxy mapping for request-level upstream overrides.

#### Scenario: Mapped upstream host
- **WHEN** the override upstream host matches an entry in the proxy map
- **THEN** the relay MUST route the upstream request through the configured proxy for that host

### Requirement: Proxy map default behavior
The system SHALL preserve existing proxy settings when no proxy map entry applies.

#### Scenario: Unmapped upstream host
- **WHEN** the override upstream host does not match any proxy map entry
- **THEN** the relay MUST use the existing proxy configuration without modification
