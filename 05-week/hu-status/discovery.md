# Discovery — LMS-LIBRARY-V1

> Deliverable for the "Project Discovery" issue: consolidates the problem definition, target
> users, market research, MVP scope, initial data model, and main user flows into one document,
> pulled from (and kept in sync with) the detailed source documents in this repo. This file is a
> summary with links — the detailed, maintained version of each section always lives in its own
> folder; if this file and a linked folder ever disagree, the folder wins.

---

## 1. The problem

University library staff manage the book loan process through paper records, spreadsheets, or
disconnected tools, with no centralized, real-time view of which books are available, who has
which book on loan, or which loans are overdue. This causes delays, registration errors, and a
lack of control over overdue material.

Full problem statement, evidence, and impact → [`03-product/problem-framing.md`](03-product/problem-framing.md)

## 2. Target users

| Persona | Role in the system |
|---------|---------------------|
| **Administrator** (primary persona) | Library staff member — the only authenticated role. Registers books/students, runs the loan/return cycle, applies suspensions. |
| Student, faculty, staff | Indirect — never log in; the Administrator manages their records and acts on their behalf. |

Full persona detail and jobs-to-be-done → [`03-product/problem-framing.md`](03-product/problem-framing.md), [`03-product/discovery-brief.md`](03-product/discovery-brief.md)

## 3. Market research (similar solutions)

Reviewed against Koha, Evergreen ILS, Follett Destiny/Aspen, Alexandria, and Accessit
Library/LibraryWorld. None fit this project's academic scope, timeline, or budget — every one
assumes multi-role/self-service access and/or a fines engine, which is exactly what this
project's MVP deliberately excludes. Full comparison table and what it validated →
[`03-product/discovery-brief.md`](03-product/discovery-brief.md)

## 4. MVP scope — Cut 1

**In scope (v1 MVP):** Administrator authentication, book catalog management, student registry,
loan management, return management, overdue tracking & suspensions, loan history & queries.

**Explicitly out of scope (v1):** student self-registration/login, multiple roles (RBAC), loan
renewals, book reservations/holds, monetary fines, exportable reports, due-date notifications,
physical-condition tracking.

Full in/out tables, assumptions, and constraints → [`01-context/scope.md`](01-context/scope.md)
Prioritized backlog → [`03-product/product-backlog.md`](03-product/product-backlog.md)

## 5. Initial data model

Four entities: `Book`, `Student`, `Loan`, `Administrator` — relational schema, one table per
entity, with `Loan` referencing `Student` and `Book`. See the ER diagram and full column-level
schema → [`06-data/models.md`](06-data/models.md); field-by-field business meaning →
[`06-data/data-dictionary.md`](06-data/data-dictionary.md)

## 6. Main user flows

- **Authentication** — login with username/password, generic error on failure.
- **Register a student, then register a book** — the two independent registration flows that
  feed the loan flow below.
- **Loan and return** (Circulation context — domain modeled, endpoint not yet implemented) —
  register a loan, then later register its return; a late return applies a 7-day suspension.

Full flow diagrams and the screen map they connect → [`12-ux-ui/navigation-map.md`](12-ux-ui/navigation-map.md)

## 7. Wireframes

🔴 **Pending visual design.** The screens that need a wireframe, in priority order, are listed
and ready to receive the actual mockup links → [`12-ux-ui/wireframes.md`](12-ux-ui/wireframes.md)

---

## Acceptance criteria checklist (Discovery issue)

- [x] Problem clearly defined and documented — `03-product/problem-framing.md`
- [x] MVP scope agreed — `01-context/scope.md`, `03-product/product-backlog.md`
- [ ] Wireframes of the main screens ready — tracker exists (`12-ux-ui/wireframes.md`), visual mockups still pending
- [x] Initial data model documented — `06-data/models.md`, `06-data/data-dictionary.md`
- [x] Documentation uploaded to the docs repo — this file and all links above live in `library-docs`

## Correlations

- Problem framing → `03-product/problem-framing.md`
- Discovery research (persona, market analysis, validated assumptions) → `03-product/discovery-brief.md`
- Product vision → `03-product/vision.md`
- Confirmed scope → `01-context/scope.md`
- Domain model → `02-domain/entities-and-rules.md`
- Data model → `06-data/models.md`
- Navigation and user flows → `12-ux-ui/navigation-map.md`
- Technology stack decision → [`stack.md`](stack.md)
