<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       09-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 09

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Oscar Mauricio Areiza Paramo
- GITHUB_USER: OscarAreiza
- TEAM: lms-library
- SPRINT_GOAL: Session 1 - harden your config: .env.example + startup validation of required vars, secrets never in git, a pre-commit secret scan, and at least one feature flag guarding a new capability. Session 2 (planning) - turn this into a secure-config and rollout plan for MVP 2.
<!-- CONFIG-END -->

> **This week's brief (as given):** Harden your config: .env.example + startup validation of
> required vars, secrets in a store/injected (never in git), a pre-commit secret scan, and at
> least one feature flag guarding a new capability. Session 2 (planning) turns this into a
> secure-config and rollout plan for MVP 2. Write your secrets plan (owners + rotation), a
> feature-flag policy (naming, owner, removal), and a canary + rollback plan for one MVP 2
> feature. Slice the hardening stories with testable acceptance criteria. Next week:
> persistence, then the MVP 2 release.

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-HARDEN-01 | `.env.example` + startup validation of required vars | done | `lms-access-api#7` (new) + `lms-catalog-api`'s existing `config.Load()` (confirmed, no change needed) — `CONFIG-HARDENING.md` §1 |
| HU-HARDEN-02 | Secrets never in git, pre-commit secret scan | doing | https://github.com/code-corhuila/lms-access-api/pull/6 — rolled out to one reference repo, not the fleet — `CONFIG-HARDENING.md` §2 |
| HU-HARDEN-03 | At least one feature flag guarding a shipped capability | done | https://github.com/code-corhuila/lms-catalog-api/pull/9 — verified live against the real shared Postgres with a real signed token — `CONFIG-HARDENING.md` §3 |
| HU-HARDEN-04 | Secrets plan (owners + rotation) | done | `SECURE-CONFIG-ROLLOUT-PLAN.md` §1 |
| HU-HARDEN-05 | Feature-flag policy (naming, owner, removal) | done | `SECURE-CONFIG-ROLLOUT-PLAN.md` §2 |
| HU-HARDEN-06 | Canary + rollback plan for one MVP 2 feature | done | `SECURE-CONFIG-ROLLOUT-PLAN.md` §3 |
| HU-HARDEN-07 | Hardening stories sliced with testable acceptance criteria | done | `SECURE-CONFIG-ROLLOUT-PLAN.md` §4 (HS-01..HS-04) |

## 2. My individual contribution
- **Session 1 - implementation:** fixed a real correctness bug found while wiring up the startup
  validation: `lms-access-api` was issuing Administrator tokens HS256-signed with a shared
  secret, which both contradicts `rules/2-anexos/C-api-hexagonal.md` 5.3.7 ("only the issuer
  holds the private key") and doesn't match what every other `-api` already validates against
  (`JWT_PUBLIC_KEY`, RS256). Fixed alongside the startup-validation work since both trace to the
  same under-exercised config path (`lms-access-api#7`). Built and live-tested a real
  gitleaks-based pre-commit hook (`lms-access-api#6`) with a correctly fingerprint-scoped
  allowlist — proved it actually catches a planted fake secret and doesn't false-positive on a
  known-safe placeholder, not just that the hook file exists. Built and live-tested a real
  feature flag on `lms-catalog-api`'s `GET /books` (`lms-catalog-api#9`), including discovering
  the correct HTTP semantics when the flag is off (`405`, since the path still exists for
  `POST` — not a bare `404` or a `5xx`). This PR replaces an earlier one (`#5`) that got closed
  as superseded once a teammate's hexagonal-restructuring PR landed on `develop` first — re-applied
  cleanly against the new code, re-verified live rather than assumed still correct.
  Full detail -> `CONFIG-HARDENING.md`.
- **Session 2 - planning:** wrote the secrets plan, feature-flag policy, and canary + rollback
  plan grounded in what Session 1 actually built and tested — the canary/rollback plan is
  written against the real, verified `405`/`201` behavior of `FEATURE_CATALOG_SEARCH_ENABLED`,
  not a hypothetical flag. Sliced four hardening stories (HS-01..HS-04) in the same Gherkin
  format `04-requirements/user-stories.md` already uses, scoped to close this week's honestly-stated
  gaps: fleet-wide rollout of both patterns, a CI backstop for the secret scan, and a real
  secrets store once a non-local environment actually exists. Full detail ->
  `SECURE-CONFIG-ROLLOUT-PLAN.md`.

## 3. Blockers and risks
- Both `lms-access-api` PRs (`#6`, `#7`) and `lms-catalog-api#9` are still open and unreviewed as
  of this write-up - nothing here is merged to `develop` yet.
- The secret-scan and feature-flag patterns exist in one repo each; every other `-api`/`-portal`/
  `-db` has neither (`HS-01`, `HS-03` exist to close this, not yet started).
- No CI pipeline runs the secret scan (`HS-02`) - the pre-commit hook is the only gate today, and
  it fails open if `gitleaks` isn't installed locally.
- The secrets plan (`SECURE-CONFIG-ROLLOUT-PLAN.md` §1) is a draft, not yet ratified by the team -
  same caveat prior weeks' first-draft planning docs (`06`-`08-week`) already carried forward.
- `INTERNAL_JWT_SECRET` has no safe rotation window (a single shared secret across three
  services, no old/new overlap possible) - flagged in the secrets plan as a real weakness, not
  glossed over.

## 4. Plan for next week
- Get `lms-access-api#6`, `#7`, and `lms-catalog-api#9` reviewed and merged.
- Persistence, then the MVP 2 release, per the brief - continuing the Anexo J single-shared-
  database-instance work already underway (`lms-infra`/`lms-infra-mongo`) and the `qa`
  promotion of the Circulation and Catalog repos (in progress, separate from this week's
  hardening scope).
- Start `HS-01`/`HS-03`: roll the secret-scan hook and the feature-flag pattern out to a second
  repo each, to prove the pattern generalizes before calling it fleet-wide.
- Take the secrets plan to the team for ratification (same open item `08-week`'s planning docs
  already carried).

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [x] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria - `SECURE-CONFIG-ROLLOUT-PLAN.md` §4, Gherkin format
- [ ] Tests added/updated (unit / integration) - config/auth changes verified live end-to-end
  (real signed tokens, real shared Postgres), not via a new unit test added this week
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables - this week's whole subject

## 6. Evidence links
- Local documents attached in this same folder (`09-week/hu-status/`):
  - `CONFIG-HARDENING.md` - Session 1: real implementation evidence, live-verified
  - `SECURE-CONFIG-ROLLOUT-PLAN.md` - Session 2: secrets plan, feature-flag policy, canary +
    rollback plan, sliced hardening stories
- Open PRs this week: https://github.com/code-corhuila/lms-access-api/pull/6,
  https://github.com/code-corhuila/lms-access-api/pull/7,
  https://github.com/code-corhuila/lms-catalog-api/pull/9
- Latest ADRs referenced: `library-docs/05-architecture/decisions/records/ADR-006-repo-per-context-decomposition.md`,
  `ADR-007-gateway-auth-and-rate-limiting.md`, `ADR-008-circulation-saga-scope.md`,
  `ADR-009-worker-scheduling-model.md`
- Governance referenced: `library-docs/00-governance/git-conventions.md`,
  `library-docs/00-governance/definition-of-done.md`
