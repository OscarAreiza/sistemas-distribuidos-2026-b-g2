# Week 5 - Session 2: Ship MVP 1

> Goal: promote to `main`, tag `v1.0.0`, verify the checklist/DoD with evidence, demo the
> running system, and run the retrospective.
>
> **Status: in progress, not shipped yet.** The `v1.0.0` tag and the `main` promotion are being
> done manually by the team once the scope question below is resolved — this document is the
> honest checklist of where things stand, not a claim that MVP 1 has already shipped.

---

## 1. What "MVP 1" was supposed to include

Per `library-docs/03-product/product-backlog.md`, Cut 1 (Must Have) is:

| HU | Title | Merged to `dev`? |
|---|---|---|
| HU-01 | Administrator Authentication | ✅ Yes |
| HU-02 | Student Registration | ✅ Yes |
| HU-04 | Book Registration | ✅ Yes |
| HU-06 | Loan Registration | ❌ No — not merged, and Circulation has no implementation at all yet |
| HU-07 | Return Registration & History Tracking | ❌ No — same as above |

**Real gap:** Circulation is this system's Core Domain (`library-docs/02-domain/domain-map.md`)
— the loan/return cycle is the actual point of a library system. Tagging `v1.0.0` today ships
login + student management + book catalog, with no way to actually lend or return a book.

## 2. Definition of Done — verification against the real repo state

Per `library-docs/00-governance/definition-of-done.md`, checked against what is actually true
today (not aspirational):

| DoD item | Status | Note |
|---|---|---|
| Code implements all acceptance criteria of the story | ✅ (for HU-01, 02, 04) | HU-06/07/08/09/03/05 exist on unmerged branches |
| Code reviewed via PR | ✅ | PRs [#8](https://github.com/OscarAreiza/lms-library/pull/8), [#9](https://github.com/OscarAreiza/lms-library/pull/9) went through review before merging |
| Linting/formatting pass in CI | ❌ Not applicable | No CI pipeline exists yet — `go vet`/`go build` are run manually before opening a PR |
| Unit tests written | ✅ | Every merged use case has a `_test.go` with fakes (`11-quality/tdd-guide.md` pattern) |
| Tests pass locally and in CI | 🟡 Partial | Pass locally (verified inside a throwaway Docker build); no CI to run them automatically |
| Integration tests pass | 🟡 Partial | Verified manually via `docker compose up` + real HTTP calls between services; no automated integration suite |
| OpenAPI contract updated | ❌ Not done | `07-api/contracts/openapi/` still has the generic template, not real per-service contracts — tracked as a known gap |
| Deployed to staging, smoke test passing | ❌ Not applicable | No staging environment exists (`05-architecture/deployment.md`) — local Docker Compose is the only environment today |
| Service README updated | ✅ | Every service has an accurate `README.md` and a `library-docs/09-microservices/services/NN-*/README.md` entry |
| ADR created for significant decisions | ✅ | `ADR-004-incremental-microservices-decomposition.md` |

## 3. Before the `v1.0.0` tag (team decision needed)

1. **Decide the real MVP 1 scope**: ship now with login + students + catalog only (dropping
   HU-06/07 from this cut), or hold the tag until Circulation is implemented and merged.
2. Whichever is chosen, update `library-docs/03-product/product-backlog.md` and
   `04-requirements/traceability-matrix.md` to reflect the actual shipped scope — not the
   originally planned one, if it changes.
3. Promote `dev -> QA -> main` per `00-governance/git-conventions.md`'s branch strategy.
4. Tag `main` as `v1.0.0`.
5. Record the demo and hold the retrospective; log any carried-over technical debt in
   `library-docs/15-project-control/technical-backlog.md`.

## 4. Retrospective input (to discuss live, not pre-decided here)

- **What went well:** the hexagonal module boundaries from `ADR-002` made extracting Access and
  Membership into real services a small, low-risk change rather than a rewrite — exactly what
  that ADR predicted.
- **What to improve:** documentation (`library-docs`) drifted from the code for a while during
  the microservices extraction — worth agreeing on updating both in the same PR going forward,
  per `ADR-004`'s own risk table.
- **What's blocking MVP 1:** Circulation (HU-06/07/08) — the team should discuss whether it was
  under-prioritized relative to its status as the Core Domain.

## Correlations

- Product backlog / MoSCoW prioritization → `library-docs/03-product/product-backlog.md`
- Definition of Done → `library-docs/00-governance/definition-of-done.md`
- Current per-context implementation state → `library-docs/05-architecture/decisions/records/ADR-004-incremental-microservices-decomposition.md`
- Traceability (what's really merged vs. built-but-not-merged) → `library-docs/04-requirements/traceability-matrix.md`
