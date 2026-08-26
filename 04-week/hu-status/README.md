<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       04-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 04

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Oscar Mauricio Areiza Paramo
- GITHUB_USER: OscarAreiza
- TEAM: lms-library
- SPRINT_GOAL: Fill in 01-context through 06-data of the SDD documentation (library-docs), including the derived requirements/data artifacts (user stories, functional/non-functional requirements, data dictionary), stand up the lms-library codebase (Go hexagonal backend + React frontend base), and implement HU-01 (Administrator Login) end to end.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-01 | Administrator Authentication | doing | https://github.com/OscarAreiza/lms-library/commit/b9929bc (branch `feat/HU-01-admin-login`, pushed — pending PR review/merge to `dev`) |

## 2. My individual contribution
- Filled in `library-docs` folders `01-context` through `06-data` following each folder's own README format, resolving several format/coherence gaps found along the way (missing fields, ID scheme mismatches, a contradiction between HU-09's acceptance criteria and the already-decided ISBN-immutable design).
- Produced the derived requirements/data artifacts specifically requested for this status:
  - `04-requirements/user-stories.md` — 9 formalized HUs (HU-01..HU-09) with Gherkin acceptance criteria.
  - `04-requirements/functional.md` — 20 functional requirements (FR-001..FR-020), each traced to its source HU.
  - `04-requirements/non-functional.md` — 8 non-functional requirement categories (NFR-001..NFR-008) with measurable metrics.
  - `06-data/data-dictionary.md` — field-by-field business meaning for every table backing the domain model.
- Set up the `lms-library` code repository: branch strategy (`main -> QA -> dev -> feat/HU-XX-...`) per `00-governance/git-conventions.md`, `.github` issue/PR templates linking each change to its HU, and the base scaffolding (hexagonal Go backend, React + Vite + Tailwind frontend, Docker Compose with automatic migrations) — pushed to `dev`.
- Implemented HU-01 end to end: JWT-based login (`AdministratorRepository`, `JWTIssuer`, `Login` use case, `AuthHandler`), seeded a local dev administrator, and verified the full flow live (Docker) — login returns a valid token and the SPA reaches `/dashboard`.

## 3. Blockers and risks
- `00-governance/agile-conventions.md` still has unfilled placeholders (sprint length, ceremony times, backlog tool) — needs a real team decision, not something I can fill in alone.
- All 9 HUs currently live on separate, unmerged feature branches (`feat/HU-01`..`feat/HU-09`) — no PR has been opened yet, so none has gone through code review.
- The 4 bounded-context modules still share a single PostgreSQL database in v1 (tracked as technical debt `AT-001`/`AT-002` in `05-architecture/overview.md`) — acceptable for now, but will need attention if the project ever needs to split a module into its own service.

## 4. Plan for next week
- Open PRs for `feat/HU-01-admin-login` (and progressively the rest) into `dev`, get them reviewed, and merge.
- Once Cut 1 (Must Have: HU-01, HU-02, HU-04, HU-06, HU-07) is merged and stable in `dev`, promote `dev -> QA` for validation.
- Fill in the real Scrum/Kanban details in `00-governance/agile-conventions.md` with the team.

## 5. Compliance self-check
- [x] Conventional Commits - `type(scope): summary`
- [ ] Per-environment HU branch + PR to that environment (hu-xxx-dev -> develop, ...)
- [x] Testable acceptance criteria
- [x] Tests added/updated (unit / integration)
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links
- Documentation (01-context..05-architecture): https://github.com/code-corhuila/library-docs/commit/98ed960
- Documentation (06-data, incl. data-dictionary.md): https://github.com/code-corhuila/library-docs/commit/104a614
- Requirements fix (user-stories.md / functional.md HU-09 correction): https://github.com/code-corhuila/library-docs/commit/509e366
- lms-library base scaffolding on `dev`: https://github.com/OscarAreiza/lms-library/commit/6457c2b
- HU-01 implementation: https://github.com/OscarAreiza/lms-library/commit/b9929bc
- Local copies attached in this same folder (`04-week/hu-status/`):
  - `user-stories.md` — HU-01..HU-09 with Gherkin acceptance criteria
  - `functional.md` — FR-001..FR-020
  - `non-functional.md` — NFR-001..NFR-008
  - `data-dictionary.md` — field-by-field data dictionary
  - `mvp.zip` — full snapshot of `feat/HU-01-admin-login` (base scaffolding + HU-01 code)
