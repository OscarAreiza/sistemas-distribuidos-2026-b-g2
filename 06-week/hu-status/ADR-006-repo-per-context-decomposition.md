# ADR-006 — Repo-per-Context Decomposition (Monorepo → 18 Repositories)

- **ID:** ADR-006
- **Date:** 2026-09-15
- **Status:** Accepted
- **Authors:** Oscar Areiza — Tech Lead
- **Reviewers:** Hermes Pascuas, Luis Alejandro Meneses — Development team

---

## Context

`ADR-004` decomposed `lms-library` incrementally, one bounded context at a time, but every
extracted service (`access-service`, `membership-service`, and now `catalog-service`) still lives
as a subfolder inside the **same** repository, sharing one `docker-compose.yml`, one frontend
SPA, and one branch/PR flow. That was a deliberate first step — "the context stays inside the
`library-api` monolith until its own extraction lands," per `ADR-004` — not the end state.

The course organization (`code-corhuila`) created 18 empty repositories on 2026-09-14, each with
a README that states the project's team (`lms-library`, Grupo 2) and a `## Migration scope`
section describing exactly what moves there from `lms-library`, or "from scratch" where nothing
exists yet. This is not a proposal the team is choosing between alternatives on — the repos, their
names, and their scopes already exist as a **course requirement**, in the same way
`00-governance/branching-policy.md` is "not negotiable." This ADR records the team's
acknowledgment of that requirement and the concrete plan to execute it, since `library-docs` is
where every one of those 18 READMEs says "the full map lives."

**Known constraints:**
- Same 3-person academic team as every prior ADR — an 18-repository migration cannot happen in
  one PR; it is executed context by context, the same incremental discipline `ADR-004` already
  established, not a big-bang cutover
- The new repos are currently empty (no commits) — this ADR is written before any code has moved,
  the same way `ADR-005` was written before `circulation-service` existed
- Two of the 18 repos (`lms-circulation-api`, `lms-circulation-db`) and two cross-cutting ones
  (`lms-workflow`, `lms-worker`) have no equivalent in `lms-library` at all — they are new scope,
  not a lift-and-shift
- No Kubernetes, no service mesh, no managed service-discovery product — same "Minimal
  Infrastructure Footprint" principle (`05-architecture/overview.md`, P5) applies across 18 repos
  just as it did across 1

---

## Decision

**We decided:** decompose `lms-library` into 18 single-purpose repositories — one API repo, one
database repo, and one UI "portal" repo per bounded context (Access, Membership, Catalog,
Circulation), plus four cross-cutting repos: `lms-api-gateway` (single entry point — its exact
auth/rate-limiting responsibilities are `ADR-007`'s decision, not this one), `lms-front` (the
micro-frontend shell that packages the four portals as remotes), `lms-workflow` (the loan saga —
scope defined in `ADR-008`), `lms-worker` (asynchronous jobs — scheduling model defined in
`ADR-009`), and `lms-infra` (assembles the per-service repos into one deployable system: compose,
environments, observability, secrets). `lms-library` itself stops receiving new bounded-context
code once its pieces are migrated; it is not deleted, since `main` there still holds the released
`v1.0.0` history.

**How the repos find each other:** no service mesh, no managed discovery — `lms-infra` keeps
composing every service onto one shared Docker Compose network (`lms-network`), exactly as
`lms-library`'s `docker-compose.yml` does today, so container-name-based calls already in the
code (e.g. `membership-service`'s client calling `http://circulation-service:8080`) keep working
unmodified. To let `lms-infra` build images from 17 sibling repos without git submodules (real
operational friction for a 3-person team) or a container registry/CI pipeline (neither exists yet
— `11-quality/testing-strategy.md` and the DoD checklist both note there is no CI today), the
convention is: **clone every repo as a sibling folder** (`../lms-access-api`, `../lms-access-db`,
...) and reference them from `lms-infra`'s compose file by relative build context path. This is
documented in `lms-infra`'s own `SETUP.md` once that repo receives its first commit.

**Justification:** this is `ADR-004`'s own "Database per Service" and "own deployable" goals
taken to their conclusion — a service's database now has its own repository lifecycle
(migrations, seeds, backups) entirely independent of its API's, and a bounded context's UI is no
longer a folder inside one shared SPA but an independently deployable "portal," composed at
runtime by `lms-front`. `ADR-002`'s hexagonal internal structure is exactly what makes each
`api` repo's migration a folder move (`cmd/`, `internal/{domain,application,infrastructure}`),
not a rewrite — the same property that made `ADR-004`'s incremental split low-risk applies again
here, one repo at a time.

**Migration order:** Access first — every other context depends on token validation, so it has no
dependency on anything else being migrated first (the same reason `ADR-004` extracted it first).
Then Membership and Catalog, in either order — both are already-working code being relocated, not
built. Then `lms-infra`, once at least two services exist to actually assemble. Then
`lms-api-gateway` and `lms-front`, once the repos they route to/compose exist. Circulation
(`lms-circulation-api`/`-db`/`-portal`), `lms-workflow`, and `lms-worker` come last — this is new
work, not a relocation, and `lms-workflow`/`lms-worker` specifically depend on the other three
contexts already being reachable as independent services. Full detail per repo →
`09-microservices/repo-migration-map.md`.

**Cut-over criterion:** a context's new repos (`lms-<ctx>-api`/`-db`/`-portal`) become
authoritative **only once `lms-infra`'s compose file routes to them instead of to
`lms-library`'s copy** — not merely once code has been copied over. Until that switch, whatever
is actually running (`lms-library`'s copy) remains the source of truth, even if a newer copy
already exists in the destination repo mid-migration. This avoids the two-copies-no-stated-owner
problem: authority follows what's deployed, not what's been copied. Once cut over,
`lms-library`'s corresponding folder is frozen (no further commits) — `library-docs` and
`09-microservices/repo-migration-map.md`'s status column are updated in the same PR that flips
`lms-infra`'s routing, not as a follow-up.

**Durable state for the cross-cutting repos:** the topology only assigns a `-db` repo to the four
domain contexts, not to `lms-workflow` or `lms-worker` — this was undecided until now. Decision:
**neither needs its own database.** Both are designed to be stateless, re-deriving what they need
by querying `lms-circulation-api`/`lms-membership-api`/`lms-catalog-api`'s own public state on
every run rather than keeping a separate execution ledger:
- `lms-workflow`'s saga (`ADR-008`) checks each participant's current state directly (is this
  loan already marked late? is the student already suspended?) and acts on the delta — both steps
  are individually idempotent and re-checkable, so there is nothing a separate saga-state store
  would add that querying the participants doesn't already give
- `lms-worker` (`ADR-009`) is the same: it finds work by querying for records in a given state
  (e.g. active loans past `dueDate`) rather than remembering what it processed last run

If a future requirement genuinely needs saga-execution history that outlives a single run (e.g.
an audit trail of every compensation attempt), that is a new decision to make explicitly when it
comes up — not something to silently bolt on to `lms-workflow`/`lms-worker` later.

---

## Evaluated alternatives

| Alternative | Pros | Cons | Reason for discarding |
|------------|------|------|-----------------------|
| **18 repos, repo-per-artifact (CHOSEN)** | Independent CI/CD, ownership, and history per artifact; database lifecycle fully decoupled from API lifecycle; enables true micro-frontend composition via `lms-front` | 18 repositories for a 3-person team to keep in sync (branch policy, CODEOWNERS, README) instead of 1; cross-repo changes (e.g. an API contract change touching its portal) now span 2+ PRs instead of 1 | Not actually a choice — the repos already exist under the course org with this exact shape; the only real decision left is *how* to execute the migration, not *whether* |
| Stay in `lms-library`, one repo per context as a folder (`ADR-004`'s current state) | Zero new repos to manage; single PR flow already working | Doesn't match the repos the course actually created and requires (`00-governance/branching-policy.md` describes the same "Course rule — not negotiable" pattern for `<abbr>-docs`) | Overtaken by the course's own infrastructure — the 18 repos exist regardless of this ADR |
| Git submodules to assemble `lms-infra` from the other 17 repos | Pinned, versioned references between repos; standard Git mechanism | Real operational friction for a 3-person academic team (detached-HEAD confusion, easy to forget `submodule update`); no CI exists yet to automate the pin bump | Sibling-directory clone + relative compose paths gives the same "assemble everything locally" outcome with none of submodules' footguns, at this team's size |
| Container registry (GHCR) + `lms-infra` pulls prebuilt images | Decouples "build" from "run"; closer to a real production setup | Requires a CI pipeline per repo to build and publish images — none exists today (a real, separately-tracked gap, not something to build as a side effect of this ADR) | Adds infrastructure this project doesn't operate yet; contradicts P5 ("Minimal Infrastructure Footprint") until CI is a deliberate, separate decision |

---

## Consequences

**Positive:**
- Each bounded context's database (`lms-<ctx>-db`) can evolve, be backed up, and be restored
  independently of its API's release cycle — the natural conclusion of `AT-001`'s "Database per
  Service" debt-payoff, applied to every context at once instead of one at a time
- `lms-front` + four portals is a real micro-frontend architecture, not a single SPA with
  route-based code-splitting — matches how the bounded contexts are already modeled in
  `02-domain/domain-map.md`
- Existing service-to-service code (container-name HTTP calls) needs **zero changes** thanks to
  the sibling-clone + shared-network convention — the migration is a folder move, not a rewrite,
  at the networking level too

**Negative / Trade-offs:**
- 18 repositories to keep individually compliant with `00-governance/branching-policy.md`
  (branch protection, CODEOWNERS, PR flow) instead of 1 — real ongoing overhead for a 3-person team
- `lms-front`'s README already flags the sharpest risk: if each portal reimplements its own HTTP
  client/session handling instead of consuming `lms-front`'s shared `lib/api.ts`/`lib/auth.ts`,
  the migration duplicates the hardest part of the frontend — must be enforced by review, not
  tooling
- Cross-repo consistency (an API's OpenAPI contract change vs. its portal's client code) now
  requires coordinating PRs across repositories, with no shared CI to catch drift automatically
- The sibling-clone convention is undocumented tooling knowledge until `lms-infra`'s `SETUP.md`
  exists — anyone cloning the system for the first time needs that file, not tribal knowledge

**Impact on the system:**
- Affected repos: all 18 new repos (currently empty) plus `lms-library` (source of the
  migration, stops receiving new bounded-context code once each piece moves out) and
  `library-docs` (this ADR, the migration map, and every cross-reference below)
- Documents that must be updated as each repo receives its first real commit:
  `09-microservices/repo-migration-map.md` (status column), `09-microservices/service-catalog.md`
  (folder/repo links), `05-architecture/overview.md` (topology diagrams), `05-architecture/deployment.md`
  (infra now assembled from `lms-infra` instead of one `docker-compose.yml` in `lms-library`)

---

## Risks

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|-----------|
| A repo's README migration-scope note is wrong (already found once: `lms-membership-portal` claims "no `pages/students`" in `lms-library`, but that folder exists with a working HU-03 implementation) | Medium | Medium — someone rebuilds already-working code from scratch | Cross-check every "from scratch" claim against the real `lms-library` tree before starting that repo's migration; report the specific error to whoever owns the course repos |
| 18 repos drift out of sync with `00-governance/branching-policy.md` (missing CODEOWNERS, wrong branch names) as they're created one at a time | Medium | Medium — inconsistent review/gating across the system | Copy the branching setup from the first fully-migrated repo (`lms-access-api`, since its README is already confirmed accurate) as the template for the rest, rather than reinventing it per repo |
| Sibling-directory convention breaks if someone clones repos into a different layout | Low | Low — `lms-infra`'s compose fails to find a build context | Document the exact expected directory layout in `lms-infra/SETUP.md` the moment that repo exists, not as an afterthought |
| Someone treats a context's new repo as authoritative once code is copied there, before `lms-infra` actually routes to it — the exact two-copies-no-owner scenario the cut-over criterion exists to prevent | Medium | Medium — bug fixes land in the copy nobody is running | Enforce the cut-over criterion above literally: a PR that copies code to `lms-<ctx>-api` does not by itself freeze `lms-library`'s folder; only the `lms-infra` routing-switch PR does |

---

## References

- Extends `ADR-004-incremental-microservices-decomposition.md` — this ADR is that decomposition
  carried from "folder inside one repo" to "repo per artifact"
- `lms-api-gateway`'s auth/rate-limiting scope → `ADR-007-gateway-auth-and-rate-limiting.md`
- `lms-workflow`'s saga scope → `ADR-008-circulation-saga-scope.md`
- `lms-worker`'s scheduling model → `ADR-009-worker-scheduling-model.md`
- Full repo-by-repo migration plan → `09-microservices/repo-migration-map.md`
- Hexagonal structure that makes each `api` repo's migration a lift, not a rewrite →
  `ADR-002-hexagonal-modular-monolith.md`
- Course-mandated branch/PR policy every new repo must follow →
  `00-governance/branching-policy.md`
