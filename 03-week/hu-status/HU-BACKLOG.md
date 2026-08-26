# Session Log — Documentation Kickoff & User Stories Definition

**Date:** 2026-08-22
**Team:** Oscar Areiza (Tech Lead), Hermes Pascuas (DevOps), Luis Alejandro Meneses (DevOps)

## What was done today

- Started filling in the `library-docs` repository following the SDD methodology: completed
  `01-context` (system overview, scope, glossary) and `02-domain` (bounded contexts, entities
  and business rules, domain event catalog).
- Updated the tech stack decision: frontend changed from Vue.js to **React** (backend Go,
  database PostgreSQL, Docker unchanged).
- Confirmed the operating model: a single **Administrator** role manages the whole system
  (no student self-registration, no multi-role access).
- Confirmed the core loan policy used across the domain model and the user stories below:
  fixed 7-day loan period, no renewals, maximum 2 active loans per student, and a flat
  7-day suspension applied on any late return.
- Defined the initial Product Backlog: 8 user stories grouped into 4 modules, each with
  acceptance criteria and checked against a shared Definition of Ready / Definition of Done.

---

## Product Backlog — User Stories

### Epics

| ID | Epic | Modules covered |
|----|------|-----------------|
| EP-001 | Security & Authentication | HU-01 |
| EP-002 | Student Management | HU-02, HU-03 |
| EP-003 | Inventory Management | HU-04, HU-05 |
| EP-004 | Loans, Returns & Penalties | HU-06, HU-07, HU-08 |

---

### HU-01 — Administrator Authentication

**Epic:** EP-001

> **As** the system administrator
> **I want** to log in with a secure username and password
> **so that** unauthorized users cannot access the library management functions

**Acceptance criteria:**
- The system has a single pre-configured administrator account in the backend/database.
- Correct credentials create a secure session (JWT) and redirect to the control panel.
- Incorrect credentials show a generic error message, without revealing whether the username
  or the password was wrong.
- Passwords are stored using a strong hashing algorithm (bcrypt).

| Field | Value |
|-------|-------|
| Priority | Must Have |
| Status | Backlog |
| Dependencies | None |

---

### HU-02 — Student Registration

**Epic:** EP-002

> **As** the administrator
> **I want** to register new students with their basic information (name, ID, email, phone)
> **so that** they are enabled in the system and can be linked to book loans

**Acceptance criteria:**
- The registration form requires: full name, unique ID/document, institutional email, and phone.
- The system validates that the document ID is not already registered.
- On success, the student appears in the general student listing.

| Field | Value |
|-------|-------|
| Priority | Must Have |
| Status | Backlog |
| Dependencies | None |

---

### HU-03 — Student Search, Editing & Deactivation

**Epic:** EP-002

> **As** the administrator
> **I want** to search, edit, or deactivate student records
> **so that** contact information stays current and I control who is actively entitled to the service

**Acceptance criteria:**
- There is a search/listing view of registered students.
- The administrator can update a student's contact information.
- The system blocks deletion of a student who has active loans or an active suspension.

| Field | Value |
|-------|-------|
| Priority | Should Have |
| Status | Backlog |
| Dependencies | HU-02 |

---

### HU-04 — Book Registration

**Epic:** EP-003

> **As** the administrator
> **I want** to register new books with title, author, ISBN, category, year, and total copies
> **so that** the library catalog is populated

**Acceptance criteria:**
- Registration requires: title, author, unique ISBN, category, publication year, and number
  of copies.
- The system automatically sets initial availability equal to the total copies entered.

| Field | Value |
|-------|-------|
| Priority | Must Have |
| Status | Backlog |
| Dependencies | None |

---

### HU-05 — Inventory Control & Book Search

**Epic:** EP-003

> **As** the administrator
> **I want** to search the catalog by title, author, ISBN, or category and see real-time availability
> **so that** I can quickly tell a student whether a book is available

**Acceptance criteria:**
- The catalog offers a text search and category filters.
- For each book, the interface shows total copies and currently available copies.

| Field | Value |
|-------|-------|
| Priority | Should Have |
| Status | Backlog |
| Dependencies | HU-04 |

---

### HU-06 — Loan Registration

**Epic:** EP-004

> **As** the administrator
> **I want** to register a book loan to a registered student with a due date
> **so that** there is a formal record of the physical loan

**Acceptance criteria:**
- The administrator selects an active student and a book with at least one available copy.
- Confirming the loan decreases the book's availability by one.
- The system records the loan start date and computes the due date as 7 days later (fixed
  policy, no renewals).
- The system blocks the loan if the student has an active suspension.
- The system blocks the loan if the student already has 2 active loans (maximum allowed).

| Field | Value |
|-------|-------|
| Priority | Must Have |
| Status | Backlog |
| Dependencies | HU-02, HU-04 |

---

### HU-07 — Return Registration & History Tracking

**Epic:** EP-004

> **As** the administrator
> **I want** to register the return of a loaned book
> **so that** the copy is added back to available inventory and the loan cycle is closed

**Acceptance criteria:**
- The administrator finds the active loan and marks it as returned.
- Book availability increases by one.
- The loan moves from "Active" to "Returned" in the history.
- The system compares the return date with the due date to determine whether the return
  was late.

| Field | Value |
|-------|-------|
| Priority | Must Have |
| Status | Backlog |
| Dependencies | HU-06 |

---

### HU-08 — Late-Return Penalty System

**Epic:** EP-004

> **As** the administrator
> **I want** the system to detect overdue loans and automatically apply a suspension
> **so that** late returns are discouraged

**Acceptance criteria:**
- The system shows an "Overdue Loans" indicator/report for loans past their due date.
- A late return triggers a flat 7-day suspension on new loans for that student (fixed
  policy — not proportional to days late, no monetary fine).
- The administrator can see a student's active suspension on their profile.

| Field | Value |
|-------|-------|
| Priority | Should Have |
| Status | Backlog |
| Dependencies | HU-07 |

---

## Definition of Ready (applies to every HU above)

- [ ] Written as **As / I want / so that**, with a specific role (not "as a user")
- [ ] At least 2 acceptance criteria, covering the happy path and one error/edge case
- [ ] Dependencies on other HUs are identified (see table per HU)
- [ ] The team has discussed and agreed on the story before it enters a sprint

## Definition of Done (applies to every HU above)

- [ ] All acceptance criteria are implemented and verified
- [ ] Code reviewed by at least one teammate
- [ ] Unit tests written for the new business logic
- [ ] Related documentation updated (`library-docs` repo) if behavior changed

---

## Status summary

| HU ID | Title | Status |
|-------|-------|--------|
| HU-01 | Administrator Authentication | Backlog |
| HU-02 | Student Registration | Backlog |
| HU-03 | Student Search, Editing & Deactivation | Backlog |
| HU-04 | Book Registration | Backlog |
| HU-05 | Inventory Control & Book Search | Backlog |
| HU-06 | Loan Registration | Backlog |
| HU-07 | Return Registration & History Tracking | Backlog |
| HU-08 | Late-Return Penalty System | Backlog |

## Next step

Move these user stories into `library-docs/04-requirements/user-stories.md`, following the
SDD weekly order, once the team is ready to continue past `02-domain`.
