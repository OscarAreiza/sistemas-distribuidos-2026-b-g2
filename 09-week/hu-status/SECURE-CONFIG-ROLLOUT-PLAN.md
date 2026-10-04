# Session 2 — Secure-Config and Rollout Plan for MVP 2

> Brief: "Write your secrets plan (owners + rotation), a feature-flag policy (naming, owner,
> removal), and a canary + rollback plan for one MVP 2 feature. Slice the hardening stories
> with testable acceptance criteria."

Builds directly on Session 1's real evidence (`CONFIG-HARDENING.md`, this same folder) — this
is the planning layer on top of what was actually implemented and tested this week, not a
separate exercise.

---

## 1. Secrets plan — owners + rotation

**Status: draft, not yet ratified by the team** — same caveat prior weeks' first-draft planning
docs (`06-week`, `07-week`, `08-week`) already carried forward. Proposed, grounded in who
actually built/owns each service so far:

| Secret | Lives in | Proposed owner | Rotation trigger | Rotation mechanism |
|---|---|---|---|---|
| `JWT_PRIVATE_KEY` (RS256, signs every Administrator token) | `lms-access-api` only — no other service holds it (`rules/2-anexos/C-api-hexagonal.md` 5.3.7) | Owner of `lms-access-api` | Suspected compromise; a departing team member had access; every 90 days as a baseline | Generate a new RSA pair, update `JWT_PUBLIC_KEY` on every consuming `-api` **before** deploying the new private key — old tokens signed with the retired key must still validate until they expire (`JWT_EXPIRY`, currently 1h), so this is a sequenced rollout, not an atomic swap |
| `INTERNAL_JWT_SECRET` (HS256, service-to-service calls) | `lms-circulation-api` (mints), `lms-catalog-api`/`lms-membership-api` (validate) | Owner of `lms-circulation-api` (the minting side) | Same triggers as above | Must be updated on **all three** services in the same deploy window — unlike the RS256 pair, there's no old/new overlap period possible with a single shared secret; this is this secret's real weakness vs. the RS256 pair above, and an argument for migrating it to asymmetric too in a future MVP |
| `DB_PASSWORD` per domain (`circulation_app`, `catalog_app`, `membership_app`, ...) | The shared infra instance (`lms-infra`/`lms-infra-mongo`) creates the login; the domain's own `-db` migration grants it a role | Owner of the shared infra repo for creation; the domain's `-db` owner for the grant | Suspected compromise; infra redeploy | Update the password in the infra instance's `env/<environment>.env` (never committed — `rules/3-Anexo-J` J.5.1) and in the consuming `-api`'s own secret store in the same window |

**Where secrets actually live today:** local `.env` files (gitignored) for development,
`.env.example` files (committed, names only — no values) documenting what each service needs.
**Not yet implemented:** a real secrets store (Vault, AWS Secrets Manager, or even encrypted
`.env` injection at deploy time) for anything beyond local dev — there is no staging/production
deploy target yet (`05-architecture/deployment.md`), so this is correctly scoped as a plan, not
a built system.

## 2. Feature-flag policy

**Naming:** `FEATURE_<DOMAIN>_<CAPABILITY>_ENABLED`, all caps, matching the env-var convention
every service already uses (`FEATURE_CATALOG_SEARCH_ENABLED` is the first real instance). Boolean
only — a flag that needs more than on/off (a percentage rollout, a variant) is a different
mechanism, not this one, and should say so explicitly when introduced.

**Owner:** whoever's PR introduces the flag owns it until removal — recorded in that PR's
description (already the case for `FEATURE_CATALOG_SEARCH_ENABLED`, owned by its author). A flag
with no PR-traceable owner doesn't get merged.

**Default:** on (`true`) for a flag guarding an already-shipped capability (this is a kill switch,
not a rollout gate) — off (`false`) only for a capability not yet ready to be universally live.

**Removal criteria — a flag is deleted, not just set permanently true/false, once:**
1. The capability has run at its final on/off state for at least one full sprint with no
   incident tied to it, **and**
2. Every environment (`develop`, `qa`, `main`) is confirmed at that same state — no drift, **and**
3. The removal PR deletes the config field, the `if` branch, and the `.env.example` line in the
   same change — a flag that's "always true" but still branches in code is tech debt, not a
   finished rollout, and gets logged in `library-docs/15-project-control/technical-backlog.md`
   if it can't be removed immediately for a documented reason.

**Registry (so a flag doesn't get forgotten once merged):**

| Flag | Service | Introduced | Default | Owner | Status |
|---|---|---|---|---|---|
| `FEATURE_CATALOG_SEARCH_ENABLED` | `lms-catalog-api` | Week 9 (`lms-catalog-api#9`) | `true` | PR author | Live, not yet a removal candidate |

## 3. Canary + rollback plan — `FEATURE_CATALOG_SEARCH_ENABLED`

Chosen because it's the one real flag that exists and was actually tested this week (not a
hypothetical) — the plan below is written against real, verified behavior, not assumed behavior.

**Canary step:** there is no multi-instance/staged deploy target yet (`05-architecture/deployment.md`
— single environment per branch tier, confirmed directly: `lms-infra`/`lms-infra-mongo` each run
one instance per environment, `dev`/`qa`/`main`, not a canary slice within one). The realistic
canary for this project's actual infrastructure is **environment-sequenced, not
traffic-sequenced**: land the flag's code in `develop` (done), exercise it in `qa` (next —
promotion pending, see `CONFIG-HARDENING.md` Gaps), only then `main`. Each stage is the canary for
the next.

**Rollout:**
1. `develop`: flag defaults `true` (ships on) — PR `lms-catalog-api#9`, merged once reviewed.
2. `qa`: promote by re-application (`git cherry-pick -x`, per `library-docs/00-governance/git-conventions.md`
   — never a direct merge). Re-run the exact live check already done in `develop`: flag `false` →
   `GET /books` is `405`, `POST /books` still `201`; flag `true` → `GET /books` is `200`.
3. `main`: only after `qa` passes with no incident, via a `release/*` branch per
   `branching-policy.md`.

**Rollback trigger:** `GET /books` returning `5xx` at a rate above baseline, or a bad query
pattern visibly hammering the shared Postgres instance (the exact scenario the flag exists for,
per its own code comment).

**Rollback action:** set `FEATURE_CATALOG_SEARCH_ENABLED=false` in that environment's config and
restart the service — **no code deploy required**, no database change, nothing to migrate back.
This is the entire reason the capability is gated by a route-registration `if`, not a deeper
change: the rollback is a config flip, verified to complete in the time it takes the process to
restart.

**Rollback verification:** confirm `GET /books` now returns `405` (not just "stopped erroring") —
a flag that silently no-ops instead of cleanly disabling the route would be the wrong
implementation; this is exactly what was verified live in Session 1.

## 4. Hardening stories, sliced with testable acceptance criteria

The honest gap from `CONFIG-HARDENING.md`: the two patterns built this week (pre-commit secret
scan, feature flag) exist in one reference repo each. These are the slices to close that,
written at the same granularity as `04-requirements/user-stories.md`'s HU-06/07/08:

**HS-01 — Roll out the pre-commit secret scan to every `-api`/`-db` repo**
```gherkin
Scenario: A repo without the hook gets it
  Given a repo that has lms-access-api's .githooks/pre-commit and .gitleaksignore pattern, but
        hasn't adopted it yet
  When  the hook is added and `git config core.hooksPath .githooks` is documented in its README
  Then  committing a fake secret (e.g. a Stripe-shaped test token) is blocked locally
  And   a known-safe placeholder already in that repo doesn't trigger a false positive
```

**HS-02 — CI secret-scanning backstop**
```gherkin
Scenario: A secret bypasses the local hook
  Given a contributor without gitleaks on PATH, or a commit made from a tool that skips hooks
  When  that commit reaches a Pull Request
  Then  a CI job runs the same gitleaks scan against the diff
  And   the PR is blocked from merging if it finds a real secret
```

**HS-03 — Roll out the feature-flag pattern to one more real capability**
```gherkin
Scenario: A second service gets a flag on a shipped capability
  Given a capability already live in one more -api (e.g. lms-circulation-api's overdue-loans
        report, HU-08)
  When  a FEATURE_CIRCULATION_OVERDUE_REPORT_ENABLED flag is added following §2's naming/default
        rules
  Then  turning it off returns the documented HTTP status for that route (405 if the method
        simply isn't registered, consistent with §3's verified pattern)
  And   the flag is entered in §2's registry table in the same PR
```

**HS-04 — A real secrets store for anything beyond local dev**
```gherkin
Scenario: A secret is rotated without touching a committed file
  Given JWT_PRIVATE_KEY currently lives only in a local, gitignored .env
  When  a secrets store (Vault, AWS Secrets Manager, or equivalent) is introduced for the first
        non-local environment (qa, once one is actually deployed somewhere real)
  Then  that environment's service reads the secret from the store at startup, not from a file
        in its own repo or image
  And   rotating it there requires no code change and no redeploy of the service's image
```

**Not sliced into this week — explicitly deferred, per the brief's own "next week: persistence,
then the MVP 2 release":** none of these four stories are blocked by persistence work, but none
of them are claimed as started either; this section is the backlog this week produces, not work
already done.

## Correlations

- Real implementation evidence for Session 1 → `CONFIG-HARDENING.md`, this same folder
- Branching/promotion mechanics referenced in §3 → `library-docs/00-governance/git-conventions.md`
- Technical debt tracking for an un-removable flag → `library-docs/15-project-control/technical-backlog.md`
- HU-06/07/08's own acceptance-criteria format, mirrored in §4 → `library-docs/04-requirements/user-stories.md`
