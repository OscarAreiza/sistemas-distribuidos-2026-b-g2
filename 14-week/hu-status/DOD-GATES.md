# Week 14 — Session 2: Mapping the DoD to pipeline gates, and CI/CD stories for MVP 3

Source of the DoD: `library-docs/00-governance/definition-of-done.md` (the team's own checklist).
Pipeline design: `CICD-PIPELINE.md`.

## 1. Every DoD item mapped to a gate

**Blocking** = the PR cannot merge (required status check or required review). **Manual** = cannot
be checked by a machine, so it has a checkbox in the PR template that the reviewer verifies.

| DoD item | Gate | Type |
|---|---|---|
| Code implements all acceptance criteria | Tests written from the criteria (`test` job), and the reviewer checks them against the story | Blocking + manual |
| Reviewed and approved by at least 1 team member | Branch protection: 1 required approval | Blocking |
| Linting and formatting pass in CI | `lint` job (`golangci-lint` / `oxlint` / `eslint`) | Blocking |
| No technical debt without registering it | PR template checkbox and link to `technical-backlog.md` | Manual |
| Unit tests for new business logic | `test` job; coverage gate stops new code without tests | Blocking |
| Coverage does not decrease from baseline | Coverage gate (`COVERAGE_MIN`, never lowered) | Blocking |
| All tests pass locally and in CI | `test` job | Blocking |
| Acceptance criteria verified | `test` job + smoke test after deploy to staging | Blocking |
| Does not break other services (integration) | `contract` job (section 4) | Blocking |
| API changes → OpenAPI contract updated | `contract` job validates the service's responses against the OpenAPI file in `07-api/contracts/` | Blocking |
| Data model changes → `data-model.md` updated | PR template checkbox; path filter warns when `migrations/` changes without a docs link | Manual |
| New/modified events → `event-catalog.md` updated | PR template checkbox | Manual |
| Mergeable to `develop` (no conflicts) | Branch protection: branch must be up to date | Blocking |
| CI/CD green on the branch | All jobs are required status checks | Blocking |
| Deployed to staging | CD job for `qa` tags deploys to staging | Blocking for the `qa` → `main` promotion |
| Smoke test passing on staging | `smoke-test.sh` runs after every deploy; a failure fails the deploy job | Blocking for the promotion |
| README updated if the public interface changed | PR template checkbox | Manual |
| ADR for significant decisions | PR template checkbox | Manual |
| *(from week 13)* No committed secrets, no vulnerable deps, authz tests | `security` job | Blocking |

## 2. Branch protection on `develop`, `qa` and `main`

Same settings in every service repo. Today none of the branches has required checks, and only
`main` requires a review.

| Setting | `develop` | `qa` | `main` |
|---|---|---|---|
| Require a PR before merging | yes | yes | yes |
| Required approvals | 1 | 1 | 1 |
| Dismiss stale approvals on new commits | yes | yes | yes |
| Required status checks | `build`, `lint`, `test`, `security`, `contract` | same | same |
| Branch must be up to date before merging | yes | yes | yes |
| Block force pushes and deletion | yes | yes | yes |
| Apply to administrators | no (instructor/maintainer may fix a blocked branch) | no | **yes** |

Applied with the GitHub API, so it is the same in every repo and can be checked:

```bash
for repo in lms-catalog-api lms-circulation-api lms-membership-api lms-access-api; do
  for b in develop qa main; do
    gh api -X PUT repos/code-corhuila/$repo/branches/$b/protection --input protection-$b.json
  done
done
```

## 3. Path-filtered pipelines per service

Each bounded context has its own repo (ADR-006), so each service already has its own pipeline.
Path filters are still needed in two places:

- **In each service repo**, documentation-only changes skip the heavy jobs:

  ```yaml
  on:
    pull_request:
      branches: [develop, qa, main]
      paths-ignore: ["**/*.md", "docs/**"]
  ```

  The required checks still have to report. For those PRs, a small `noop` workflow with the same
  job names and the opposite filter reports success.

- **In `lms-pod`** (one `compose.yml` for all services), a change only re-deploys and smoke-tests
  the services whose block or `env/` file changed:

  ```yaml
  jobs:
    changes:
      runs-on: ubuntu-latest
      outputs:
        services: ${{ steps.filter.outputs.changes }}
      steps:
        - uses: actions/checkout@v4
        - uses: dorny/paths-filter@v3
          id: filter
          with:
            filters: |
              circulation: ['env/circulation*', 'compose.yml']
              catalog:     ['env/catalog*', 'compose.yml']
              membership:  ['env/membership*', 'compose.yml']
              access:      ['env/access*', 'compose.yml']
    smoke:
      needs: changes
      if: needs.changes.outputs.services != '[]'
      strategy:
        matrix:
          service: ${{ fromJSON(needs.changes.outputs.services) }}
      runs-on: ubuntu-latest
      steps:
        - run: echo "deploy + smoke ${{ matrix.service }}"
  ```

## 4. Cross-service contract tests

The only synchronous calls between services are **circulation → membership** (`GET
/students/{id}`) and **circulation → catalog** (`POST /books/{id}/loan-copy`). These are the
contracts that break HU-06 when one side changes.

| Side | What runs | Where |
|---|---|---|
| **Provider** (catalog, membership) | Starts the service with a test DB, calls each endpoint in its OpenAPI file (`07-api/contracts/<service>.yaml`), and validates every response (status, body schema) against that file. | `contract` job in the provider's `ci.yml` |
| **Consumer** (circulation) | `catalogclient` and `membershipclient` tests run against a stub server generated from the **same** OpenAPI files, including the error cases circulation depends on (`404`, `409 NO_COPIES_AVAILABLE`, suspended student). | `contract` job in circulation's `ci.yml` |

The OpenAPI files are the single source of truth. A provider change that breaks a response fails
the provider's CI, and a consumer that expects a field the provider does not have fails the
consumer's CI. Both happen before merge.

## 5. Smoke test

Runs after every deploy (staging and production). If it fails, the deploy job fails, and the
rollback from week 12 applies.

```bash
#!/usr/bin/env bash
# smoke-test.sh — minimal end-to-end check of a deployed lms-pod
set -euo pipefail
BASE=${BASE:-http://localhost}

curl -fsS "$BASE/health" >/dev/null                         # gateway up
for s in access membership catalog circulation; do           # each API ready (internal)
  docker compose exec -T "$s-service" wget -qO- http://localhost:8080/health/ready >/dev/null
done

TOKEN=$(curl -fsS -X POST "$BASE/api/v1/auth/login" -H 'Content-Type: application/json' \
  -d "{\"username\":\"$SMOKE_USER\",\"password\":\"$SMOKE_PASSWORD\"}" | jq -r .token)

curl -fsS -H "Authorization: Bearer $TOKEN" "$BASE/api/v1/books?limit=1" >/dev/null
curl -fsS -H "Authorization: Bearer $TOKEN" "$BASE/api/v1/students?limit=1" >/dev/null
curl -fsS -H "Authorization: Bearer $TOKEN" "$BASE/api/v1/loans?limit=1" >/dev/null
echo "smoke OK"
```

`SMOKE_USER` / `SMOKE_PASSWORD` belong to a dedicated read-only smoke-test Administrator. They are
stored as CI secrets, never in the repo.

## 6. CI/CD stories for MVP 3 (sliced, with testable acceptance criteria)

| ID | Story | Acceptance criteria |
|---|---|---|
| CI-01 | As the team, we want the four PR gates in every Go API, so that broken code cannot merge. | A PR with a `gofmt` error, a failing test, or coverage below `COVERAGE_MIN` shows a red required check and the Merge button is disabled. A clean PR shows 4 green checks. |
| CI-02 | As the team, we want the same gates in `lms-access-api` and the portals, so that every repo follows one DoD. | Each of the 7 repos without a workflow has a `ci.yml`. An `oxlint`/`eslint` error in a PR blocks the merge. |
| CI-03 | As the team, we want branch protection on `develop`, `qa` and `main`, so that the DoD is enforced, not optional. | `gh api repos/code-corhuila/<repo>/branches/<b>/protection` returns the settings in section 2 for every service repo and branch. A direct `git push` to `main` is rejected. |
| CI-04 | As the team, we want security gates in CI, so that secrets and vulnerable deps never reach `main`. | A test PR with a fake AWS key fails `gitleaks`. A test PR that adds a dependency with a known critical CVE fails `govulncheck` / `npm audit`. |
| CI-05 | As the Circulation owner, I want contract tests between circulation and catalog/membership, so that HU-06 does not break when a provider changes. | Renaming a response field in catalog's `loan-copy` handler without updating the OpenAPI file fails catalog's `contract` job. Circulation's `contract` job passes against the stub for `200`, `404` and `409`. |
| CI-06 | As the team, we want each tag to build and push one image, so that the same artifact is promoted through environments. | Pushing tag `v3.0.0` produces `ghcr.io/code-corhuila/<repo>:v3.0.0`. Deploying to staging and production uses that exact image digest. |
| CI-07 | As the team, we want the tagged image deployed automatically with a rolling update, so that releases need no manual steps. | After approval in the `production` environment, the deploy job runs `deploy-rolling.sh` with the tag, and `GET /health` returns `200` during the whole deploy (no failed request in a 1 req/s probe). |
| CI-08 | As the team, we want a smoke test after every deploy, so that a broken release is caught at once. | `smoke-test.sh` exits `0` on a healthy pod. Stopping `catalog-service` makes it exit non-zero and marks the deploy job as failed. |
| CI-09 | As the team, we want path filters, so that docs-only PRs are fast and `lms-pod` only re-deploys what changed. | A PR that only changes `README.md` finishes in under 1 minute with all required checks green. A change to `env/catalog.env` in `lms-pod` runs the smoke job only for `catalog`. |

## Next week

Testing at scale, then the MVP 3 release.
