<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       02-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 02

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Oscar Mauricio Areiza Paramo
- GITHUB_USER: OscarAreiza
- TEAM: lms-library
- SPRINT_GOAL: Define the team's agile working method (Kanban & Scrum) and draft the initial LMS Library product brief (tech stack, scope, needs, expected process, glossary)
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| N/A | No formal user stories yet — the initial backlog was defined the following week (see `03-week`) | — | — |

## 2. My individual contribution
- Wrote `AGILE.md`: comparison of Kanban and Scrum as the team's reference working methodologies.
- Wrote `LMS_LIBRARY.md`: the initial product brief for the Loan Management System (LMS) — tech stack (React, Go, PostgreSQL, Docker), context, scope, needs/problems, expected process, open questions, and business glossary.
- Reviewed the Git Flow reference (feat-per-environment, hotfix, cherry-pick) used for the code repositories.

## 3. Blockers and risks
- Several open questions remained in the product brief (user roles, loan period, renewals, fines, reservations) — resolved the following week when defining the user stories.

## 4. Plan for next week
- Start the `library-docs` SDD documentation (`01-context`, `02-domain`) and define the initial user story backlog with acceptance criteria.

## 5. Compliance self-check
- [ ] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [ ] Testable acceptance criteria
- [ ] Tests added/updated (unit / integration)
- [ ] DDD / hexagonal boundaries respected (domain has no I/O)
- [ ] No secrets; config via environment variables

> Not applicable yet — this week produced planning/brief documents, not code or user stories.

## 6. Evidence links
- `02-week/hu-status/AGILE.md`
- `02-week/hu-status/LMS_LIBRARY.md`
