# Service Level Indicators (SLI) - LAWNE

### 1. **Availability**

- **Indicator**: Uptime percentage.
- **Measurement Tool**: Grafana and other monitoring tools.
- **Formula**: 
  \[
  \text{Availability} = \frac{\text{Total Uptime}}{\text{Total Time}} \times 100
  \]
- **Target**: 99.5% uptime (no more than 0.5% downtime).

### 2. **Response Time**

- **Indicator**: Percentage of requests served in less than 2 seconds.
- **Measurement Tool**: Monitoring tools integrated with the backend (e.g., Grafana, Go logging).
- **Formula**:
  \[
  \text{Response Time} = \frac{\text{Requests Served in < 2s}}{\text{Total Requests}} \times 100
  \]
- **Target**: 95% of requests served in under 2 seconds.

### 3. **Incident Management**

- **Indicator**: Response and resolution times for incidents at different priority levels.
- **Measurement Tool**: Incident tracking system.
- **Formula**: 
  - **High Priority**: Response within 30 minutes, resolution within 4 hours.
  - **Medium Priority**: Response within 2 hours, resolution within 24 hours.
  - **Low Priority**: Response within 8 hours, resolution within 72 hours.
- **Target**: Adherence to the respective times for response and resolution.

### 4. **Security Vulnerabilities**

- **Indicator**: Number of security vulnerabilities discovered.
- **Measurement Tool**: Security audits, penetration testing, automated vulnerability scanning (e.g., OWASP ZAP).
- **Formula**:
  \[
  \text{Security Vulnerabilities} = \text{Total Number of Found Vulnerabilities in the System}
  \]
- **Target**: No security vulnerabilities detected.
