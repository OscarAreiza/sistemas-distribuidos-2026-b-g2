# ADR-010 — Liquibase for every `-db` repository's migrations

- **ID:** ADR-010
- **Date:** 2026-09-26
- **Status:** Accepted
- **Authors:** Oscar Areiza — Tech Lead
- **Reviewers:** Hermes Pascuas, Luis Alejandro Meneses — Development team

---

## Context

`rules/1-Norma-Repositorios-Sistemas-Distribuidos-2026B.pdf` (numeral 4.2.2) requires every
`-db` repository to manage its schema with **Liquibase or Flyway** — not the ad hoc tool each
repo currently uses. Today, `lms-membership-db` and `lms-catalog-db` each run `golang-migrate`
(numbered `.up.sql`/`.down.sql` files plus a `migrate/migrate` container), and `lms-circulation-db`
has no migration tooling at all: `loans`' indexes are created by `lms-circulation-api` itself at
startup (`internal/infrastructure/mongodb.EnsureIndexes`), which a real review on
`lms-circulation-api#2` already flagged as a critical violation of "database structure lives in
the `-db` repo, never in the `-api`" — now a written rule (`rules/2-anexos/B-db-mongo.md`,
"Qué va en cada parte").

None of the three `-db` repos yet meet Anexo A (`A-db-postgres.md`, PostgreSQL) or Anexo B
(`B-db-mongo.md`, MongoDB): no DDL/DML/DCL/TCL family folders, no rollback mirror, no
reconstruction-verification CI, tables/collections outside a per-domain schema, and (for
Postgres) `VARCHAR(n)` columns where the norm requires `text` with an explicit `CHECK`.

**Known constraints:**
- Same 3-person academic team as `ADR-002`/`ADR-004`/`ADR-005` — one migration tool to learn,
  not two, keeps the operational surface small
- The norm lets the team choose **per domain** (numeral 4.2.2), but nothing in this project's
  three domains has a data shape that favors Flyway specifically (`Anexo A`, "Con Flyway" section:
  a real advantage only where the team wants Flyway's single-file-per-version simplicity and
  is fine hand-writing every rollback, since Flyway Community never reverts on its own)
- `lms-circulation-db` is MongoDB, and the norm is explicit that Mongo's only supported tool is
  **Liquibase with its MongoDB extension** (`Anexo B`) — Flyway is not an option there at all

---

## Decision

**We decided:** all three `-db` repositories (`lms-membership-db`, `lms-catalog-db`,
`lms-circulation-db`) use **Liquibase**, replacing `golang-migrate` in the two PostgreSQL repos
entirely.

**Justification:** `lms-circulation-db` already has no choice — Liquibase is the only tool the
norm allows for MongoDB (`Anexo B`). Once one of the three repos must run Liquibase, using it for
the other two as well means the team operates **one** migration tool, one *changeset* format, one
CI verification shape (`update` → `update` again expecting zero → `rollback-count 999` → `update`
again) across every `-db` repo, instead of splitting attention between Liquibase's YAML changelogs
and Flyway's `V<n>__` file-per-version convention for no functional gain — nothing in
`students` or `books`' shape needs Flyway's simplicity enough to justify a second tool, and
Liquibase's self-contained rollback (declared in the same changeset) is a real advantage over
Flyway Community, which never reverts on its own and would need every `U<n>` script hand-verified.

---

## Evaluated alternatives

| Alternative | Pros | Cons | Reason for discarding |
|------------|------|------|-----------------------|
| **Liquibase everywhere — CHOSEN** | One tool, one team habit, across all three `-db` repos; only option the norm allows for MongoDB anyway; rollback is declared alongside the change it reverts | More verbose than Flyway (a YAML changelog entry plus the SQL file, instead of just a numbered SQL file); the team has to install and learn the MongoDB extension (`liquibase-mongodb` + the `mongodb` driver, both required — `Anexo B`) | — (chosen) |
| Flyway for `lms-membership-db`/`lms-catalog-db`, Liquibase only for `lms-circulation-db` | Postgres migrations stay a single numbered `.sql` file, arguably the simplest format to read | Splits the team across two tools and two CI verification shapes for no functional gain; Flyway Community cannot revert on its own, so every rollback still has to be hand-written as a `U<n>` script — the "simplicity" only applies to the forward migration, not the full lifecycle the norm requires | The two PostgreSQL domains have no data shape that specifically rewards Flyway's format, and running two tools costs more team attention than it saves |
| Keep `golang-migrate` for Postgres, add Liquibase only where MongoDB forces it | Zero migration cost for the two repos already using it | Does not satisfy the norm at all — `golang-migrate` is not Liquibase or Flyway (numeral 4.2.2); no DDL/DML/DCL/TCL families, no rollback mirror, no reconstruction-verification CI shape the norm's "Cómo se verifica" checklist expects | Already disqualified by the norm itself, independent of team preference |

---

## Related decisions this ADR also registers (per `Anexo A`/`Anexo B`)

- **Schema ownership, per domain:** `lms-membership-db` creates and owns schema `membership`;
  `lms-catalog-db` owns schema `catalog`. Neither lives in `public` (`Anexo A`, rule 9). Circulation
  has no relational schema — its "schema" is the `loans` collection's `$jsonSchema` validator in
  `loan_db` (`Anexo B`).
- **Embed or reference (`lms-circulation-db`):** `Loan` **references** `studentId`/`bookId` by
  identifier; it does not embed a Student or Book copy. Both have their own lifecycle in a
  different bounded context (`ADR-006`), so nothing about them is embedded — consistent with
  `06-data/models.md`'s already-documented "no foreign key across services" design this ADR does
  not reopen.
- **`validationLevel` (`lms-circulation-db`):** starts at `strict` (`Anexo B` default) — every
  write is validated, not only new documents.
- **What can't be reverted:** nothing yet — all three repositories start from an empty schema, so
  every initial changeset has a real rollback (`DROP TABLE`/`dropCollection`). This line gets
  revisited the first time a `-db` repo ships a changeset that can't be cleanly reverted (a data
  migration with no inverse, for example), per `Anexo A`'s rule 12.

---

## Consequences

**Positive:**
- All three `-db` repos gain the reconstruction-verification CI shape the norm requires
  (`db-ci.yml`: build from empty, no-op on a second run, full rollback, rebuild) — something none
  of them has today
- `lms-circulation-api` loses its `EnsureIndexes()` responsibility, closing the exact gap a real
  PR review (`lms-circulation-api#2`) already flagged as critical
- One Liquibase habit, one CI shape, reused instead of relearned per repo

**Negative / Trade-offs:**
- Every existing `golang-migrate` migration in `lms-membership-db`/`lms-catalog-db` gets rewritten
  in Liquibase's changelog format — not a drop-in tool swap, a real one-time migration cost
- The team installs and maintains the Liquibase MongoDB extension (`liquibase-mongodb` +
  `mongodb` driver) alongside its Postgres usage — two `lpm add` targets to keep working, not one
- `text` + `CHECK` replaces every `VARCHAR(n)` column already written in `students`/`books` — a
  real schema change, not just a tooling one, required by `Anexo A` rule 6

**Impact on the system:**
- Affected repositories: `lms-membership-db`, `lms-catalog-db`, `lms-circulation-db` — their
  migration tooling, folder structure, and (for the two Postgres repos) schema/column types
- Affected indirectly: `lms-circulation-api` (removes `EnsureIndexes()` once
  `lms-circulation-db` owns index creation)
- Documents updated as part of this ADR: `06-data/migration-strategy.md` (tool and process, all
  three domains), `09-microservices/services/*/data-model.md` (schema name per domain, `text`
  columns), `05-architecture/overview.md` (any mention of `golang-migrate`)

---

## Risks

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|-----------|
| Rewriting already-shipped `students`/`books` migrations as new Liquibase changesets, instead of editing the old ones, temporarily doubles the visible schema history | Medium | Low — a one-time, well-understood cost | Old `golang-migrate` files are removed in the same PR that adds their Liquibase replacement, not left alongside it |
| Team unfamiliarity with the Liquibase MongoDB extension specifically (as opposed to Liquibase-for-Postgres, which is at least a known category of tool) | Medium | Low | `Anexo B`'s "Instalar la extensión" section already documents the exact two packages needed and the error message that shows up if only one is installed |
| `lms-circulation-api` still calling `EnsureIndexes()` after `lms-circulation-db` starts owning indexes creates a race or a duplicate-index error | Low | Low | Removing `EnsureIndexes()` from `lms-circulation-api` is tracked as a required follow-up in the same body of work, not left as a dangling gap |

---

## References

- Requires Liquibase or Flyway per `-db` repo → `rules/1-Norma-Repositorios-Sistemas-Distribuidos-2026B.pdf`, numeral 4.2.2
- PostgreSQL structure and rules this ADR commits `lms-membership-db`/`lms-catalog-db` to → `rules/A-db-postgres.md`
- MongoDB structure and rules this ADR commits `lms-circulation-db` to, including the "no schema management in the `-api`" rule → `rules/B-db-mongo.md`
- The prior review finding this ADR closes → `lms-circulation-api#2`, review comment "database structure lives in the `-db` repository, never in the `-api`"
- Related to: `ADR-005-mongodb-for-circulation-service.md` (chose MongoDB for Circulation; this ADR chooses how its schema is versioned)
- Related to: `ADR-006-repo-per-context-decomposition.md` (each `-db` repo is one domain's own schema, never a foreign key across them)
