# ADR-007 — `lms-api-gateway`'s Auth Boundary and Rate Limiting

- **ID:** ADR-007
- **Date:** 2026-09-15
- **Status:** Accepted
- **Authors:** Oscar Areiza — Tech Lead
- **Reviewers:** Hermes Pascuas, Luis Alejandro Meneses — Development team

---

## Context

`lms-api-gateway`'s own README describes it as "Single entry point: **authentication**, routing
and rate limiting" — a broader mandate than today's NGINX, which `ADR-003` explicitly keeps out
of authn/authz ("does not perform authn/authz"). Read literally, "authentication" at the gateway
could mean the gateway becomes the sole place a JWT is validated, and every backend service stops
validating it locally.

That reading breaks an already-working, already-documented pattern:
`09-microservices/services/02-access-service/decisions.md` and `09-microservices/communication-patterns.md`
both establish that **service-to-service calls never pass through the gateway** — they go
container-to-container over the internal Docker network (e.g. `membership-service`'s
`circulation/client.go` calling `http://circulation-service:8080` directly). If the gateway were
the only thing that validates a JWT, every one of those direct calls would be completely
unauthenticated, a real regression, not a refactor.

Separately, `05-architecture/overview.md` and `05-architecture/security-threat-model.md` already
*claim* NGINX has rate limiting configured (`limit_req`) — verified false this session: no
`limit_req` directive exists in the real `infra/nginx/nginx.conf`. `00-governance/security-rules.md`
independently flags login rate limiting as "a real, un-mitigated gap." Both problems belong to
the same repo (`lms-api-gateway`) and the same decision: what does the gateway actually do.

**Known constraints:**
- Same 3-person team, no budget for a managed API gateway/WAF product (`01-context/scope.md`)
- HS256 with a shared `JWT_SECRET` remains the signing model (`ADR-002/access-service/decisions.md`) —
  this ADR does not revisit that
- No dedicated service-to-service auth scope exists (`ADR-004`'s accepted trade-off) — this ADR
  does not close that gap either, it only decides what the gateway itself is responsible for

---

## Decision

**We decided:** `lms-api-gateway` does **not** replace per-service JWT validation. Every backend
service keeps validating its own JWT locally against the shared `JWT_SECRET`, exactly as today.
The gateway's "authentication" responsibility is limited to a cheap, stateless pass-through check:
reject a request to a protected route with no `Authorization` header at all (fast 401, no JWT
parsing, no trip to a backend service) — a DoS-cheapening filter, not the system's authority on
who is logged in. Real, authoritative validation stays exactly where `ADR-004` already put it.

**Rate limiting is implemented for real**, using NGINX's `limit_req_zone`/`limit_req`:

| Zone | Scope | Rate | Burst |
|---|---|---|---|
| `login` | `POST /api/v1/auth/login` | 5 req/min per IP | 2 |
| `api` | Every other `/api/v1/*` route | 60 req/min per IP | 20 |

These numbers close the gap `00-governance/security-rules.md` already flagged as real and
un-mitigated, and give the whole system (not just login) a basic floor against accidental or
malicious request floods, sized for a single-admin back-office panel, not a public consumer app.

**Justification:** the gateway's job in this system was always "single entry point, routing,
TLS" (`ADR-003`) — adding a cheap reject-if-missing-header filter and real rate limiting extends
that job without contradicting it. Making the gateway the sole authenticator would either (a)
leave direct service-to-service calls unauthenticated, or (b) force every service-to-service call
back through the gateway, which is a much bigger, unrelated architectural change nobody decided to
make. Keeping validation duplicated (gateway does a cheap check, each service does the real one)
costs nothing extra to build — the per-service logic already exists and works.

---

## Evaluated alternatives

| Alternative | Pros | Cons | Reason for discarding |
|------------|------|------|-----------------------|
| **Gateway does a cheap header-presence check only; real validation stays per-service (CHOSEN)** | No change to already-working, already-tested per-service JWT logic; service-to-service calls stay protected exactly as today | The gateway's own README's "authentication" wording overstates what it actually does — must be documented clearly so nobody assumes more than this | — (chosen) |
| Gateway becomes the sole JWT validator; services trust an internal header it sets | Single place to change auth logic; services get simpler | Breaks service-to-service auth entirely (those calls skip the gateway) unless every direct call is rerouted through it — a much larger change than this ADR's scope | Introduces a real security regression as a side effect of a docs/repo-split ADR, not a deliberate, separately-justified decision |
| No rate limiting at all (leave the gap as documented, unaddressed) | Zero work | Login brute-forcing stays a real, named, un-mitigated risk (`security-rules.md`) | The gap is already flagged; this is the moment to close it since `lms-api-gateway`'s own scope statement includes rate limiting |
| Rate limit at the application layer (each Go service) instead of NGINX | Finer-grained, per-endpoint control possible | Duplicates the same logic 3-4 times across services instead of once at the single entry point; NGINX already has a mature, battle-tested `limit_req` directive for exactly this | The gateway is the natural, single chokepoint for this — no reason to reimplement it per service |

---

## Consequences

**Positive:**
- Zero changes to any already-working service's authentication code
- Closes a real, previously-flagged security gap (login rate limiting) with concrete, documented
  numbers instead of leaving it open-ended
- The distinction (gateway = cheap filter, services = real authority) is simple enough for a
  3-person team to reason about without a diagram

**Negative / Trade-offs:**
- JWT validation logic still lives in 3-4 places (once per service) — this ADR does not reduce
  that duplication, it only decides not to make it worse by centralizing incorrectly
- Rate-limit values (5/min login, 60/min general) are reasoned defaults for this project's scale,
  not load-tested — revisit if real traffic patterns disagree

**Impact on the system:**
- Affected repo: `lms-api-gateway` only, once it receives its first commit
- Documents to update once implemented: `lms-infra`'s compose (no change — NGINX config lives in
  `lms-api-gateway`), `00-governance/security-rules.md`'s A07 rate-limiting note (currently
  describes the gap as unmitigated — update once this ADR's values are live, coordinated with the
  currently-open PR touching that same file)

---

## Risks

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|-----------|
| Someone reads `lms-api-gateway`'s README literally and builds real gateway-level JWT validation anyway, duplicating/conflicting with per-service validation | Medium | Medium — inconsistent auth behavior between paths that do and don't go through the gateway | This ADR is the canonical answer; link it from `lms-api-gateway`'s own `decisions.md` once that repo exists |
| Rate-limit values are wrong for real usage (too strict, locks out the single Administrator during legitimate rapid retries; too loose, doesn't actually stop brute-forcing) | Low | Low — single-admin system, easy to adjust | Treat as a starting point, not a permanent constant; adjust based on the demo/real usage, not in the abstract |

---

## References

- Part of the repo decomposition this ADR's target (`lms-api-gateway`) belongs to →
  `ADR-006-repo-per-context-decomposition.md`
- Amends `ADR-003-nginx-reverse-proxy.md` — the gateway's role grows from "routing only" to
  "routing + cheap auth filter + real rate limiting," still not a full auth authority
- Why service-to-service calls bypass the gateway today → `09-microservices/communication-patterns.md`
- The real per-service JWT validation this ADR does not change → `09-microservices/services/02-access-service/decisions.md`
- The previously-flagged, now-closed gap → `00-governance/security-rules.md`, section A07
