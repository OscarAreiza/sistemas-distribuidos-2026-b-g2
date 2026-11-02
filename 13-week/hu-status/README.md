<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       13-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 13

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Oscar Mauricio Areiza Paramo
- GITHUB_USER: OscarAreiza
- TEAM: lms-library
- SPRINT_GOAL: Secure the LMS system - authenticate at the edge, validate tokens and authorize in each service, secure service-to-service calls, hash secrets and validate inputs, checked against the real code of each repo (Session 1) - then build a STRIDE threat model with trust boundaries, a prioritized security backlog with testable criteria, and CI security gates for MVP 3 (Session 2).
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| W13-S1 | Secure the system: edge auth, token validation and authorization per service, service-to-service security, hashed secrets, input validation | doing | `SECURITY-CONTROLS.md` (this folder) - audit of each control in the code is done; edge auth is on an unmerged gateway branch and per-service authorization does not exist yet (SEC-02, SEC-03, SEC-04) |
| W13-S2 | Threat model (data flow + trust boundaries + STRIDE), prioritized security backlog with testable criteria, CI security gates | done | `THREAT-MODEL.md` (this folder) - not yet committed/pushed |

## 2. My individual contribution
- Checked every control in the real code instead of describing it in general:
  - the `RequireAuth` middleware of the Go APIs (RS256 for Administrator sessions and HS256 for
    internal calls, each with its own key, chosen by the token's algorithm to block algorithm
    confusion);
  - the internal 1-minute service token minted by circulation's clients;
  - bcrypt cost 12 and the non-enumerating login in `lms-access-api`;
  - the gateway's header check and rate limits (5 req/min for login, 60 req/min for the API);
  - the `VALIDATION_ERROR` envelope.
- Found and documented the real gaps:
  - no authorization by role or scope;
  - one internal secret shared by all services, and the callee does not check the caller;
  - the APIs are published on host ports `8081`-`8084`, so a client can skip the gateway;
  - no `iss`/`aud` check and no TLS;
  - no body size limit and no login field length limit.
- Built the threat model:
  - a data-flow diagram with 4 trust boundaries;
  - 10 STRIDE threats, each with its current control and a risk level;
  - a prioritized backlog SEC-01..SEC-10 with testable acceptance criteria;
  - a CI `security` job (gitleaks, govulncheck / npm audit, authorization tests) ready to add to the
    existing `ci.yml`.

## 3. Blockers and risks
- High risks open in the code today: T1 (internal token spoofing), T2 (no authorization), T3 (gateway
  bypass through host ports), T5 (no TLS once exposed outside localhost).
- The gateway config with rate limiting is on `feat/gateway-routing-and-rate-limiting`, not merged
  to `develop`, so a deploy from `develop` does not include it.
- The CI security gates only block merges once they are required status checks. Today no branch in
  the service repos has required checks (week 14).

## 4. Plan for next week
- Week 14: CI/CD end to end - make the PR gates (including this week's `security` job) blocking,
  protect `develop`/`qa`/`main`, and automate deployment of the tagged image.
- Start the backlog with SEC-01 (remove direct API ports) and SEC-02 (merge the gateway config),
  which are cheap and close the gateway-bypass threat.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria
- [ ] Tests added/updated (unit / integration)
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- Session documents in this same folder (`13-week/hu-status/`):
  - `SECURITY-CONTROLS.md` - Session 1: each security control, where it lives in the code, and its gaps
  - `THREAT-MODEL.md` - Session 2: data-flow diagram, trust boundaries, STRIDE, prioritized security
    backlog, CI security gates
- Code reviewed: `lms-api-gateway` (`deploy/nginx.conf` on `feat/gateway-routing-and-rate-limiting`),
  `lms-*-api/internal/adapter/in/httpapi/middleware/auth.go`,
  `lms-circulation-api/internal/adapter/out/catalogclient/client.go`,
  `lms-access-api/src/domain/access/administrator.js`, `lms-pod/compose.yml`
