# Session 2 — MVP 2 Release

> Brief: "Ship MVP 2: promote to main, tag v2.0.0, verify the checklist (including a
> failure/compensation path), demo the integrated system with an injected failure, and
> retrospect. Ensure each member's hu-status evidence is complete. On to Corte 3."

Written day 1 of this week, as a plan — not a record of a release that happened. Every item below
is `todo`, in the order it has to happen, because each one is a real prerequisite for the next.

---

## 1. Why nothing here can start yet

Per `library-docs/00-governance/branching-policy.md`: a release is cut from `main` and filled
**gradually**, one commit per user story **already validated in `qa`**. `main` is downstream of
`qa`, `qa` is downstream of `develop` — skipping a stage isn't a shortcut, it's a different
(non-compliant) process. Two things currently block even starting:

1. **`qa` itself isn't settled.** Every domain's `qa` promotion PR is still open, unreviewed, as
   of this write-up:

   | Repo | PR | Commits |
   |---|---|---|
   | `lms-circulation-api` | `#6` | 78 |
   | `lms-circulation-db` | `#6` | 23 |
   | `lms-circulation-portal` | `#7` | 16 |
   | `lms-catalog-api` | `#8` | 59 |
   | `lms-catalog-db` | `#5` | 50 |
   | `lms-catalog-portal` | `#4` | 26 |
   | `lms-membership-api` | `#13` | 11 |
   | `lms-membership-db` | `#6` | 4 |
   | `lms-membership-portal` | `#7` | 6 |

   Promoting to `main` while `qa` is still mid-review would mean releasing code nobody but its
   own author has looked at — exactly what the three-tier branch model exists to prevent.

2. **The failure/compensation path doesn't exist yet to demo.** The release checklist explicitly
   requires verifying it — see `PERSISTENCE-SAGA-OUTBOX.md`, this same folder: `lms-workflow` is
   still zero code. A release that claims to demo a failure/compensation path that isn't built
   would be a fabricated checklist item, not a shortcut.

**Sequencing, not a side effect:** Session 1's work (this same week) is the actual critical path
to Session 2, not parallel work. Nothing in this document can be marked `done` honestly before
`PERSISTENCE-SAGA-OUTBOX.md`'s four items are.

## 2. The release checklist (todo, in order)

- [ ] **Get every open `qa` promotion PR reviewed and merged** (table above) — `qa` has to
  actually reflect what's being released, with real review, not just CI green.
- [ ] **Scaffold and implement `lms-workflow`**'s saga (per `PERSISTENCE-SAGA-OUTBOX.md` §2-4) —
  the one new capability this release adds over MVP 1, and the one the failure-path demo depends
  on.
- [ ] **Promote `lms-workflow`'s saga to `qa` the same way** (`cherry-pick -x`, per
  `library-docs/00-governance/git-conventions.md` — never a direct merge) and validate it there,
  per `library-docs/11-quality/qa-promotion-standard.md`'s Stage 1-4.
- [ ] **Cut `release/2.0.0`** from `main`, filled gradually — one commit per story already
  validated in `qa`, each with its `cherry picked from` trail, per `branching-policy.md`'s
  release-branch model.
- [ ] **Verify the Definition of Done** (`library-docs/00-governance/definition-of-done.md`)
  against the actual release contents — not a rubber stamp; the same honest standard this status
  repo has used all along (see week 9's `qa-promotion-standard.md` for what "verify" means in
  practice here).
- [ ] **Verify the failure/compensation path specifically**, per `ADR-008`'s own compensation
  rule: force the membership-suspension call to fail after the loan is marked late, confirm the
  saga retries on its next scheduled pass (not that it rolls back the "marked late" step — per
  `PERSISTENCE-SAGA-OUTBOX.md` §2, that step is correct to leave in place), and confirm a
  re-triggered pass doesn't double-suspend (the `loan_id` idempotency check, §4 of that same
  document).
- [ ] **Merge `release/2.0.0` into `main`** via PR, 1 approval from `ariel5253` (course rule, no
  exceptions, per `branching-policy.md`).
- [ ] **Tag `v2.0.0`** on `main`, right after that PR merges — not a separate later step
  (`branching-policy.md`'s own tagging rule).
- [ ] **Demo the integrated system with an injected failure** — the same scenario just verified
  above, run live for the class/instructor, not just asserted in a document.
- [ ] **Retrospective** — what worked, what didn't, carried into Corte 3's own planning (not
  written yet; this document isn't the retrospective, it's the pre-release plan).
- [ ] **Confirm every team member's `hu-status` evidence is complete** — Hermes's and Luis's own
  forks, not assumed complete because this one is.

## 3. What "complete" does NOT mean here

Carrying forward the same discipline `08-week`/`09-week` already used: a checked box needs real
evidence attached (a PR link, a command's real output, a screenshot of the live demo) — not "done
because it's on the list." None of the boxes above are checked yet, on purpose: checking one
before its evidence exists would be exactly the kind of gap this status repo's whole approach has
been built to avoid.

## Correlations

- What each unchecked box depends on technically → `PERSISTENCE-SAGA-OUTBOX.md`, this same folder
- Release branch mechanics → `library-docs/00-governance/branching-policy.md`
- Definition of Done → `library-docs/00-governance/definition-of-done.md`
- QA-stage verification standard → `library-docs/11-quality/qa-promotion-standard.md`
