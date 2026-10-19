# Week 11 Workshop — Observability, SLOs and Incident Response

Project: **LMS Library** (team `lms-library`). Scope of this workshop: the Circulation
bounded context, which I own (`lms-circulation-api`, `lms-circulation-db`,
`lms-circulation-portal`), and the two services it calls synchronously
(`lms-membership-api`, `lms-catalog-api`).

<!-- CONFIG-START -->
- FULL_NAME: Oscar Mauricio Areiza Paramo
- GITHUB_USER: OscarAreiza
- TEAM: lms-library
- DOCS_REPO: https://github.com/code-corhuila/library-docs
- GOOGLE_DOC: PENDING (link added as soon as the Doc is shared)
- JOURNEY: HU-06 Register a loan — POST /api/v1/loans (lms-circulation-portal -> nginx -> lms-circulation-api -> lms-membership-api + lms-catalog-api + MongoDB loan_db)
<!-- CONFIG-END -->

## Checkpoints

| CP | File | Status |
|---|---|---|
| CP0 | `README.md` | done |
| CP1 | `01-journey.md` | done |
| CP2 | `02-log-contract.md` | done |
| CP3 | `03-instrumentation/` | done |
| CP4 | `04-metrics.md` | done |
| CP6 | `05-slo.md` | todo |
| CP7 | `06-alerts.yml`, `06-alert-policy.md` | todo |
| CP8 | `07-runbook.md` | todo |
| CP9 | `08-incident.md` | todo |
