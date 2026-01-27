## ADDED Requirements
### Requirement: Wildcard and suffix proxy mapping
The system SHALL support wildcard and suffix matching for `request_upstream_proxy_map` keys.

#### Scenario: Wildcard match
- **WHEN** the proxy map contains `*.example.com` and the upstream host is `api.example.com`
- **THEN** the relay MUST apply the mapped proxy

#### Scenario: Suffix match
- **WHEN** the proxy map contains `.example.com` and the upstream host is `example.com` or `api.example.com`
- **THEN** the relay MUST apply the mapped proxy

### Requirement: Proxy map precedence
The system SHALL apply proxy map entries based on precedence: exact host match, then wildcard, then suffix; longer suffix wins.

#### Scenario: Exact match wins
- **WHEN** both `api.example.com` and `*.example.com` are present
- **THEN** the relay MUST use the proxy mapped to `api.example.com`
