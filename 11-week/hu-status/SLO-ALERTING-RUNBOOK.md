# Week 11 — Session 2: SLOs, alerting and runbook

Planning for the Circulation context (`lms-circulation-api`, `lms-workflow`), built on the metrics
from Session 1 (`OBSERVABILITY-INSTRUMENTATION.md`, `workshop/04-metrics.md`).

## 1. SLOs

All SLOs use a **28-day rolling window**.

| # | SLO | SLI (good / valid) | Target |
|---|---|---|---|
| SLO-1 | **Availability** of loan registration | Requests to `POST /api/v1/loans` that did **not** end in `outcome="error"` / all requests. Business rejections (`STUDENT_SUSPENDED`, `LOAN_LIMIT_REACHED`, `NO_COPIES_AVAILABLE`) count as good: the service worked. | 99.5 % |
| SLO-2 | **Latency** of loan registration | Non-rejected requests answered in ≤ 0.6 s / all non-rejected requests (0.6 s is an NFR-001 bucket boundary). | 99 % |
| SLO-3 | **Freshness** of overdue loans | Minutes in which the last successful overdue scan of `lms-workflow` is less than 30 min old / all minutes. | 99 % |

PromQL for each SLI (28d):

```promql
# SLO-1 availability
1 - (
  sum(increase(circulation_loan_registrations_total{outcome="error"}[28d]))
  /
  sum(increase(circulation_loan_registrations_total[28d]))
)

# SLO-2 latency
sum(increase(circulation_loan_registration_duration_seconds_bucket{outcome!="rejected",le="0.6"}[28d]))
/
sum(increase(circulation_loan_registration_duration_seconds_count{outcome!="rejected"}[28d]))

# SLO-3 freshness (needs the new gauge from story OBS-03)
avg_over_time(
  ((time() - workflow_overdue_scan_last_success_timestamp_seconds) < bool 1800)[28d:1m]
)
```

### Traffic estimate

There is no production traffic yet. **Assumption:** about 400 loan registrations per day (one
campus library at peak term), which is about 11,200 per 28 days, or about 17 per hour. This
assumption must be checked against real numbers after the MVP 3 release.

## 2. Error budget

| SLO | Budget per 28 days (requests) | Budget per 28 days (time) |
|---|---|---|
| SLO-1 availability 99.5 % | 0.5 % × 11,200 ≈ **56 failed registrations** | 0.5 % × 40,320 min ≈ **202 min** |
| SLO-2 latency 99 % | 1 % × 11,200 ≈ **112 slow registrations** | — |
| SLO-3 freshness 99 % | — | 1 % × 40,320 min ≈ **403 min** stale |

## 3. Error-budget policy

| Budget consumed | Who decides | What is frozen | What is prioritized |
|---|---|---|---|
| 50 % | Circulation owner (Oscar Areiza) | Nothing yet. The team is told in the daily. | Read the burn: which `error_code` / `downstream` causes it. Open a story if the cause is known. |
| 75 % | Circulation owner + team | Risky changes to `lms-circulation-api` and its callers (new features, migrations). | Reliability stories for the cause found, ahead of new features in the sprint. |
| 100 % | Whole team at sprint planning | All feature releases of the Circulation context until the budget is back above 0 for the window. Only fixes and rollbacks are deployed. | Root-cause fix, postmortem action items, and a test that reproduces the failure. |

## 4. Symptom-based alerts

We alert on the **symptom the user sees** (failed or slow loan registrations, stale overdue
status), not on causes. Multi-window, multi-burn-rate on SLO-1:

| Alert | Burn rate | Long window | Short window | Budget consumed when it fires | Severity |
|---|---|---|---|---|---|
| `LoanRegistrationBudgetBurnFast` | 14.4 | 1 h | 5 m | 14.4 × 1 h / 672 h ≈ **2.1 %** | page |
| `LoanRegistrationBudgetBurnSlow` | 6 | 6 h | 30 m | 6 × 6 h / 672 h ≈ **5.4 %** | ticket |
| `OverdueScanStale` | — | — | — | the scan is more than 30 min old for 10 min | ticket |

**Owner:** Circulation owner (Oscar Areiza). **Runbook:** section 5 of this document.

At about 17 requests per hour, 1 or 2 errors can trip a burn-rate alert. Each rule therefore also
requires a minimum number of requests in the window (`> 10`), so one failed loan at night does not
page anyone.

```yaml
groups:
  - name: lms-circulation-slo
    rules:
      - record: circulation:loan_error_ratio:rate5m
        expr: sum(rate(circulation_loan_registrations_total{outcome="error"}[5m])) / sum(rate(circulation_loan_registrations_total[5m]))
      - record: circulation:loan_error_ratio:rate1h
        expr: sum(rate(circulation_loan_registrations_total{outcome="error"}[1h])) / sum(rate(circulation_loan_registrations_total[1h]))
      - record: circulation:loan_error_ratio:rate30m
        expr: sum(rate(circulation_loan_registrations_total{outcome="error"}[30m])) / sum(rate(circulation_loan_registrations_total[30m]))
      - record: circulation:loan_error_ratio:rate6h
        expr: sum(rate(circulation_loan_registrations_total{outcome="error"}[6h])) / sum(rate(circulation_loan_registrations_total[6h]))

      - alert: LoanRegistrationBudgetBurnFast
        expr: |
          circulation:loan_error_ratio:rate1h > (14.4 * 0.005)
          and circulation:loan_error_ratio:rate5m > (14.4 * 0.005)
          and sum(increase(circulation_loan_registrations_total[1h])) > 10
        for: 2m
        labels:
          severity: page
          team: lms-circulation
        annotations:
          summary: "Loan registration is burning the 28-day error budget 14x too fast"
          runbook_url: "https://github.com/OscarAreiza/sistemas-distribuidos-2026-b-g2/blob/main/11-week/hu-status/SLO-ALERTING-RUNBOOK.md#5-runbook-loanregistrationbudgetburnfast"

      - alert: LoanRegistrationBudgetBurnSlow
        expr: |
          circulation:loan_error_ratio:rate6h > (6 * 0.005)
          and circulation:loan_error_ratio:rate30m > (6 * 0.005)
          and sum(increase(circulation_loan_registrations_total[6h])) > 10
        for: 15m
        labels:
          severity: ticket
          team: lms-circulation
        annotations:
          summary: "Loan registration is burning the 28-day error budget 6x too fast"
          runbook_url: "https://github.com/OscarAreiza/sistemas-distribuidos-2026-b-g2/blob/main/11-week/hu-status/SLO-ALERTING-RUNBOOK.md#5-runbook-loanregistrationbudgetburnfast"

      - alert: OverdueScanStale
        expr: (time() - workflow_overdue_scan_last_success_timestamp_seconds) > 1800
        for: 10m
        labels:
          severity: ticket
          team: lms-circulation
        annotations:
          summary: "Overdue loan status is more than 30 minutes old"
```

### Cause-based alerts we will NOT create

1. **High CPU or memory** on a container. Users only notice it if loans fail or slow down, and
   SLO-1 / SLO-2 already catch that.
2. **MongoDB / PostgreSQL pool above 80 %.** It stays on the dashboard as a diagnosis aid in the
   runbook, not as a page.
3. **One downstream timeout.** Retries and the budget absorb single timeouts; a sustained problem
   shows up as a burn-rate alert.

## 5. Runbook: `LoanRegistrationBudgetBurnFast`

### Symptoms

Administrators get errors (`INTERNAL_ERROR`) when they register a loan in `lms-circulation-portal`,
or the request times out. The alert fired because more than 7.2 % of loan registrations failed in
the last hour and the last 5 minutes.

### How to confirm

```promql
# Error ratio right now, and which error code / downstream causes it
sum by (error_code) (rate(circulation_loan_registrations_total{outcome="error"}[5m]))
sum by (downstream, operation, result) (rate(circulation_downstream_requests_total{result=~"timeout|error"}[5m]))
histogram_quantile(0.95, sum by (le, downstream) (rate(circulation_downstream_request_duration_seconds_bucket[5m])))
```

```bash
# Take one failing trace_id from an exemplar or an error log, then follow it through every service
docker compose -f lms-pod/compose.yml logs --since 15m circulation-service membership-service catalog-service \
  | grep '"trace_id":"<TRACE_ID>"'
```

### First mitigations (ordered by risk, lowest first)

1. **Downstream is the cause** (`downstream="lms-catalog-api"` or `lms-membership-api` shows
   timeouts or errors): check that service's `/health/ready` and its database pool. If the
   container is unhealthy, restart only that container (`docker compose restart <service>`).
2. **A recent deploy is the cause** (errors started right after a release): roll back
   `lms-circulation-api` to the previous tag with the week 12 script:
   `./deploy-rolling.sh circulation-service <previous-tag>` (see `12-week/hu-status/deployment-plan.md`).
3. **MongoDB `loan_db` is the cause** (`downstream="mongodb.loan_db"`, pool `in_use` at its
   maximum): restart `circulation-service` to release connections, then check `loan_db` disk and
   locks.
4. **Copies decremented with no Loan** (catalog committed, circulation timed out): after recovery,
   list the copies to restore and fix them by hand with the catalog owner. Never re-run the
   requests blindly, because `loan-copy` is not idempotent.

### Escalation path

| Order | Role | Person | When |
|---|---|---|---|
| 1 | Circulation owner (alert owner) | Oscar Areiza | When the alert fires |
| 2 | Teammates working on the failing downstream service (catalog or membership) | Luis Meneses, Hermes Pascuas | When the cause is in catalog or membership, or after 30 min with no recovery |
| 3 | Whole team | lms-library team | Budget fully consumed, or data needs manual repair |

### Verify recovery

- `circulation:loan_error_ratio:rate5m` is back below 0.5 % for 15 minutes, and the alert has
  resolved.
- A test loan registered from the portal returns `201` and its `trace_id` appears in all three
  services' logs.
- Write down the budget consumed by the incident and open the postmortem if it was over 5 %.

## 6. Observability stories (sliced, with testable acceptance criteria)

| ID | Story | Acceptance criteria |
|---|---|---|
| OBS-01 | As the Circulation owner, I want `lms-circulation-api` to expose the SLI metrics, so that SLO-1 and SLO-2 can be measured. | Given the service is running, when `GET /metrics` is called on the internal network, then it returns `circulation_loan_registrations_total` and `circulation_loan_registration_duration_seconds` with only the labels in the metric catalog. Given a loan is registered, then the counter for `outcome="registered"` increases by exactly 1. |
| OBS-02 | As the Circulation owner, I want the SLO recording and alert rules loaded in Prometheus, so that budget burn is detected. | `promtool check rules` passes on the rules file. A `promtool test rules` case with 10 % errors for 1 h fires `LoanRegistrationBudgetBurnFast`; a case with 0.3 % errors does not. |
| OBS-03 | As the Circulation owner, I want `lms-workflow` to publish the time of its last successful overdue scan, so that SLO-3 can be measured. | After a successful scan, `workflow_overdue_scan_last_success_timestamp_seconds` equals the scan's finish time (±5 s). After a failed scan the value does not change. |
| OBS-04 | As an on-call person, I want each alert to link to a runbook, so that I can act without guessing. | Every alert rule has a `runbook_url` that opens an existing section. A drill on the local pod (stop `catalog-service` for 10 min) fires the fast-burn alert, and following the runbook restores the error ratio below 0.5 % within 15 min. |
| OBS-05 | As the Circulation owner, I want the real `catalogclient` and `membershipclient` to forward `traceparent`, so that one `trace_id` follows HU-06 across services in the deployed system. | When `POST /api/v1/loans` is called with a given `traceparent`, the log lines of circulation, membership and catalog for that request all carry the same `trace_id`. An integration test asserts this. |
