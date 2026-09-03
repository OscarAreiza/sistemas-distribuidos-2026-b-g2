# Week 5 - Session 2: Ship MVP 1

> Goal: promote to `main`, tag `v1.0.0`, verify the checklist/DoD with evidence, demo the
> running system, and run the retrospective.
>
> **Status: in progress, not shipped yet.** The `v1.0.0` tag and the `main` promotion are being
> done manually by the team. **The scope question below is now resolved** (see `README.md` in
> this same folder) — Circulation ships after the tag, not before — so what remains is executing
> steps 3–5 below, not deciding them.

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
— the loan/return cycle is the actual point of a library system. Tagging `v1.0.0` on the resolved
scope ships login + student management + book catalog, with no way to actually lend or return a
book yet — an accepted, deliberate trade-off (see Section 3), not an oversight.

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

## 3. Before the `v1.0.0` tag

1. **Resolved:** MVP 1 ships with login + students + catalog (HU-01/02/03/04/05). Circulation
   (HU-06/07/08) is **not** part of this cut — it's a fast-follow implemented and merged after
   the `v1.0.0` tag, per the decision recorded in this folder's `README.md` (Section 3, "Blockers
   and risks") and formalized architecturally in `library-docs/05-architecture/decisions/records/ADR-005-mongodb-for-circulation-service.md`
   (Circulation's own database engine — MongoDB, decided as part of finally scoping its build).
2. Update `library-docs/03-product/product-backlog.md` and
   `04-requirements/traceability-matrix.md` to reflect this actual shipped scope (HU-01–05), not
   the originally planned Cut 1 (which included HU-06/07).
3. Promote `dev -> QA -> main` per `00-governance/git-conventions.md`'s branch strategy, once
   HU-03/HU-05 (student/book editing) are merged alongside HU-01/02/04.
4. Tag `main` as `v1.0.0`.
5. Record the demo and hold the retrospective; log Circulation as the carried-over scope in
   `library-docs/15-project-control/technical-backlog.md`, not as unplanned technical debt.

## 4. Retrospective input (to discuss live, not pre-decided here)

- **What went well:** the hexagonal module boundaries from `ADR-002` made extracting Access and
  Membership into real services a small, low-risk change rather than a rewrite — exactly what
  that ADR predicted.
- **What to improve:** documentation (`library-docs`) drifted from the code for a while during
  the microservices extraction — worth agreeing on updating both in the same PR going forward,
  per `ADR-004`'s own risk table.
- **What's blocking MVP 1:** nothing anymore — Circulation (HU-06/07/08) was deliberately moved
  out of the MVP 1 cut (Section 3) rather than blocking the tag; still worth discussing live
  whether the Core Domain being the last thing built was the right call for future cuts.

## Correlations

- Product backlog / MoSCoW prioritization → `library-docs/03-product/product-backlog.md`
- Definition of Done → `library-docs/00-governance/definition-of-done.md`
- Current per-context implementation state → `library-docs/05-architecture/decisions/records/ADR-004-incremental-microservices-decomposition.md`
- Traceability (what's really merged vs. built-but-not-merged) → `library-docs/04-requirements/traceability-matrix.md`
