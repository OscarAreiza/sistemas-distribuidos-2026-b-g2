<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       15-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 15

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Oscar Mauricio Areiza Paramo
- GITHUB_USER: OscarAreiza
- TEAM: lms-library
- SPRINT_GOAL: Close the v3.0.0 release gaps found while packaging the one-click local installer — a missing PostgreSQL role set for the 3 newest services (lms-iam, lms-workflow, lms-worker) and a broken Dockerfile — verify the fix against a real, running stack, and ship the installer as a release package.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| W15-S1 | Extend the one-click installer to provision `lms-iam`, `lms-workflow`, `lms-worker` (compose services, env files, Postgres app roles) | done | `release-v3.0.0.zip` (this folder); `library-docs` commits to `10-devops/installer/install-lms.cjs` on `main` |
| W15-S2 | Consolidate 9 stale, mutually-conflicting `lms-infra-postgres` PRs into one clean `access_app`/`workflow_app`/`iam_app` role + compose-include change | done | https://github.com/code-corhuila/lms-infra-postgres/pull/14 |
| W15-S3 | Find and fix a real build defect (`lms-iam`'s Dockerfile pinned to `golang:1.25-alpine`, but `go.mod` requires `1.26.1`) by actually building and running the container, not by static review | done | https://github.com/OscarAreiza/lms-iam/pull/8 |

## 2. My individual contribution
- Noticed the installer never grew with the project: `lms-iam`, `lms-workflow`, `lms-worker` existed
  as real services but were missing from `REPOS`, `lms-pod/compose.yml`, and the generated env
  files — the installer would have stood up an incomplete system.
- Traced *why* those three couldn't work yet: `lms-infra-postgres`'s `develop` never actually got the
  `access_app`/`workflow_app`/`iam_app` roles. Found 9 open PRs (#5–#13) that each added a slice of
  this independently, off the same stale base, conflicting with each other — none ever merged.
  Wrote one clean, non-conflicting change instead of trying to merge nine broken ones (PR #14).
- Ran the actual verification instead of trusting the YAML: brought up a fresh PostgreSQL instance
  from this branch, ran all 5 migration executors (`membership`, `catalog`, `access`, `workflow`,
  `iam`) to completion (exit 0 each), confirmed the 5 app roles exist in the database, then built and
  ran the real `lms-iam`/`lms-workflow`/`lms-worker` containers against it.
- That real run caught a defect static review missed: `lms-iam`'s `docker build` failed outright
  (`go.mod requires go >= 1.26.1`, Dockerfile builds with `golang:1.25-alpine`). Fixed and verified:
  rebuilt, `GET /health/ready` answered `{"status":"ok","database":"connected"}`.
- Packaged the installer (`install-lms.cjs`, recompiled `lms-installer.exe`, `README.md`, and a new
  `RELEASE-NOTES-v3.0.0.md`) into `release-v3.0.0.zip`, attached to this folder as evidence.

## 3. Blockers and risks
- `lms-infra-postgres#14` and `lms-iam#8` are open, not yet merged — anyone cloning straight from
  `code-corhuila`/`OscarAreiza` today still hits the same role-missing/build-broken state; everything
  above was verified against the locally-patched branches, not the merged remote.
- `lms-iam`'s `/authorize` endpoint still has no real caller wired into any service (a separate,
  already-tracked gap from `ADR-012`) — standing the service up does not by itself turn on
  authorization.
- No Go toolchain in this environment for most of this session's other v3.0.0 work, so several
  sibling PRs (new `lms-circulation-api`/`lms-catalog-api` endpoints) are unverified by `go vet`/
  `go test` — flagged explicitly in their own PR descriptions, not hidden.

## 4. Plan for next week
- Get `lms-infra-postgres#14` and `lms-iam#8` reviewed and merged, then re-run the same verification
  straight from a clean clone of the real repos (no local patches) to confirm the installer works
  end to end with zero manual intervention.
- Wire at least one real caller into `lms-iam`'s `/authorize` endpoint.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...) — these are infra/tooling fixes, not a domain HU, so they went straight to `develop`/`main` per each repo's own convention for that kind of change
- [x] Testable acceptance criteria — "the 5 migrations exit 0, the 5 roles exist, `lms-iam` answers `database:connected`" is exactly what was checked
- [ ] Tests added/updated (unit / integration) — this was infra/config + a real manual run, not new application logic
- [x] DDD / hexagonal boundaries respected (domain has no I/O) — n/a, infra-only deliverable
- [x] No secrets; config via environment variables — `env/*.env.example` pattern followed for the 3 new passwords

## 6. Evidence links
- `15-week/hu-status/release-v3.0.0.zip` — the packaged installer (source, compiled `.exe`, README, release notes)
- https://github.com/code-corhuila/lms-infra-postgres/pull/14 — consolidated role/compose fix
- https://github.com/OscarAreiza/lms-iam/pull/8 — Dockerfile Go-version fix
- `library-docs/10-devops/installer/install-lms.cjs`, `RELEASE-NOTES-v3.0.0.md` (commits on `main`)
