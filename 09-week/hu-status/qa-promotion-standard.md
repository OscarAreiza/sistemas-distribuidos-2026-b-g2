# QA Promotion Standard

> **What this is:** the concrete, step-by-step checklist every team member runs **after** a
> `promote-qa/*` PR opens against `qa` (per `00-governance/git-conventions.md`'s
> promotion-by-re-application model) and **before** approving it. It turns the TDD discipline in
> `tdd-guide.md` and the test tiers in `testing-strategy.md` into one shared sequence everyone
> follows the same way, so a review from Hermes, Luis, or Oscar means the same thing.
>
> **Does not replace:** `tdd-guide.md` (how to write a test, Red-Green-Refactor, test doubles) or
> `testing-strategy.md` (what each test tier covers project-wide). This document is "what to
> actually run, in what order, for a `qa` promotion" — the missing operational layer on top of
> both.

---

## Why this exists

Promoting to `qa` is the first point where a story's tests run against something closer to a
real deployment than a laptop — the real inter-service calls (`circulation-api` → `catalog-api`/
`membership-api`), the real shared Mongo instance (Anexo J), the real Module Federation remote
loading. A story that passed `go test ./...` on `develop` can still fail here: a wrong env var, a
missing index, a CORS config that only breaks when two real containers talk to each other. This
is exactly the gap Tier 2 (Integration) and Tier 4 (E2E) in `testing-strategy.md` exist to catch
— this document is where that actually happens, on a schedule, by a named person.

---

## Before opening the `promote-qa/*` PR (author's own responsibility)

Not re-litigated here — already covered on `develop`, this just confirms it didn't regress:

- [ ] Every acceptance criterion the promoted commits implement has a test that was written
      **before** the implementation (TDD — `tdd-guide.md`'s Red-Green-Refactor). If a commit adds
      logic with no corresponding test, that's a gap to close before promoting, not after.
- [ ] `go test ./... -cover` (Go repos) passes locally with no coverage drop vs. the baseline in
      `tdd-guide.md`'s table (domain ≥ 90%, application ≥ 80%, infrastructure ≥ 60%)
- [ ] The PR body's **Test plan** checklist (from `.github/pull_request_template.md`) is filled
      in with real evidence (command output, a screenshot, a curl transcript) — not left as
      unchecked boilerplate

## Stage 1 — Automated, same for every repo

Whoever reviews the `promote-qa/*` PR runs this **before** reading a single line of diff:

| Repo kind | Command | Must show |
|---|---|---|
| Go service (`-api`) | `go vet ./... && go test ./... -cover` | 0 vet issues, all tests green, coverage per `tdd-guide.md`'s table |
| Go migrator (`-db`) | `docker compose -f deploy/compose.yml up --build --abort-on-container-exit` | Liquibase executor exits `0`; no `${VAR}` left unsubstituted in the log (the Anexo J bug class — grep the output for a literal `$`) |
| Frontend (`-portal`) | `npm ci && npm run build && npm run lint` | Clean build, no lint errors, bundle emitted under `dist/` |

**This is a required status check, same spirit as `branching-policy.md`'s CI gate on `develop` —
if it's red, stop here and request changes. Do not move to Stage 2 against red Stage 1.**

## Stage 2 — Integration, against the real `qa` stack

Spin up the promoted service against the **shared** infrastructure it will actually run
against — not an isolated local mock. `lms-infra`/`lms-infra-mongo` run one instance **per
environment** (`env/dev.env`, `env/qa.env`, `env/main.env` — each its own
`COMPOSE_PROJECT_NAME`, its own volume): testing a `qa` promotion against the `dev` instance is a
stand-in, not the real thing — use `env/qa.env` (`docker compose --env-file env/qa.env -p
lms-mongo-qa up -d` / the Postgres equivalent) once that environment's instance needs exercising
for real, not just the one already running locally.

Copy the service's own `.env.example` to a **local, never-committed** file and override only the
connection fields — same variables the course's secrets-hardening rules already forbid putting
real values for in a tracked file (`rules/2-anexos/I-github-common.md`):

```bash
cp .env.example .env.qa-test   # confirm it matches the repo's own `.env*` gitignore pattern before relying on it
# override only:
#   DB_HOST / DB_PORT  -> the shared instance's container name or published port
#   DB_USER / DB_PASSWORD -> the environment's admin creds (or the domain's own
#                            least-privilege user, once its authSource is fixed —
#                            see the Anexo J note below)
#   DB_NAME            -> this domain's database/schema name
#   PORT               -> a free local port if the service's default is already taken
```

```bash
docker run -d --name <service>-qa-test --network platform -p <PORT>:<PORT> \
  -v "$(pwd)":/app -w /app --env-file .env.qa-test golang:1.25 sh -c "go run ./cmd/api"
```

Delete `.env.qa-test` and remove the container the moment the check is done — it holds real
credentials for a real shared environment, and a stopped container is dead weight. Neither one is
evidence; the Stage 2/3 checklist items below, with their output pasted into the PR, are the
evidence.

**Known Anexo J gap (as of this promotion round):** a domain's own least-privilege user (e.g.
`circulation_app`) authenticates against **its own database**, not `admin` — if the service code
being tested still hardcodes `authSource=admin`, root admin credentials are the only way to get a
real connection for Stage 2 until that service's own Anexo J authSource fix lands on `develop`.
Note which one you used in the PR's evidence; don't silently test with root and call it equivalent
to testing with the real least-privilege path.

- [ ] Service starts and `/health` / `/health/ready` return 200 against the **real** shared
      Mongo/Postgres instance (Anexo J) — not a local throwaway container
- [ ] If the service calls another service (e.g. `circulation-api` → `catalog-api`'s
      `/books/{id}/loan-copy`), that call is exercised for real — both services up, real HTTP,
      real HS256 internal token minted — not stubbed
- [ ] Any feature flag the service declares (`.env.example`) is tested in **both** positions
      (on and off), per `09-week/hu-status/FEATURE-FLAG-POLICY.md`

## Stage 3 — Manual acceptance check, against the HU's real Gherkin

Pick the scenarios straight from `04-requirements/user-stories.md` — do not write new ones here,
and do not skip a scenario because "it's basically the same as the other one":

**HU-06 — Loan Registration** (`circulation-api` + `circulation-portal`)
- [ ] Scenario 1: eligible student (< 2 active loans, no suspension) + available copy → loan
      created, due date = start + 7 days, book's available-copy count decreases by one
- [ ] Scenario 2: suspended student, or student already at 2 active loans → loan rejected, with
      the reason stated (`STUDENT_SUSPENDED` / `LOAN_LIMIT_REACHED` → HTTP 422, per
      `690441f`'s error mapping)
- [ ] Retry the same registration with the same `Idempotency-Key` → no duplicate loan, no double
      decrement of available copies (this is the scenario the durable `IdempotencyStore` exists
      for — don't skip it because Scenario 1 already passed)

**HU-07 — Return Registration & History Tracking** (`circulation-api` + `circulation-portal`)
- [ ] Scenario 1: on-time return → loan moves `Active` → `Returned`, available-copy count
      increases by one
- [ ] Scenario 2: late return → return marked late, history records the late status

**HU-08 — Overdue Loans Report** (`circulation-api` + `circulation-portal`)
- [ ] An active loan past its due date, not yet returned, appears in the Overdue Loans report

Each checked box needs a one-line evidence note in the PR (a curl transcript, a screenshot of the
portal, a log line) — "tested, works" with nothing attached does not count, same rule the HU
status README already uses.

## Stage 4 — Sign-off

- Reviewer is **never** the author of the promoted commits — same rule
  `00-governance/git-conventions.md` already sets for `qa` review (minimum 1 approval)
- Approve the PR only after Stages 1–3 are all checked with evidence attached
- If Stage 2 or 3 fails: request changes on the `promote-qa/*` PR itself — fix forward on that
  branch (new commits, or a fresh `cherry-pick -x` of a fix commit from `develop`), never patch
  directly on `qa`

---

## Known gap: the Pact requirement in `branching-policy.md`

`branching-policy.md` (non-negotiable course rule) lists **"contract tests (Pact)"** as a
required status check on `qa`. As of this promotion round, **no repo in this project has a Pact
suite** — `testing-strategy.md`'s own Tier 3 explicitly says contract testing is "not applicable
to this project" (no second consuming team yet, no OpenAPI spec written per service).

Both documents are correct about their own scope and still disagree with each other. Until one
of them is updated, this standard's interim rule — the `00-governance/definition-of-done.md`
"Allowed exceptions" clause requires exactly this: an explicit, written, Tech-Lead-agreed
exception, not a silently skipped gate:

- **Interim substitute for Stage 1 on `qa`:** a reviewer manually checks the promoted service's
  real request/response shapes against `07-api/contracts/openapi/` by hand (the PR's Stage 2
  evidence already has real curl transcripts — diff those against the spec)
- **Action item, not yet assigned:** write a minimal Pact suite for at least one producer/consumer
  pair (`circulation-api` as consumer of `catalog-api`'s `/books/{id}/loan-copy` is the obvious
  first candidate — it's already a real cross-service call) and decide, as a team, whether to
  update `branching-policy.md`'s wording or actually build it before the next MVP's `qa`
  promotion round

This gap is declared here, in writing, rather than assumed closed — don't check Stage 1's Pact box
on any PR until this is resolved for real.

---

## Correlations

- Red-Green-Refactor, test doubles, FIRST principles → `11-quality/tdd-guide.md`
- Test tiers and what each one covers → `11-quality/testing-strategy.md`
- Promotion mechanics (`cherry-pick -x`, branch naming) → `00-governance/git-conventions.md`
- The non-negotiable gates per branch, incl. the Pact requirement → `00-governance/branching-policy.md`
- Required checklist fields on every PR → `00-governance/definition-of-done.md`
- HU-06/07/08's real acceptance criteria → `04-requirements/user-stories.md`
- Anexo J (shared DB instance) and the bug classes it introduced → `rules/2-anexos/3-anexo-j.md`
