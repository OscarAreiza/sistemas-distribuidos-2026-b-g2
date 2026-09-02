# Week 1 — Session 2: Profile Repo, Docs Repo, and ADR-001

> **Added retroactively on 2026-09-02**, closing the same gap as
> `PROBLEM-AND-CONSISTENCY-MODEL.md` in this folder — see that file's header note for why.

---

## Profile repo and fork

Done that week, as required: the fork of this course repo exists at
`OscarAreiza/sistemas-distribuidos-2026-b-g2`, and `01-week/hu-status/README.md` was added with
the CONFIG block (`FULL_NAME`, `GITHUB_USER`, `TEAM`) matching the profile repo.

## The docs repo

The team's `docs` repo is `code-corhuila/library-docs` — created and structured around the
15 numbered sections (`00-governance` through `15-project-control`) the course template
defines. It is the declared source of truth for the whole project: "if code and docs disagree,
the docs repo wins until this README (or a PR) updates it" (`lms-library`'s own root
`README.md`).

## ADR-001 — a numbering note, told straight

Session 2 asked specifically for **"ADR-001 (chosen architecture style)."** In this team's
actual `library-docs` repo, that is not what ADR-001 is:

| ADR | Title | What it actually decides |
|---|---|---|
| `ADR-001` | Documentation Language | English for all code/docs, over Spanish or a split |
| `ADR-002` | **Architectural Style: Hexagonal Modular Monolith** | The chosen architecture style — this is the ADR Session 2 was really asking for |
| `ADR-003` | NGINX as Reverse Proxy / Edge Layer | The edge/gateway layer in front of the architecture ADR-002 describes |
| `ADR-004` | Incremental Microservices Decomposition | Records that Access and Membership have since been extracted into their own services, one bounded context at a time, superseding ADR-002's "ship as one deployable" framing while keeping its internal hexagonal structure |

**Why it landed this way:** the team wrote ADR-001 for the language decision first because it
was the most immediate, lowest-stakes decision to formalize before anything else got written —
not because the numbering was planned around Session 2's specific ask. The architecture-style
decision came right after, as ADR-002, following the same template
(`05-architecture/decisions/_template-adr.md`) and the same rigor (evaluated alternatives,
consequences, risks) the assignment asks for — it is simply the second ADR chronologically, not
the first.

**The chosen architecture style itself** (`ADR-002`, in full): a **Hexagonal Modular
Monolith** — one deployable Go service (`library-api`), internally organized into one module
per bounded context (Access, Circulation, Catalog, Membership), each built with Ports &
Adapters, chosen over (a) full microservices from day one — rejected for requiring 4
deployables, a broker, and 4 databases a 3-person academic team cannot operate — and (b) a
plain layered monolith without ports/adapters — rejected for coupling business logic to the
framework/DB and making a later split much harder. See `ADR-002` for the full evaluated-
alternatives table and consequences.

**Status as of this writing:** that decision has since evolved — Access and Membership are
already extracted into their own services (`ADR-004`), ahead of the "v2, not scheduled"
timeline `ADR-002` originally assumed, because the team's own module boundaries from `ADR-002`
made the split low-risk sooner than expected. The chosen architecture *style* (hexagonal,
ports & adapters per bounded context) never changed — only the deployment topology it produced.

## Backlog with testable acceptance criteria

Drafted and refined since: `library-docs/03-product/product-backlog.md` (MoSCoW-prioritized)
and `library-docs/04-requirements/user-stories.md` (Gherkin acceptance criteria per story).

---

## Correlations

- The architecture ADR this document is about → `library-docs/05-architecture/decisions/records/ADR-002-hexagonal-modular-monolith.md`
- Its current, evolved state → `library-docs/05-architecture/decisions/records/ADR-004-incremental-microservices-decomposition.md`
- Full ADR register → `library-docs/05-architecture/decisions/README.md`
- Companion document (Session 1: problem, backlog, consistency/delivery semantics) → `PROBLEM-AND-CONSISTENCY-MODEL.md` in this same folder
