# LAWNE App Development & Service Layer Agreement (SLA)
Project Overview

The LAWNE App, developed by a team of DevOps students at KEA, is a web-based application focused on user authentication, data storage, and real-time analytics. The backend is developed in Go, with a front-end built using TypeScript, CSS, and HTML. The app has undergone significant restructuring, including adopting PostgreSQL and switching from MD5 to bcrypt for security.

**Service Name**: LAWNE App
**Effective Date**: January 4, 2025
**Revision Date**: January 15, 2026

---

## 1. Service Scope
This SLA applies to the LAWNE App, which provides the following services:

* Secure user authentication.
* Real-time analytics dashboards.
* Cloud-based data storage.

---

## 2. Performance Metrics
| **Metric**                         | **Target**        | **Exclusions**           |
|------------------------------------|-------------------|--------------------------|
| Uptime                             | 99.5%             | Planned maintenance      |
| Response Time                      | 2s for 95%        | Network disruptions      |

---

## 3. Incident Management

| **Priority Level**                 | **Response Time** | **Resolution time**      |
|------------------------------------|-------------------|--------------------------|
| High                               | 30 minutes        | 4 hours                  |
| Medium                             | 2 hours           | 24 hours                 |
| Low                                | 8 hours           | 72 hours                 |

---

## 4. Security Measures
1. Data is encrypted using TLS 1.2+ and AES-256 for secure communications.
3. The backend is designed with secure practices, including moving from MD5 to bcrypt for password hashing.

---

## 5. Compliance Standards
The service adheres to the following standards:

1. GDPR (General Data Protection Regulation).
2. ISO 27001 standards for information security.

---

## 6. Compensation Scheme
This section is omitted as the SLA is for academic purposes only.

---

## 7. Revision Clause
This SLA may be revised as needed, with updates communicated to all stakeholders.

---

# Code Quality and Backend Development
## Frameworks and Technologies
* Backend: The backend is built using Go, a programming language chosen for its simplicity and ease of learning for our team. The Go standard package is used extensively, with plans to minimize the use of third-party libraries like Gorilla. We have transitioned to using Gorm, an ORM for Go, which has greatly simplified database interactions.

* Frontend: The frontend is developed using TypeScript, CSS, and HTML. Initially, we considered frameworks like Angular but opted to stick with these technologies due to familiarity.

* Database: PostgreSQL has replaced the previous database due to its superior query performance and flexibility. We have integrated Gorm for database operations, making it easier to interact with the PostgreSQL database.

## Refactor and Code Improvements
The backend code was initially monolithic, with all features in a single large file. This made collaboration difficult, so the project structure was refactored to separate concerns. Key changes include:

File Structure: The backend is now organized into multiple directories such as api, database, models, and security.
MD5 to Bcrypt: The outdated and insecure MD5 hashing algorithm was replaced with bcrypt, improving security.
PostgreSQL Migration: The project switched to PostgreSQL, resulting in faster queries and improved search functionality, especially with pattern matching using ILIKE.

## Commit and Coding Conventions
To maintain consistent code quality, the following conventions are being adhered to:

### General Conventions

| **Concept/Context**                | **Convention**    | **Example**              |
|------------------------------------|-------------------|--------------------------|
| Go Variables and Functions         | Snake case        | my_variable, my_function |
| JavaScript Variables and Functions | Camel case        | 24 hours                 |
| CSS Classes                        | Kebab case        | .my-class, .user-profile |
| Environment Variables              | Upper snake case  | DATABASE_URL, API_KEY    |
| CSS Classes                        | Upper snake case  | MAX_SIZE, DEFAULT_COLOR  |

### Code Quality Issues Identified
* Code Climate detected several issues, some of which we addressed:
    * Missing language tags in HTML files.
    * Minor code duplication was ignored as it didn’t have an immediate impact.
    * We fixed missing imports in test files and replaced MD5 with bcrypt.

## Git Workflow and Branching Strategy
Initially, the team adopted the Gitflow branching strategy, which involves three main branches: main, dev, and qa. However, due to practical constraints, we sometimes merged directly into main when features were deemed production-ready.

### Branching Workflow
1. Feature Branches: Each member worked on their own feature branch.
2. Dev Branch: Features were merged into dev for integration.
3. QA Branch: After integration, the dev branch was reviewed in qa for testing.
4. Main Branch: Once tested, changes were merged into main, which triggered our CI/CD pipeline.

## DevOps Practices
We believe we're making progress in becoming more DevOps-oriented, particularly in the following areas:

* Collaboration
Our team fosters a blame-free environment and focuses on resolving issues, not assigning blame.

* Version Control
Git is used for version control with a structured branching strategy.

* Deployment Automation
We've set up CI/CD pipelines to automate deployments without manual intervention.

However, there are still areas for improvement, such as:
* Monitoring: We’ve only just begun to implement comprehensive monitoring.
* Automation: More automation is needed for provisioning and testing.

## Monitoring and Observability
Our Grafana dashboard has provided valuable insights into our system’s performance. Notably:

*Memory Usage
We’ve realized that our application consumes more memory than expected, which could pose issues if scaling is required.
* Disk Space
Our VM only has 8GB of disk space, which may become restrictive as we scale up the app.

We plan to continuously monitor these metrics and adjust our infrastructure accordingly to ensure smooth scaling in the future.
