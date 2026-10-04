# Session 1 — Config Hardening (Week 9)

> Brief: ".env.example + startup validation of required vars, secrets in a store/injected (never
> in git), a pre-commit secret scan, and at least one feature flag guarding a new capability."

Scope: the two services that most directly needed it — `lms-access-api` (mints every
Administrator token in the system; its private key is the single most sensitive secret in the
project) and `lms-catalog-api` (chosen for the feature flag because `GET /books` is real,
shipped capability, not a placeholder — gating a stub would prove nothing).

## 1. `.env.example` + startup validation of required vars

Both services already load config exclusively from environment variables
(`library-docs/05-architecture/cross-cutting.md`); this session closed the two real gaps:

- **`lms-access-api` had no `.env.example` at all**, and its `config.Load()` booted successfully
  with an **empty JWT signing secret** if `JWT_SECRET` was unset — a silent misconfiguration, not
  a startup failure. Fixed in
  [lms-access-api#7](https://github.com/code-corhuila/lms-access-api/pull/7):
  `config.Load()` now returns `"JWT_PRIVATE_KEY must be set"` and the process never binds a port
  if it's missing or unparseable. (This PR also fixes a real correctness bug found while doing
  this: the service was signing Administrator tokens **HS256** with a shared secret, when
  `rules/2-anexos/C-api-hexagonal.md` 5.3.7 and every other `-api`'s `JWT_PUBLIC_KEY`
  validation already assume **RS256** with a private/public key pair only this service holds.
  Both fixes share one root cause — the config layer was never actually exercised against the
  real cross-service auth contract — so they went in the same PR.)
- **`lms-catalog-api` already had `.env.example` and startup validation** for `JWT_PUBLIC_KEY`/
  `INTERNAL_JWT_SECRET` (fails fast if either is empty or unparseable) — confirmed by reading
  `internal/config/config.go`, no change needed there.

## 2. Secrets never in git, pre-commit secret scan

[lms-access-api#6](https://github.com/code-corhuila/lms-access-api/pull/6) — real tool, not a
hand-rolled regex scanner:

- **gitleaks** (v8.30.1) run via a tracked `.githooks/pre-commit` hook, activated per clone with
  `git config core.hooksPath .githooks`. Scans only staged changes (`gitleaks protect --staged`),
  so it stays fast on every commit.
- **Fails open** if `gitleaks` isn't on `PATH` (skips the scan, doesn't block the commit) — a
  missing local scanner blocking every contributor's every commit is a worse failure mode than
  one missed local-only scan; CI is meant to be the real backstop (not wired yet — see Gaps,
  below).
- **`.gitleaksignore`** allowlists exactly one fingerprint — the throwaway dev RSA keypair
  documented in `.env.example` (`.env.example:private-key:16`), by fingerprint, not a blanket
  file/path ignore. Confirmed directly that a *different*, newly-introduced secret anywhere in
  that same file still gets caught: tested with a fake Stripe-shaped token
  (`sk_live_51H8xK2...`), correctly flagged under rule `stripe-access-token`.
- Verified against AWS's own canonical placeholder (`AKIAIOSFODNN7EXAMPLE`) **not** tripping a
  false positive — confirms the scanner's stopword list behaves as documented, not a bug in this
  setup.

## 3. At least one feature flag guarding a new capability

[lms-catalog-api#9](https://github.com/code-corhuila/lms-catalog-api/pull/9) —
`FEATURE_CATALOG_SEARCH_ENABLED` (default `true`) gates **route registration** for `GET /books`,
not just request handling:

```go
if cfg.FeatureCatalogSearchEnabled {
    books.Get("/", cfg.Books.List) // HU-04, search half
}
```

Verified live against the real shared dev Postgres instance, with a real RS256-signed
Administrator token (not just read from the diff):

| Flag | `GET /books` | `POST /books` |
|---|---|---|
| `false` | `405 Method Not Allowed` (the path exists for `POST`, not a `5xx` from a half-wired handler) | `201 Created` — real row inserted |
| `true` (default) | `200 OK` with real data | unaffected |

Naming, ownership, and removal criteria for this and future flags → `SECURE-CONFIG-ROLLOUT-PLAN.md`
§2 (Feature-Flag Policy), this same folder.

## Gaps — stated honestly, not hidden

- **Rolled out to one reference repo each, not the whole fleet.** The gitleaks hook exists only
  in `lms-access-api`; the feature-flag pattern only in `lms-catalog-api`. Every other `-api`/
  `-portal`/`-db` repo has neither yet.
- **No CI secret-scanning backstop.** The pre-commit hook is the only gate today; a contributor
  without `gitleaks` on `PATH`, or who commits from a tool that bypasses hooks, isn't caught.
- **Both `lms-access-api` PRs (#6, #7) are still open**, unreviewed as of this write-up — not
  "done," "in review."

## Correlations

- Secrets ownership + rotation → `SECURE-CONFIG-ROLLOUT-PLAN.md` §1
- Feature-flag policy → `SECURE-CONFIG-ROLLOUT-PLAN.md` §2
- Canary + rollback plan for an MVP 2 feature → `SECURE-CONFIG-ROLLOUT-PLAN.md` §3
- Hardening stories sliced with acceptance criteria → `SECURE-CONFIG-ROLLOUT-PLAN.md` §4
