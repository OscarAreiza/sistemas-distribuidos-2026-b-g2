# Week 13 — Session 1: Securing the system

What the LMS system does today for each security control, checked in the code of each repo
(`develop` branch unless noted), and what is still missing. The missing items become the security
backlog in `THREAT-MODEL.md`.

## Summary

| Control | Where | Status |
|---|---|---|
| Authenticate at the edge | `lms-api-gateway` (nginx) | Partial. Written on a feature branch, not merged to `develop`. |
| Validate tokens in each service | `RequireAuth` middleware in each Go API | Done |
| Authorize in each service | Go APIs | Gap. Any valid token gets full access; there are no roles or scopes. |
| Secure service-to-service calls | Internal HS256 token minted by each client | Done, with known limits |
| Hash secrets | `lms-access-api` (bcrypt) | Done |
| Validate inputs | Handlers of each API | Partial |

## 1. Authenticate at the edge

`lms-api-gateway`, branch `feat/gateway-routing-and-rate-limiting`, file `deploy/nginx.conf`
(ADR-007):

- **One entry point.** `/api/v1/*` is routed by path prefix to the owning service, and the SPA is
  served from the same origin.
- **Fast 401.** Protected routes (`/students`, `/books`, `/loans`) without an `Authorization` header
  get `401` at the gateway, before they reach a backend. The gateway does **not** verify the JWT;
  that stays in each service (ADR-003, ADR-007).
- **Rate limiting:** login is limited to 5 requests/min per IP (burst 2), and the rest of the API to
  60 requests/min per IP (burst 20). Over the limit the gateway answers `429` with `Retry-After`.

**Gaps:**
- The branch is not merged into `develop`, so a deploy from `develop` does not include this
  config.
- `lms-pod/compose.yml` also publishes each API directly on the host (`8081`–`8084`), so a client
  can skip the gateway and its rate limits.

## 2. Validate tokens in each service

Each Go API (`lms-membership-api`, `lms-catalog-api`, `lms-circulation-api`) wraps its protected
routes with `internal/adapter/in/httpapi/middleware/auth.go` → `RequireAuth`:

- It requires `Authorization: Bearer <jwt>`, and answers `401 UNAUTHORIZED` otherwise.
- It accepts **only RS256 and HS256** (`jwt.WithValidMethods`). It picks the key **by the token's
  algorithm**: the RS256 public key of `lms-access-api` for Administrator sessions, and a separate
  internal secret for service calls. This blocks the algorithm-confusion attack, where an RS256
  public key is reused as an HMAC secret.
- `exp` is enforced by the library: an expired token is rejected.
- The `sub` claim (Administrator ID) is put in the request context for the handlers and logs.
- The service refuses to start without `JWT_PUBLIC_KEY` and `INTERNAL_JWT_SECRET` (`config.go`).

Tokens are issued by `lms-access-api` (`infrastructure/auth/jwtIssuer.js`): RS256, signed with a
private key that only access-service holds, and expiring after `JWT_EXPIRY` (default 1 h).

**Gaps:**
- `iss` and `aud` are not checked, so a token is not tied to the LMS issuer or audience.
- There is no revocation or logout. A stolen token is valid until it expires.

## 3. Authorize in each service

**Gap.** There is one role (Administrator) and no permission check. Any valid token, including an
internal service token, can call every protected route of every service. For example, the internal
token that circulation uses to call catalog could also call `DELETE` routes in membership.

## 4. Secure service-to-service calls

`lms-circulation-api` → `lms-catalog-api` / `lms-membership-api`
(`internal/adapter/out/catalogclient/client.go`, `mintServiceToken`):

- Each outbound call carries a fresh **HS256 token** signed with `INTERNAL_JWT_SECRET`, with
  `sub` set to the calling service (`circulation-service`), and expiring after **1 minute**.
- The internal secret is separate from the RS256 key used for user sessions, so a user token and a
  service token cannot be swapped.
- Calls go over the Docker internal network, by service name.

**Gaps:**
- **One secret is shared by all services.** Any service, or anyone who reads the secret, can mint
  a token that passes as any service, or as any `sub`.
- The callee does not check **which** service is calling. `sub` is not compared against an
  allow-list per route.
- Traffic inside the Docker network is plain HTTP, without TLS.

## 5. Hash secrets

`lms-access-api/src/domain/access/administrator.js`:

- Passwords are hashed with **bcrypt, cost 12** (`BCRYPT_COST = 12`), and only `password_hash` is
  stored (`access.administrator`). Invariant INV-001: a password is never stored in plain text.
- Login uses `bcrypt.compareSync`. A wrong user and a wrong password return the same
  `401 INVALID_CREDENTIALS` with the same message, so usernames cannot be enumerated.
- Keys and secrets (`JWT_PUBLIC_KEY`, `INTERNAL_JWT_SECRET`, DB passwords) come from environment
  variables. Per-environment secret files are git-ignored (commit `a628fba` in this repo).

## 6. Validate inputs

- **Go APIs:** a body that fails JSON decoding returns `400 VALIDATION_ERROR`. Field errors use one
  `ValidationError` envelope with a `FieldDetail` per invalid field (`response.go`). IDs are parsed
  before they reach the domain.
- **`lms-access-api` login:** `username` and `password` must be strings, otherwise `400
  VALIDATION_ERROR`.

**Gaps:**
- No **body size limit** (`http.MaxBytesReader`), so a very large body is read fully into memory.
- Unknown JSON fields are accepted silently (no `DisallowUnknownFields`).
- Login does not limit the length of `username` and `password`. bcrypt only uses the first 72
  bytes, and a long string costs CPU on every attempt.

## What is next

Each gap above is a threat in `THREAT-MODEL.md` and a story in its prioritized security backlog.
