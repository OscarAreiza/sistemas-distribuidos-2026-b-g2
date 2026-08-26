# Functional Requirements (FR) — LMS-LIBRARY-V1

> Each FR is derived from an acceptance criterion in `04-requirements/user-stories.md`.
> Module names match each HU's "Affected service(s)" field.

| ID | Module | Description | Source (HU) | Priority |
|----|--------|-------------|------------|---------|
| FR-001 | Auth | The system must authenticate the administrator with username and password, creating a secure (JWT) session on success | HU-01 | Must Have |
| FR-002 | Auth | The system must reject invalid credentials with a generic error that does not reveal whether the username or the password was wrong | HU-01 | Must Have |
| FR-003 | Students | The system must allow registering a new student with full name, document ID, institutional email, and phone | HU-02 | Must Have |
| FR-004 | Students | The system must reject student registration when the document ID is already registered | HU-02 | Must Have |
| FR-005 | Students | The system must allow searching and editing a registered student's contact information | HU-03 | Should Have |
| FR-006 | Students | The system must block deactivation of a student who has an active loan or an active suspension | HU-03 | Should Have |
| FR-007 | Catalog | The system must allow registering a new book with title, author, ISBN, category, year, and total copies, setting initial availability equal to total copies | HU-04 | Must Have |
| FR-008 | Catalog | The system must reject book registration when the ISBN is already registered | HU-04 | Must Have |
| FR-009 | Catalog | The system must allow searching the catalog by title, author, ISBN, or category and return each match's total and available copies | HU-05 | Should Have |
| FR-010 | Catalog | The system must return an empty result set (not an error) when a catalog search has no matches | HU-05 | Should Have |
| FR-011 | Catalog | The system must allow editing a registered book's title, author, category, and year without affecting its existing loan history | HU-09 | Should Have |
| FR-012 | Catalog | The system must keep a book's ISBN immutable through the edit action — ISBN is set once at registration (HU-04) and is not part of the editable fields | HU-09 | Should Have |
| FR-013 | Loans | The system must allow registering a loan for an eligible student (not suspended, fewer than 2 active loans) and an available book, setting the due date to 7 days after the loan date | HU-06 | Must Have |
| FR-014 | Loans | The system must decrement the book's available-copy count by one when a loan is registered | HU-06 | Must Have |
| FR-015 | Loans | The system must reject registering a loan for a suspended student or a student with 2 active loans, indicating the reason | HU-06 | Must Have |
| FR-016 | Returns | The system must allow registering the return of a loaned book, moving the loan from Active to Returned | HU-07 | Must Have |
| FR-017 | Returns | The system must increment the book's available-copy count by one when a return is registered | HU-07 | Must Have |
| FR-018 | Returns | The system must detect and record whether a return was late (return date after due date) | HU-07 | Must Have |
| FR-019 | Penalties | The system must automatically apply a fixed 7-day suspension to a student when their return is registered as late | HU-08 | Should Have |
| FR-020 | Penalties | The system must allow querying all currently overdue loans (active, due date already passed) | HU-08 | Should Have |

---

## Correlations

- Source user stories → `04-requirements/user-stories.md`
- Non-functional requirements (quality attributes) → `04-requirements/non-functional.md`
- Traceability (FR ↔ HU ↔ tests) → `04-requirements/traceability-matrix.md`
- FRs are grouped by responsible service in `09-microservices/` once that section is filled in
