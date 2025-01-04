# Service Level Objectives (SLO) - LAWNE 

## Introduction

The LAWNE App is a web-based application focused on user authentication, data storage, and real-time analytics. Developed by DevOps students at KEA, the app is built using Go for the backend and TypeScript, CSS, and HTML for the frontend. This document outlines the Service Level Objectives (SLOs) to ensure the app delivers reliable, secure, and high-performance service.

## Service Level Objectives (SLOs)

### 1. **Availability**

- **Objective**: The service should be available 99.5% of the time.
- **Measurement**: The uptime will be tracked using monitoring systems such as Grafana.
- **Target**: A maximum of 0.5% downtime per month, excluding planned maintenance.

### 2. **Response Time**

- **Objective**: The system should provide responses to user requests within 2 seconds for 95% of the time.
- **Measurement**: Response times will be monitored using backend logging tools integrated with Grafana.
- **Target**: 95% of responses should be served in less than 2 seconds, excluding network disruptions.

### 3. **Incident Management Response Time**

- **Objective**: Incidents of varying priority levels should be handled promptly.
- **Measurement**: Response and resolution times are tracked in the incident management system.
- **Target**:
  - **High Priority**: Response in 30 minutes and resolution within 4 hours.
  - **Medium Priority**: Response in 2 hours and resolution within 24 hours.
  - **Low Priority**: Response in 8 hours and resolution within 72 hours.

### 4. **Security**

- **Objective**: Sensitive data must be protected, and security practices should be up-to-date.
- **Measurement**: Security assessments, code reviews, and penetration testing will be carried out regularly.
- **Target**: No security vulnerabilities should be present, and all data must be encrypted using TLS 1.2+ and AES-256.

---

## Conclusion

The SLOs ensure that the LAWNE App meets high standards for availability, response times, and security. Monitoring and incident management procedures will help track performance and handle any issues that arise.
