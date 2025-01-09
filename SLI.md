# Service Level Indicators (SLI) - LAWNE

## 1. Availability

- **Indicator**: Uptime percentage.
- **Measurement Tool**: Grafana and other monitoring tools.
- **Formula**:  
  `Availability = (Total Uptime / Total Time) × 100`
- **Target**: 99.5% uptime (no more than 0.5% downtime).

## 2. Response Time

- **Indicator**: Percentage of requests served in less than 2 seconds.
- **Measurement Tool**: Backend monitoring tools (e.g., Grafana, Go logging).
- **Formula**:  
  `Response Time = (Requests Served in < 2s / Total Requests) × 100`
- **Target**: 95% of requests served in under 2 seconds.

## 3. Incident Management

- **Indicator**: Response and resolution times for incidents at different priority levels.
- **Measurement Tool**: Incident tracking system.
- **Target**:
  - **High Priority**: Response within 30 minutes, resolution within 4 hours.
  - **Medium Priority**: Response within 2 hours, resolution within 24 hours.
  - **Low Priority**: Response within 8 hours, resolution within 72 hours.

## 4. Security Vulnerabilities

- **Indicator**: Number of security vulnerabilities discovered.
- **Measurement Tool**: Security audits, penetration testing, automated vulnerability scanning (e.g., OWASP ZAP).
- **Formula**:  
  `Security Vulnerabilities = Total Number of Found Vulnerabilities in the System`
- **Target**: No security vulnerabilities detected.
