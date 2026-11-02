<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       14-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 14

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Oscar Mauricio Areiza Paramo
- GITHUB_USER: OscarAreiza
- TEAM: lms-library
- SPRINT_GOAL: Set up CI/CD for the LMS service repos - blocking PR gates (build, lint, type, tests, coverage, security), a protected main and automated deployment of the tagged image (Session 1) - then map every Definition of Done item to a gate, define branch protection for develop/qa/main, path-filtered pipelines, cross-service contract tests and a smoke test, and slice the CI/CD stories for MVP 3 (Session 2).
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| W14-S1 | Set up CI/CD: blocking PR gates (build, lint, type, tests, coverage, security), protect main, automate deployment of the tagged image | doing | `CICD-PIPELINE.md` (this folder) - current state audited and target `ci.yml` / `cd.yml` written; not yet applied to the service repos |
| W14-S2 | Map every DoD item to a blocking gate, branch protection on develop/qa/main, path-filtered pipelines, contract tests, smoke test, CI/CD stories | done | `DOD-GATES.md` (this folder) - not yet committed/pushed |

## 2. My individual contribution
- Audited the real pipelines:
  - only the 3 Go APIs and 3 portals have a `ci.yml` (gofmt/vet/test, or npm build);
  - 7 repos have no workflow at all;
  - branch protection has no required status checks on any branch, and only `main` requires a
    review, so CI is advisory today.
- Wrote the target pipeline:
  - four blocking PR jobs (build, lint with golangci-lint, tests with a coverage gate, security
    from week 13);
  - a protected `main`;
  - a tag-triggered CD that builds one image, pushes it to GHCR and deploys it with the week 12
    `deploy-rolling.sh` behind a `production` approval, followed by a smoke test.
- Mapped every item of the team's own Definition of Done
  (`library-docs/00-governance/definition-of-done.md`) to a blocking gate, or to a manual PR-template
  check where a machine cannot verify it.
- Defined:
  - branch protection settings for `develop`/`qa`/`main`;
  - path filters (docs-only PRs in each repo, per-service re-deploy in `lms-pod`);
  - OpenAPI-based contract tests for the only synchronous calls (circulation -> membership and
    catalog);
  - a `smoke-test.sh`;
  - 9 CI/CD stories (CI-01..CI-09) with testable acceptance criteria.

## 3. Blockers and risks
- There is no deploy host reachable from GitHub Actions yet, so the CD `deploy` job cannot run end
  to end; until then the image is built by CI and deployed by hand with the same script and tag.
- There is no container registry in use yet (ADR-006 gap). GHCR is proposed, and
  `deploy-rolling.sh` must change from "build from git tag" to "pull image tag".
- Changing branch protection on the `code-corhuila` repos needs admin rights on the organization.
- The coverage baseline (`COVERAGE_MIN`) is a placeholder until the first run on each repo.

## 4. Plan for next week
- Testing at scale, then the MVP 3 release.
- Apply CI-01 and CI-03 first (gates + branch protection) so that every later PR for MVP 3 already
  goes through them.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria
- [ ] Tests added/updated (unit / integration)
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- Session documents in this same folder (`14-week/hu-status/`):
  - `CICD-PIPELINE.md` - Session 1: current state of CI and branch protection, target `ci.yml` with
    blocking gates, protected main, tag-triggered build and deploy
  - `DOD-GATES.md` - Session 2: DoD-to-gate map, branch protection for develop/qa/main, path filters,
    contract tests, smoke test, CI/CD stories for MVP 3
- Existing workflows reviewed: `lms-circulation-api/.github/workflows/ci.yml` (same in catalog and
  membership), `lms-circulation-portal/.github/workflows/ci.yml`
- Team Definition of Done: `library-docs/00-governance/definition-of-done.md`
