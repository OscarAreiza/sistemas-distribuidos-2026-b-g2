# Non-Functional Requirements (NFR) — LMS-LIBRARY-V1


## NFR-001: Performance

| Attribute | Metric | Test condition |
|---|---|---|
| P95 latency — critical endpoints | < 300ms | Under 50 concurrent RPS |
| P99 latency — critical endpoints | < 600ms | Under 50 concurrent RPS |
| P95 latency — non-critical endpoints | < 1000ms | Normal load |
| Minimum throughput | 20 RPS | Without degradation |
| Service startup time | < 30 seconds | Cold start |

**Defined critical endpoints:**
- `POST /loans` — registers a loan; must respond quickly so it doesn't bottleneck the circulation desk (FR-010).
- `POST /returns` — registers a return; same reasoning, and unlocks availability for the next borrower (FR-015).
- `GET /books?search=` — catalog search; used constantly by administrators to locate books (FR-004).

**Load testing tools:** k6, Apache JMeter, Locust, Gatling.
**Where is it validated?** CI/CD in the staging pipeline before deployment.

> **Provisional assumption:** the endpoint paths above are illustrative — the final API contract will live in `07-api/contracts/openapi/` once defined. Load figures (RPS, concurrency) reuse the academic-scale baseline assumed elsewhere in this document, since the Product Brief does not specify expected usage volume.

## NFR-002: Availability

| Environment | SLO | Maintenance window | Max downtime/month |
|---|---|---|---|
| Production (staging/demo deployment) | 95% | Flexible, coordinated with the team | ~36 hours |
| Local / development | Best-effort | No restriction | N/A |

**Monthly error budget:** ~36 hours.
**Error budget policy:** if more than 50% of the monthly error budget is consumed in the first half of the month, feature deploys are paused and stability work is prioritized until the next month.

**Health checks:**
- `GET /health` — liveness: responds 200 if the process is alive.
- `GET /health/ready` — readiness: responds 200 only if it can process traffic (DB connected, dependencies OK).

> **Note:** 95% (not 99.9%) is used deliberately — a self-hosted academic deployment cannot realistically promise enterprise-grade SLOs, per the Product Brief's guidance against unrealistic targets (§16).

## NFR-003: Scalability

| Scenario | Expected behavior |
|---|---|
| Gradual load growth | Horizontal scaling of stateless service instances when CPU > 70%, if the deployment target supports it |
| Sudden spike (e.g., peak enrollment week) | System does not crash under a 2x burst; graceful queuing/backpressure preferred over dropped requests |
| Load reduction | Scale-down without interrupting active requests |
| Horizontal scaling limit | Up to 3 instances per service for V1 |

**Strategy:** stateless horizontal scaling — no instance stores state in memory. Persistent data lives in PostgreSQL; any cache/session state introduced later goes in an external store, not in-process memory.

> **Provisional assumption:** true autoscaling assumes a container-orchestration target beyond plain `docker-compose` (e.g., Kubernetes or a managed equivalent). The Product Brief only confirms Docker + "cloud-hosted / virtualized infrastructure" (§3); whether V1 runs under an orchestrator with autoscaling is a **pending architecture decision**. This section documents the design property (statelessness) that keeps that option open — it is not a confirmed autoscaling implementation for V1.

## NFR-004: Security

**Authentication and Authorization**
- All administrative operations require a valid token in the `Authorization: Bearer <token>` header.
- Token expires in 1 hour; refresh token (if implemented) valid for 7 days.
- RBAC: not applicable in v1 — a single Administrator role exists, with no role hierarchy (confirmed in `01-context/scope.md` and `02-domain/domain-map.md`); revisit only if a future version introduces multiple roles.

**Data transmission**
- HTTPS mandatory in production (TLS 1.2+).
- HTTP allowed only in local development.

**Sensitive data**
- Passwords: hashing with bcrypt (cost factor ≥ 12) or Argon2id.
- PII (library-user personal data): encrypted at rest.
- Secrets/keys: only in environment variables or a vault, never in code.

**OWASP Top 10**
Code should be reviewed against the OWASP Top 10 on each release. Suggested tools: SAST (SonarQube/Snyk), dependency scanning, DAST in staging.

**Regulatory compliance**
Habeas Data (Colombian personal-data protection law, Ley 1581 de 2012) applies given the university context and the storage of library-user PII. Exact compliance measures beyond encryption-at-rest are a **pending business decision**.

> **Provisional assumption:** a token-based (JWT-style) authentication mechanism is assumed for concreteness. A session-based alternative is equally compatible with FR-024/FR-025 — the final choice belongs to the architecture stage.

## NFR-005: Observability

| Pillar | Requirement | Suggested tool |
|---|---|---|
| Logs | Structured JSON format + Correlation ID | `zap` or `logrus` (Go) |
| Metrics | RED (Rate, Errors, Duration) per endpoint | Prometheus + Grafana |
| Traces | End-to-end distributed traces | OpenTelemetry + Jaeger |
| Alerts | Alert in < 5 min when an SLI violates its SLO | Alertmanager (or manual dashboard review, given academic scale) |

**Correlation ID:** each external request generates a UUID `correlationId` propagated through all logs and spans of that transaction.

## NFR-006: Maintainability

| Metric | Target |
|---|---|
| Test coverage | ≥ 80% of lines (≥ 90% in the domain/business-logic layer) |
| Cyclomatic complexity | ≤ 10 per function (e.g., checked with `gocyclo` for Go) |
| Technical debt | Resolution time < 1 sprint from registration |
| Onboarding time | A new team member can run the project locally in < 1 hour following `10-devops/local-setup.md` |
| Average build time | < 5 minutes in CI |

## NFR-007: Portability

- All services are deployed as Docker images.
- Images are runnable via `docker-compose` for local/dev, and compatible with Kubernetes 1.28+ if/when the team adopts container orchestration for deployment.
- No service depends on the host operating system.
- Environment variables are the only source of environment-specific configuration.

> **Provisional assumption:** Kubernetes compatibility is a forward-looking target, not a confirmed V1 deployment choice — see NFR-003.

## NFR-008: Disaster Recovery (DR / Recovery)

| Scenario | RTO (Recovery Time Objective) | RPO (Recovery Point Objective) |
|---|---|---|
| Single service failure | < 5 minutes (manual or automated restart) | 0 (stateless) |
| Primary database failure | < 30 minutes (restore from latest backup) | ≤ 24 hours (daily backup) |
| Full infrastructure loss | < 1 business day (redeploy from repo + restore latest backup) | ≤ 24 hours |

> **Note:** the classic multi-availability-zone / multi-region DR scenario is not realistic for a single-region academic Docker deployment, per the Product Brief's guidance against unrealistic enterprise targets (§16, §21). This table keeps the RTO/RPO framework but scales the scenarios to what the actual deployment can promise. Backup frequency (daily) is a **provisional assumption**, pending Open Question #15.

---

## NFR priority matrix

| NFR | Priority (P1/P2/P3) | Validated in CI? | Owner |
|---|---|---|---|
| Performance | P2 | Planned (k6 in staging) | [To be assigned] |
| Availability | P2 | Partial (health checks only) | [To be assigned] |
| Security | P1 | Planned (SAST + auth tests) | [To be assigned] |
| Scalability | P3 | Manual, not automated | [To be assigned] |
| Observability | P3 | Manual (log review) | [To be assigned] |
| Maintainability | P1 | Yes (coverage report in CI) | [To be assigned] |
| Portability | P2 | Yes (container build test) | [To be assigned] |
| Disaster Recovery | P3 | Manual (drill, not automated) | [To be assigned] |

---

## Correlations

- Detailed SLOs and SLAs → `13-operations/README.md`
- Pipeline that validates NFRs → `10-devops/README.md`
- Incidents related to NFR violations → `13-operations/incident-management.md`
- Security checklist → `00-governance/security-policy.md`
