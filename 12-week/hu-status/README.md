<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       12-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 12

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Oscar Mauricio Areiza Paramo
- GITHUB_USER: OscarAreiza
- TEAM: lms-library
- SPRINT_GOAL: Deploy the tagged v2.0.1 cut to a production-like host with a zero-downtime strategy and real health checks (Session 1), then formalize it as a written deployment plan for MVP 3 — artifact promotion, IaC, the chosen strategy, and a rehearsed/timed rollback (Session 2).
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| W12-S1 | Deploy with a zero-downtime (rolling) strategy and health checks | doing | `deploy-rolling.sh`, `deployment-plan.md` §3-4 (this folder) — script written and reviewed, not yet run against a live Docker daemon |
| W12-S2 | Write the MVP 3 deployment plan: artifact promotion, IaC, strategy, timed rollback, sliced stories | done | `deployment-plan.md` (this folder) |

## 2. My individual contribution
- Chose **rolling** over **blue-green** for the zero-downtime strategy, and wrote the justification
  in terms of this project's own real constraints (single Docker host, no cloud budget, no spare
  capacity for a doubled environment — `01-context/scope.md`), not a generic comparison.
- Wrote `deploy-rolling.sh`: builds a given git tag into an image, starts it alongside the running
  container, polls the service's own real `GET /health/ready` endpoint (already implemented in
  every `-api`, verified in each repo's health handler) until it reports ready, only then stops the
  old container. The same script doubles as the rollback mechanism, pointed at the previous tag.
- Wrote the MVP 3 deployment plan (`deployment-plan.md`): artifact promotion reusing this project's
  already-accepted `develop`→`qa`→`main` cherry-pick model (no parallel process invented), IaC as
  the real `compose.yml` files this team already has (`lms-pod`, `lms-infra-postgres`,
  `lms-infra-mongo` — not a Terraform rewrite disproportionate to a 3-person team), and 5
  deployment stories with concrete, checkable acceptance criteria.

## 3. Blockers and risks
- The Docker daemon was not reachable in the environment this was written in (`docker info` fails
  to connect) — `deploy-rolling.sh` has **not** actually been run yet, and the execution log /
  rollback timings in `deployment-plan.md` are explicitly marked as illustrative, not captured
  real output. Running it for real and replacing those sections is next.
- No container registry and no CI exist yet (an accepted, already-documented gap — `ADR-006`), so
  "artifact promotion" today means promoting a git tag and rebuilding locally per environment, not
  a pushed image. Closing this is explicitly the CI/CD milestone, not this week's.

## 4. Plan for next week
- Start Docker locally and actually run `deploy-rolling.sh` against `lms-pod`; replace the
  illustrative execution log and rollback timing table in `deployment-plan.md` with the real
  output and timestamps.
- Per this week's own syllabus note: security, then CI/CD.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria
- [ ] Tests added/updated (unit / integration) — no automated test for the deploy script itself yet
- [x] DDD / hexagonal boundaries respected (domain has no I/O) — n/a, infra-only deliverable
- [x] No secrets; config via environment variables

## 6. Evidence links
- Local documents attached in this same folder (`12-week/hu-status/`):
  - `deployment-plan.md` — the full MVP 3 deployment plan (artifact promotion, IaC, strategy,
    rollback, sliced stories) and Session 1's execution narrative
  - `deploy-rolling.sh` — the rolling deploy/rollback script itself, not just a description of one
