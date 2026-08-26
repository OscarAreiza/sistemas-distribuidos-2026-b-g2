<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       03-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 03

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Oscar Mauricio Areiza Paramo
- GITHUB_USER: OscarAreiza
- TEAM: lms-library
- SPRINT_GOAL: Start the library-docs SDD documentation (01-context, 02-domain) and define the initial user story backlog with acceptance criteria
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-01 | Administrator Authentication | todo | `03-week/hu-status/HU-BACKLOG.md` |
| HU-02 | Student Registration | todo | `03-week/hu-status/HU-BACKLOG.md` |
| HU-03 | Student Search, Editing & Deactivation | todo | `03-week/hu-status/HU-BACKLOG.md` |
| HU-04 | Book Registration | todo | `03-week/hu-status/HU-BACKLOG.md` |
| HU-05 | Inventory Control & Book Search | todo | `03-week/hu-status/HU-BACKLOG.md` |
| HU-06 | Loan Registration | todo | `03-week/hu-status/HU-BACKLOG.md` |
| HU-07 | Return Registration & History Tracking | todo | `03-week/hu-status/HU-BACKLOG.md` |
| HU-08 | Late-Return Penalty System | todo | `03-week/hu-status/HU-BACKLOG.md` |

## 2. My individual contribution
- Updated the LMS tech stack decision: frontend changed from Vue.js to React.
- Filled in `library-docs/01-context/` (system overview, scope, glossary) and `library-docs/02-domain/` (bounded contexts, entities and business rules, domain event catalog), committed and pushed directly to `main` per the repo's documentation-only Git policy.
- Defined the initial Product Backlog: 8 user stories across 4 epics, each with acceptance criteria, in `03-week/hu-status/HU-BACKLOG.md`.

## 3. Blockers and risks
- None blocking — remaining open items (renewals, reservations, monetary fines) were deliberately deferred out of v1 scope, see `library-docs/01-context/scope.md`.

## 4. Plan for next week
- Move the HU-BACKLOG.md user stories into `library-docs/04-requirements/user-stories.md` and continue the SDD weekly order (`03-product`, `04-requirements`, non-functional requirements).

## 5. Compliance self-check
- [ ] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [ ] Testable acceptance criteria
- [ ] Tests added/updated (unit / integration)
- [ ] DDD / hexagonal boundaries respected (domain has no I/O)
- [ ] No secrets; config via environment variables

> Conventional Commits were used for the `library-docs` commits below. Branch/PR flow, tests,
> and hexagonal boundaries are not applicable yet — no code exists, only documentation and
> domain modeling.

## 6. Evidence links
- `03-week/hu-status/HU-BACKLOG.md`
- https://github.com/code-corhuila/library-docs/commit/b3caefaf8d8c92e1fcbd4ac5244b2961fb683b94 — `docs(domain): fill in bounded contexts, entities, and event catalog`
- https://github.com/code-corhuila/library-docs/commit/d97e680d1b60eefbab5a92298b7b1a226f50c9e2 — `docs(context): fill in system overview, scope, and glossary`
