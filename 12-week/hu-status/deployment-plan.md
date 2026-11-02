# Week 12 — Deployment Plan for MVP 3

Evidence attachment for `README.md` in this same folder (`12-week/hu-status/`). Covers both of
this week's sessions: Session 1 (deploy with a zero-downtime strategy + health checks) and
Session 2 (the written deployment plan), which are documented together here rather than split
across separate folders.

Project: **LMS Library**. Deploys the already-tagged `v2.0.1` cut (12 microservice repos —
`library-docs/12-ux-ui/mvp-release/v2.0.1.md`) to a production-like host, with a rolling,
zero-downtime strategy and real health checks gating the cutover.

**Target environment:** a single Docker host (`lms-pod/`, this team's local sibling-directory
assembly — `ADR-006-repo-per-context-decomposition.md`), standing in for "a production-like
host." No cloud account or managed orchestrator exists for a 3-person academic team
(`01-context/scope.md`, "no monetary budget"), and no container registry exists yet (`ADR-006`'s
Evaluated Alternatives explicitly discarded one as "infrastructure this project doesn't operate
yet") — so "the tagged image" is built locally from each repo's own `v2.0.1` git tag, not pulled
from a registry.

## 1. Artifact promotion

Reuses the project's own, already-accepted promotion model
(`library-docs/00-governance/branching-policy.md`) rather than inventing a parallel one:

| Stage | What "the artifact" is | How it gets there |
|---|---|---|
| `develop` | A Docker image built from `develop`'s HEAD, tagged `:dev-<short-sha>` | Automatic on every merge (once CI exists — `11-quality/testing-strategy.md` confirms none does yet; manual `docker build` today) |
| `qa` | The same image, re-tagged `:qa-<short-sha>`, after `develop`'s commits are re-applied to `qa` via `git cherry-pick -x` | Manual promotion PR (`qa/*` branch), per the existing model |
| `main` | A `vX.Y.Z` git tag (e.g. `v2.0.1`), and the image built from exactly that tag | Only after `qa`'s commits are re-applied to `main` via a `release/*` PR, approved by `ariel5253` |

No registry means the artifact that is "promoted" is the **git tag**, not a pushed image — each
environment builds its own image locally from the tag it is pinned to. Consistent with the
project's accepted "Minimal Infrastructure Footprint" principle (P5, `05-architecture/
overview.md`); closing this gap with a real registry + CI-built images is explicitly the CI/CD
milestone two weeks after this one (syllabus note below).

## 2. IaC for each environment

"Infrastructure as code" here is Docker Compose, not Terraform/Pulumi — proportionate to a
3-person team with no cloud account, the same scale justification `ADR-007` used for choosing
NGINX over a managed API gateway product:

| Environment | IaC file | What it assembles |
|---|---|---|
| Local / dev | `lms-pod/compose.yml` (sibling-directory assembly) | Every `-api`, `shell`, every portal, `lms-api-gateway` |
| Shared Postgres instance | `lms-infra-postgres/compose.yml` | One Postgres 16 instance, one schema per domain (Anexo J) |
| Shared MongoDB instance | `lms-infra-mongo/compose.yml` | Circulation's `loan_db` (`ADR-005`) |

A real "qa" and "main" environment (separate hosts/VMs) does not exist yet — a current gap this
plan's own acceptance criteria (§5) close, not papered over.

## 3. Zero-downtime strategy: rolling, not blue-green

| Strategy | What it needs | Why it was (not) chosen |
|---|---|---|
| **Rolling (chosen)** | One extra container at a time, health-checked before the old one stops | Fits a single Docker host with no load balancer beyond `lms-api-gateway` (built this cycle) and no spare capacity to run two full environments side by side |
| Blue-green | A second, full parallel environment (2x every service) plus a router that atomically flips 100% of traffic | Would double this team's already-minimal infra footprint for a 3-person project with no budget — the same reasoning that already rejected a container registry |

Rolling is possible here specifically *because* `lms-api-gateway` now exists: it is the one place
that decides whether a request reaches a service, so a service can be replaced one container at a
time behind it without the browser ever seeing a dropped connection.

**Health checks:** every `-api` service already exposes the two endpoints this deploy gates on
(verified in each repo's health handler, Go or — for `lms-access-api` post-`ADR-011` — Node):

- `GET /health` — liveness: the process is up (`200 {"status":"ok"}`)
- `GET /health/ready` — readiness: its database is reachable (`200`/`503`, 2s-bounded ping)

## 4. Deploy mechanism: `deploy-rolling.sh` (this folder)

Builds the requested git tag into a new image, starts it alongside the running container, polls
its own `/health/ready` until it answers `200`, **then** stops the old container — the old
container never goes down before the new one is proven ready. A rollback uses the identical
mechanism, pointed at the previous tag:

```bash
./deploy-rolling.sh circulation-service v2.0.1   # deploy
./deploy-rolling.sh circulation-service v2.0.0   # rollback — same script, previous tag
```

### Execution log and timed rollback — NOT YET RUN

**Honesty note:** the Docker daemon was not reachable in the sandbox this was written in
(`docker info` fails to connect — no Desktop/dockerd running there), so this script has **not**
actually been executed yet. What follows is the *expected* output given the script's logic and
every service's real, verified `/health`/`/health/ready` behavior — not a captured real run. This
must be re-run for real against `lms-pod` and this section replaced with actual output/timestamps
before being counted as real evidence:

```
$ ./deploy-rolling.sh circulation-service v2.0.1
[12-week] Building circulation-service from tag v2.0.1 (lms-circulation-api)...
[12-week] Starting circulation-service-next alongside circulation-service...
[12-week] Waiting for circulation-service-next /health/ready ... attempt 1/10: 503
[12-week] Waiting for circulation-service-next /health/ready ... attempt 2/10: 200 OK
[12-week] circulation-service-next is ready — cutting lms-api-gateway over
[12-week] Stopping old circulation-service container
[12-week] Renaming circulation-service-next -> circulation-service
[12-week] Done. Zero requests through lms-api-gateway were refused during the switch.
```

| Rollback step (illustrative) | Elapsed |
|---|---|
| Build `circulation-service:v2.0.0` from the previous tag | 0:00 – 0:41 |
| Start `-next`, first `/health/ready` poll (503 expected — DB pool warming) | 0:41 – 0:44 |
| Second poll — expected `200 OK` | 0:44 – 0:47 |
| Stop old container, rename `-next` → real name | 0:47 – 0:49 |
| **Target total rollback time** | **< 60 seconds** |

Real pass/fail criterion once actually run: zero refused requests to `lms-api-gateway`'s
`/health`, polled continuously from a separate terminal throughout the cutover.

## 5. Deployment stories, sliced with testable acceptance criteria

| ID | Story | Acceptance criteria |
|---|---|---|
| DEPLOY-001 | As an operator, I can deploy a specific git tag of any `-api` service with zero downtime | `./deploy-rolling.sh <service> <tag>` exits 0; a continuous poll of `lms-api-gateway`'s `/health` during the run shows 0 non-200 responses; the running container's `/health` reports the new tag's build |
| DEPLOY-002 | As an operator, a deploy that never becomes healthy aborts without touching the running service | Point `deploy-rolling.sh` at a tag with a deliberately broken `DB_HOST`; script exits 1 after `HEALTH_RETRIES` attempts; the original container is still running and still serving traffic throughout |
| DEPLOY-003 | As an operator, I can roll back to the previous tag in under 2 minutes | Timed rehearsal (§4) completes in < 120s end-to-end, with the same zero-dropped-request evidence as DEPLOY-001 |
| DEPLOY-004 | As a reviewer, every environment's assembly is defined in version-controlled IaC, not manual steps | `lms-pod/compose.yml`, `lms-infra-postgres/compose.yml`, `lms-infra-mongo/compose.yml` are the only inputs `docker compose up --build` needs; no undocumented manual step is required to reach a running system from a clean clone |
| DEPLOY-005 | As an operator, the artifact that gets promoted between environments is unambiguous | A given git tag (e.g. `v2.0.1`) produces a byte-identical image regardless of which environment builds it — verified by comparing `docker image inspect --format '{{.Id}}'` output built independently in two different directories from the same tag |

## Next

Per this week's own syllabus note: security, then CI/CD. CI/CD is what eventually replaces §1's
manual tag-promotion step and §2's "no registry" gap with real automation — flagged here as the
explicit next milestone, not solved in this plan.
