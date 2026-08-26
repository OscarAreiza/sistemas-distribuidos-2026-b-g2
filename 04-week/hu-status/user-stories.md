# User Stories — Backlog — LMS-LIBRARY-V1

---

> **ID status:** `HU-01`..`HU-09` is the canonical ID scheme for this backlog — no service prefix, epics group related stories instead (see Epics table below).
> **Persona note:** all stories use the *system administrator* as the sole actor. The operating model confirms a single Administrator role manages the whole system (no student self-registration, no multi-role access).
> **Sprint/estimate note:** Story Points, sprint assignment, and "Assigned to" below are **proposed defaults** based on relative complexity and dependency order, not a confirmed team decision — adjust freely in Sprint Planning.
> **Policy note:** core loan policy referenced throughout — fixed 7-day loan period, no renewals, maximum 2 active loans per student, flat 7-day suspension on any late return.

## Backlog status

| Cut | Sprint | Total HUs | Refined | In progress | Completed |
|---|---|---|---|---|---|
| Cut 1 | Sprint 1-2 | 5 | 5 | 0 | 0 |
| Cut 2 | Sprint 3-4 | 4 | 4 | 0 | 0 |

Cut 1 covers the Must Have core loan lifecycle (authentication, student and book registration, loan registration, return registration). Cut 2 covers Should Have management/reporting conveniences (student editing/deactivation, book editing, catalog search, penalty reporting).

## Epics

| ID | Epic | Description |
|---|---|---|
| EP-001 | Security & Authentication | Authenticate the administrator and protect access to the system |
| EP-002 | Student Management | Register and maintain library-student records |
| EP-003 | Inventory Management | Register, edit, search, and monitor availability of the book catalog |
| EP-004 | Loans, Returns & Penalties | Register loans and returns, and enforce the suspension policy for late returns |

---

## HU-01 — Administrator Authentication {#HU-01}

**Epic:** EP-001

> **As** the system administrator
> **I want** to log in with a secure username and password
> **so that** unauthorized users cannot access the library management functions

**Acceptance Criteria:**
```gherkin
Scenario 1: Log in with valid credentials
  Given the pre-configured administrator account exists in the backend/database
  When  the administrator submits the correct username and password
  Then  a secure session (JWT) is created
  And   the administrator is redirected to the control panel

Scenario 2: Reject invalid credentials
  Given an administrator submits a username and password
  When  the submitted credentials do not match the stored account
  Then  the system shows a generic error message
  And   the message does not reveal whether the username or the password was incorrect
```

**Definition of Done:**
- [ ] Code reviewed and approved
- [ ] Unit tests written
- [ ] Acceptance criteria verified (manual or automated)
- [ ] API contract updated if applicable
- [ ] Deployed to staging

| Field | Value |
|---|---|
| Story Points | 5 |
| Priority | Must Have |
| Target sprint | Sprint 1 |
| Assigned to | [Unassigned] |
| Status | Backlog |
| Dependencies | — |
| Affected service(s) | Auth |

---

## HU-02 — Student Registration {#HU-02}

**Epic:** EP-002

> **As** the administrator
> **I want** to register new students with their basic information (name, ID, email, phone)
> **so that** they are enabled in the system and can be linked to book loans

**Acceptance Criteria:**
```gherkin
Scenario 1: Register a student with valid data
  Given the administrator provides full name, unique ID/document, institutional email, and phone
  When  the administrator submits the student registration
  Then  the student is stored
  And   the student appears in the general student listing

Scenario 2: Reject duplicate document ID
  Given a student with document ID "1075300000" already exists in the system
  When  the administrator submits a new student with the same document ID
  Then  the system rejects the registration
  And   an error message indicates the document ID is already registered
```

**Definition of Done:**
- [ ] Code reviewed and approved
- [ ] Unit tests written
- [ ] Acceptance criteria verified (manual or automated)
- [ ] API contract updated if applicable
- [ ] Deployed to staging

| Field | Value |
|---|---|
| Story Points | 3 |
| Priority | Must Have |
| Target sprint | Sprint 1 |
| Assigned to | [Unassigned] |
| Status | Backlog |
| Dependencies | — |
| Affected service(s) | Students |

---

## HU-03 — Student Search, Editing & Deactivation {#HU-03}

**Epic:** EP-002

> **As** the administrator
> **I want** to search, edit, or deactivate student records
> **so that** contact information stays current and I control who is actively entitled to the service

**Acceptance Criteria:**
```gherkin
Scenario 1: Edit a student's contact information
  Given a registered student
  When  the administrator updates their contact information
  Then  the student record reflects the new values

Scenario 2: Block deactivation of a student with active loans or an active suspension
  Given a registered student with an active loan or an active suspension
  When  the administrator attempts to deactivate that student
  Then  the system rejects the deactivation
  And   indicates that the student has active loans or an active suspension
```

**Definition of Done:**
- [ ] Code reviewed and approved
- [ ] Unit tests written
- [ ] Acceptance criteria verified (manual or automated)
- [ ] API contract updated if applicable
- [ ] Deployed to staging

| Field | Value |
|---|---|
| Story Points | 3 |
| Priority | Should Have |
| Target sprint | Sprint 3 |
| Assigned to | [Unassigned] |
| Status | Backlog |
| Dependencies | HU-02 |
| Affected service(s) | Students |

---

## HU-04 — Book Registration {#HU-04}

**Epic:** EP-003

> **As** the administrator
> **I want** to register new books with title, author, ISBN, category, year, and total copies
> **so that** the library catalog is populated

**Acceptance Criteria:**
```gherkin
Scenario 1: Register a book with valid data
  Given the administrator provides a title, author, ISBN, category, publication year, and number of copies
  When  the administrator submits the book registration
  Then  the book is stored
  And   the book's initial availability is set equal to the total copies entered

Scenario 2: Reject duplicate ISBN
  Given a book with ISBN "978-3-16-148410-0" already exists in the catalog
  When  the administrator submits a new book with the same ISBN
  Then  the system rejects the registration
  And   an error message indicates the ISBN is already registered
```

**Definition of Done:**
- [ ] Code reviewed and approved
- [ ] Unit tests written
- [ ] Acceptance criteria verified (manual or automated)
- [ ] API contract updated if applicable
- [ ] Deployed to staging

| Field | Value |
|---|---|
| Story Points | 3 |
| Priority | Must Have |
| Target sprint | Sprint 1 |
| Assigned to | [Unassigned] |
| Status | Backlog |
| Dependencies | — |
| Affected service(s) | Catalog |

---

## HU-05 — Inventory Control & Book Search {#HU-05}

**Epic:** EP-003

> **As** the administrator
> **I want** to search the catalog by title, author, ISBN, or category and see real-time availability
> **so that** I can quickly tell a student whether a book is available

**Acceptance Criteria:**
```gherkin
Scenario 1: Search the catalog by title, author, ISBN, or category
  Given registered books in the catalog
  When  the administrator searches or filters by title, author, ISBN, or category
  Then  matching books are returned with their total and currently available copies

Scenario 2: Search with no matches
  Given no registered book matches the search term or filter
  When  the administrator executes the search
  Then  the system returns an empty result set, not an error
```

**Definition of Done:**
- [ ] Code reviewed and approved
- [ ] Unit tests written
- [ ] Acceptance criteria verified (manual or automated)
- [ ] API contract updated if applicable
- [ ] Deployed to staging

| Field | Value |
|---|---|
| Story Points | 2 |
| Priority | Should Have |
| Target sprint | Sprint 3 |
| Assigned to | [Unassigned] |
| Status | Backlog |
| Dependencies | HU-04 |
| Affected service(s) | Catalog |

---

## HU-09 — Book Editing {#HU-09}

**Epic:** EP-003

> **As** the administrator
> **I want** to edit a registered book's information (title, author, category, year)
> **so that** the catalog stays accurate without losing the book's existing loan history

**Acceptance Criteria:**
```gherkin
Scenario 1: Edit a book's information
  Given a registered book
  When  the administrator updates its title, author, category, or year
  Then  the book record reflects the new values
  And   existing and past loans referencing this book are unaffected

Scenario 2: ISBN is not editable through this action
  Given a registered book
  When  the administrator edits its title, author, category, or year
  Then  the ISBN is not part of the editable fields
  And   the book's ISBN remains unchanged after the update
```

> ISBN is the book's stable business key (`02-domain/entities-and-rules.md`, modeling note on
> Book) — it is set once at registration (HU-04) and is not part of this action's editable
> fields, so there is no "duplicate ISBN on edit" case to handle.

**Definition of Done:**
- [ ] Code reviewed and approved
- [ ] Unit tests written
- [ ] Acceptance criteria verified (manual or automated)
- [ ] API contract updated if applicable
- [ ] Deployed to staging

| Field | Value |
|---|---|
| Story Points | 2 |
| Priority | Should Have |
| Target sprint | Sprint 3 |
| Assigned to | [Unassigned] |
| Status | Backlog |
| Dependencies | HU-04 |
| Affected service(s) | Catalog |

---

## HU-06 — Loan Registration {#HU-06}

**Epic:** EP-004

> **As** the administrator
> **I want** to register a book loan to a registered student with a due date
> **so that** there is a formal record of the physical loan

**Acceptance Criteria:**
```gherkin
Scenario 1: Register a loan when eligible
  Given an active student with fewer than 2 active loans and no active suspension, and a book with at least one available copy
  When  the administrator registers a loan
  Then  the loan is created with a start date and a due date 7 days later
  And   the book's available-copy count decreases by one

Scenario 2: Reject a loan for a suspended or over-limit student
  Given a student with an active suspension, or with 2 active loans already registered
  When  the administrator attempts to register a new loan for that student
  Then  the system rejects the loan
  And   indicates the reason (active suspension or maximum active loans reached)
```

**Definition of Done:**
- [ ] Code reviewed and approved
- [ ] Unit tests written
- [ ] Acceptance criteria verified (manual or automated)
- [ ] API contract updated if applicable
- [ ] Deployed to staging

| Field | Value |
|---|---|
| Story Points | 5 |
| Priority | Must Have |
| Target sprint | Sprint 2 |
| Assigned to | [Unassigned] |
| Status | Backlog |
| Dependencies | HU-01, HU-02, HU-04 |
| Affected service(s) | Loans |

---

## HU-07 — Return Registration & History Tracking {#HU-07}

**Epic:** EP-004

> **As** the administrator
> **I want** to register the return of a loaned book
> **so that** the copy is added back to available inventory and the loan cycle is closed

**Acceptance Criteria:**
```gherkin
Scenario 1: Register an on-time return
  Given an active loan whose return date is on or before its due date
  When  the administrator registers its return
  Then  the loan moves from "Active" to "Returned" in the history
  And   the book's available-copy count increases by one

Scenario 2: Register a late return
  Given an active loan whose return date is after its due date
  When  the administrator registers its return
  Then  the system marks the return as late
  And   the loan moves to "Returned" with the late status recorded in the history
```

**Definition of Done:**
- [ ] Code reviewed and approved
- [ ] Unit tests written
- [ ] Acceptance criteria verified (manual or automated)
- [ ] API contract updated if applicable
- [ ] Deployed to staging

| Field | Value |
|---|---|
| Story Points | 5 |
| Priority | Must Have |
| Target sprint | Sprint 2 |
| Assigned to | [Unassigned] |
| Status | Backlog |
| Dependencies | HU-06 |
| Affected service(s) | Returns |

---

## HU-08 — Late-Return Penalty System {#HU-08}

**Epic:** EP-004

> **As** the administrator
> **I want** the system to detect overdue loans and automatically apply a suspension
> **so that** late returns are discouraged

**Acceptance Criteria:**
```gherkin
Scenario 1: Apply a suspension after a late return
  Given a loan is registered as returned late (per HU-07)
  When  the return is processed
  Then  a flat 7-day suspension is applied to the student
  And   the suspension is visible on the student's profile

Scenario 2: Identify loans currently overdue
  Given an active loan whose due date has passed and it has not yet been returned
  When  the administrator views the "Overdue Loans" report
  Then  that loan appears in the list
```

**Definition of Done:**
- [ ] Code reviewed and approved
- [ ] Unit tests written
- [ ] Acceptance criteria verified (manual or automated)
- [ ] API contract updated if applicable
- [ ] Deployed to staging

| Field | Value |
|---|---|
| Story Points | 5 |
| Priority | Should Have |
| Target sprint | Sprint 4 |
| Assigned to | [Unassigned] |
| Status | Backlog |
| Dependencies | HU-07 |
| Affected service(s) | Penalties |

---

## Rules for writing HUs

**1. The role matters.** Do not write "As a user" — that says nothing. Use the specific role:
✓ As the administrator · ✗ As a user · ✗ As a person

**2. The benefit justifies the work.** The "so that" must describe a business benefit, not redescribe the action:
✓ so that they are enabled in the system and can be linked to book loans · ✗ so that I can register the student (only restates the feature)

**3. ACs are verifiable.** Each AC must be verifiable manually or automatable as a test:
✓ Then the loan is created with a due date 7 days later · ✗ Then the system works well (not verifiable)

**4. One HU = one unit of value.** If an HU has 15 ACs, it is probably several HUs. The team must be able to complete it within one sprint.

## Ready-to-copy HU template

```
### HU-0X — [Name] {#HU-0X}
**Epic:** EP-00X
> **As** [role]
> **I want** [action]
> **so that** [benefit]
**Acceptance Criteria:**
```gherkin
Scenario 1: [name]
  Given [context]
  When  [action]
  Then  [result]
```
| Field | Value |
|-------|-------|
| Story Points | |
| Priority | |
| Target sprint | |
| Status | Backlog |
| Dependencies | |
```

## Correlations

- System overview, scope, glossary → `01-context/`
- Bounded contexts, entities, business rules, domain events → `02-domain/`
- Definition of Ready → `00-governance/definition-of-ready.md`
- Definition of Done → `00-governance/definition-of-done.md`
- Functional requirements derived from these HUs → `04-requirements/functional.md`
- Traceability (HU ↔ FR ↔ tests) → `04-requirements/traceability-matrix.md`
