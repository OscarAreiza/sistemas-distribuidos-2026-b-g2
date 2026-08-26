# Data Dictionary — LMS-LIBRARY-V1

> Exact business meaning of every field that could otherwise be ambiguous. If this dictionary
> disagrees with the code, the dictionary is wrong and must be corrected — not the other way
> around (per `00-sdd-guide.md`, living documentation).

| Field | Module | Table | Type | Detailed description | Possible values |
|-------|--------|-------|------|----------------------|-----------------|
| `username` | Access | `administrators` | VARCHAR | Login identifier for the Administrator. Not an email — see `01-context/glossary.md` (Administrator). | Any unique string, 1–100 chars |
| `password_hash` | Access | `administrators` | VARCHAR | bcrypt hash (cost ≥ 12) of the Administrator's password. Never the plaintext password. | bcrypt hash string |
| `document_id` | Membership | `students` | VARCHAR | The student's institutional ID/document code — the business key used at the circulation desk (not the surrogate `id`). | Any unique string, per HU-02 |
| `suspended_until` | Membership | `students` | TIMESTAMPTZ \| NULL | Date until which the student cannot receive a new loan. Set to `loan_return_time + 7 days` automatically on a late return (INV-006 on Loan). `NULL` or a past date means the student is currently eligible. | `NULL`, or any future timestamp |
| `deactivated_at` | Membership | `students` | TIMESTAMPTZ \| NULL | When the Administrator deactivated this student (HU-03). `NULL` = active. Deactivation is blocked while the student has an active loan or an active suspension (INV-002 on Student). | `NULL`, or a past timestamp |
| `isbn` | Catalog | `books` | VARCHAR | International Standard Book Number identifying a specific edition. Unique per book record — see `01-context/glossary.md` (ISBN). | Standard ISBN-10/ISBN-13 format |
| `total_copies` | Catalog | `books` | INTEGER | Physical copies of this book registered in the library, as a count — there is no per-copy entity in v1 (`02-domain/entities-and-rules.md`, modeling note). | Integer ≥ 1 |
| `available_copies` | Catalog | `books` | INTEGER | Copies of this book not currently tied to an active loan. Recalculated on every loan/return, never edited directly by an Administrator. | Integer, `0 ≤ available_copies ≤ total_copies` |
| `category` | Catalog | `books` | VARCHAR | Free-text subject classification used for catalog search (HU-05). Not a fixed enum in v1 — no controlled vocabulary was defined in the product brief. | Any string, e.g. `"Fiction"`, `"Computer Science"` |
| `loan_date` | Circulation | `loans` | TIMESTAMPTZ | When the copy was physically handed to the student. Set once at creation, immutable. | Any timestamp ≤ `now()` at creation |
| `due_date` | Circulation | `loans` | TIMESTAMPTZ | Always `loan_date + 7 days` — computed, never independently settable (INV-001 on Loan; no renewals in v1, `01-context/scope.md`). | `loan_date + 7 days` |
| `return_date` | Circulation | `loans` | TIMESTAMPTZ \| NULL | When the copy was actually handed back. `NULL` while the loan is `ACTIVE`. | `NULL`, or a timestamp ≥ `loan_date` |
| `status` | Circulation | `loans` | VARCHAR | Lifecycle state of the loan. | `ACTIVE`, `RETURNED` |
| `was_late` | Circulation | `loans` | BOOLEAN \| NULL | Whether `return_date > due_date` at the moment the return was registered. `NULL` while `status = ACTIVE`. Drives the `StudentSuspended` policy (`02-domain/domain-events.md`). | `NULL`, `true`, `false` |
| "Overdue loan" (derived, not a column) | Circulation | `loans` | — | A loan is overdue when `status = 'ACTIVE' AND due_date < now()`. Never stored as its own status — see `06-data/models.md`, Circulation modeling decisions. | Computed at query time |

---

## Correlations

- Schema these fields belong to → `06-data/models.md`
- Business rules these fields enforce → `02-domain/entities-and-rules.md`
- Canonical term definitions → `01-context/glossary.md`
