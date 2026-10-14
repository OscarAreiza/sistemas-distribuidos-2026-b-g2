# Release — LMS Library

_Generado: 2026-10-13_

## Resumen de versiones

- **Primer tag (y único publicado hasta ahora) en los 12 repos: `v2.0.1`.** Release real de GitHub publicado el 2026-10-13 (antes solo existía el tag crudo) — ver tabla de enlaces abajo.
- **Tag actual: sigue siendo `v2.0.1`** — todavía no se ha creado `v2.0.2`. 7 de los 12 repos (`lms-front`, `lms-catalog-api`, `lms-catalog-portal`, `lms-membership-portal`, `lms-circulation-api`, `lms-circulation-portal`, `lms-circulation-db`) ya tienen en `main` trabajo nuevo del 2026-10-13, sin etiquetar todavía. Los otros 5 (`lms-membership-api`, `lms-membership-db`, `lms-catalog-db`, `lms-infra-postgres`, `lms-infra-mongo`) no se tocaron en esta ronda — su `main` es exactamente el mismo commit que `v2.0.1`.

| Repositorio | Commit del tag `v2.0.1` | Commit actual en `main` |
|---|---|---|
| lms-front | `b012ea2a59db` | `d510df0054de` |
| lms-catalog-api | `5a2ccbbc77b1` | `05df113d9fce` |
| lms-catalog-portal | `b3b578dde029` | `aff85793c5a8` |
| lms-catalog-db | `8b0582795a7f` | `8b0582795a7f` |
| lms-membership-api | `4181abdc67df` | `4181abdc67df` |
| lms-membership-portal | `f62ed7fd153e` | `285796619f66` |
| lms-membership-db | `33a74cd04fd0` | `33a74cd04fd0` |
| lms-circulation-api | `54bbaa192f2e` | `4b6d64e2c76b` |
| lms-circulation-portal | `1659c628ceb2` | `b806df158a75` |
| lms-circulation-db | `0e0c932bf6ea` | `d29da28424a0` |
| lms-infra-postgres | `15fe583064a6` | `15fe583064a6` |
| lms-infra-mongo | `f1a359ecfe31` | `f1a359ecfe31` |

_(Los dos commits mostrados arriba son los commits reales — no el SHA del objeto tag anotado, que es un identificador distinto y no comparable directamente contra un commit.)_

---

## GitHub Release — enlaces reales (v2.0.1)

| Repositorio | Release |
|---|---|
| lms-front | https://github.com/code-corhuila/lms-front/releases/tag/v2.0.1 |
| lms-catalog-api | https://github.com/code-corhuila/lms-catalog-api/releases/tag/v2.0.1 |
| lms-catalog-portal | https://github.com/code-corhuila/lms-catalog-portal/releases/tag/v2.0.1 |
| lms-catalog-db | https://github.com/code-corhuila/lms-catalog-db/releases/tag/v2.0.1 |
| lms-membership-api | https://github.com/code-corhuila/lms-membership-api/releases/tag/v2.0.1 |
| lms-membership-portal | https://github.com/code-corhuila/lms-membership-portal/releases/tag/v2.0.1 |
| lms-membership-db | https://github.com/code-corhuila/lms-membership-db/releases/tag/v2.0.1 |
| lms-circulation-api | https://github.com/code-corhuila/lms-circulation-api/releases/tag/v2.0.1 |
| lms-circulation-portal | https://github.com/code-corhuila/lms-circulation-portal/releases/tag/v2.0.1 |
| lms-circulation-db | https://github.com/code-corhuila/lms-circulation-db/releases/tag/v2.0.1 |
| lms-infra-postgres | https://github.com/code-corhuila/lms-infra-postgres/releases/tag/v2.0.1 |
| lms-infra-mongo | https://github.com/code-corhuila/lms-infra-mongo/releases/tag/v2.0.1 |

---
## PRs por repositorio y por autor

Incluye únicamente PRs de **Oscar Areiza**, **Luis Meneses** (cuenta `MENESESLUIS-bot`) y **Hermes Pascuas** (cuenta `Hermesss29`) — se excluyen bots de dependencias u otros colaboradores si los hubiera.

### `lms-front`

**Oscar Areiza** — 8 PR(s)

| # | Título | Estado | Base | Fusionado | Commit de merge |
|---|---|---|---|---|---|
| #3 | chore(front): scaffold config, tooling and README | Merged | `develop` | 2026-10-08 | `b935bd9` |
| #4 | feat(front): shared http client, session, and remote-loading contract | Merged | `develop` | 2026-10-08 | `8248a3f` |
| #5 | feat(front): layout and shared UI components | Merged | `develop` | 2026-10-08 | `1cc358f` |
| #6 | feat(front): app shell routes, pages, and bootstrap | Merged | `develop` | 2026-10-08 | `ee903ff` |
| #7 | chore(front): add npm lockfile | Merged | `develop` | 2026-10-08 | `4a39789` |
| #16 | fix(mf): missing remoteEntry.js filename + UI depth pass | Merged | `develop` | 2026-10-12 | `448a12d` |
| #17 | promote(qa): mf filename fix + UI depth pass | Merged | `qa` | 2026-10-12 | `a5ec7e2` |
| #18 | release: v2.0.2 — mf filename fix + UI depth pass | Merged | `main` | 2026-10-12 | `d510df0` |

**Luis Meneses** — 3 PR(s)

| # | Título | Estado | Base | Fusionado | Commit de merge |
|---|---|---|---|---|---|
| #9 | test(front): add unit and component test suite with Vitest | Merged | `develop` | 2026-10-08 | `3e32e55` |
| #12 | test(front): add connection, stress and performance tests | Merged | `test/shell-unit-tests` | 2026-10-08 | `e680a4b` |
| #10 | qa: promote shell + test suite (unit, connection, stress, perf) to qa | Merged | `qa` | 2026-10-08 | `2ed1d0b` |

**Hermes Pascuas** — 3 PR(s)

| # | Título | Estado | Base | Fusionado | Commit de merge |
|---|---|---|---|---|---|
| #13 | test(front): add connection, stress and performance tests | Merged | `develop` | 2026-10-08 | `55be6fd` |
| #14 | chore(release): promote qa to release/2.0.1 — HU-01 shell | Merged | `release/2.0.1` | 2026-10-08 | `ccfa7ee` |
| #15 | chore(release): v2.0.1 — lms-front | Merged | `main` | 2026-10-08 | `b012ea2` |

---

### `lms-catalog-api`

**Oscar Areiza** — 7 PR(s)

| # | Título | Estado | Base | Fusionado | Commit de merge |
|---|---|---|---|---|---|
| #3 | fix(anexo-c): align with the course's -api repository norm | Merged | `develop` | 2026-09-26 | `ef11c88` |
| #4 | fix: RS256/HS256 auth, durable idempotency, GET /books, schema/UUID fixes | Merged | `develop` | 2026-10-03 | `92da7cd` |
| #9 | feat: feature flag for GET /books (week 9, config hardening) — v2 | Merged | `develop` | 2026-10-05 | `eef2ea5` |
| #12 | feat: PATCH/DELETE /books/{id} — edit (incl. copy count) and delete | Merged | `develop` | 2026-10-12 | `a68afc1` |
| #13 | promote(qa): PATCH/DELETE /books/{id} — edit + delete | Merged | `qa` | 2026-10-12 | `97ef5f6` |
| #14 | release: v2.0.2 — PATCH/DELETE /books/{id}, search feature flag | Merged | `main` | 2026-10-12 | `05df113` |
| #5 | feat: feature flag for GET /books (week 9, config hardening) | Closed | `develop` | — | `—` |

**Luis Meneses** — 6 PR(s)

| # | Título | Estado | Base | Fusionado | Commit de merge |
|---|---|---|---|---|---|
| #2 | chore(catalog): bootstrap catalog-service skeleton | Merged | `develop` | 2026-09-26 | `30172a7` |
| #8 | qa: promote catalog service skeleton from develop | Merged | `qa` | 2026-10-05 | `4dd992f` |
| #10 | release: 1.0.0 — promote validated qa stories to main | Merged | `main` | 2026-10-06 | `5a2ccbb` |
| #11 | fix(catalog): add PATCH /books/{id} for HU-09 book editing | Open | `develop` | — | `—` |
| #7 | chore(release): promote develop to qa | Closed | `qa` | — | `—` |
| #6 | chore(qa): promote catalog service skeleton and Annex C fixes from develop | Closed | `qa` | — | `—` |

**Hermes Pascuas** — 0 PR(s)

_(sin PRs de este autor en este repositorio)_

---

### `lms-catalog-portal`

**Oscar Areiza** — 3 PR(s)

| # | Título | Estado | Base | Fusionado | Commit de merge |
|---|---|---|---|---|---|
| #7 | fix: remoteEntry.js path/CORS, book delete, editable copies, UI polish | Merged | `develop` | 2026-10-12 | `8f7321f` |
| #8 | promote(qa): remoteEntry.js path/CORS, book delete, editable copies, UI polish | Merged | `qa` | 2026-10-12 | `b429009` |
| #9 | release: v2.0.2 — remoteEntry.js path/CORS, book delete, editable copies, UI polish | Merged | `main` | 2026-10-12 | `aff8579` |

**Luis Meneses** — 5 PR(s)

| # | Título | Estado | Base | Fusionado | Commit de merge |
|---|---|---|---|---|---|
| #2 | chore: scaffold catalog portal with book pages and deploy files | Merged | `develop` | 2026-09-29 | `5c973fa` |
| #3 | feat(catalog): add book pages, ui components and deploy files | Merged | `develop` | 2026-10-03 | `6a7797f` |
| #4 | chore(qa): promote project scaffold and catalog pages from develop | Merged | `qa` | 2026-10-06 | `8e36870` |
| #6 | release: 1.0.0 catalog portal mvp1 | Merged | `main` | 2026-10-06 | `b3b578d` |
| #5 | release: 0.1.0 catalog portal mvp1 | Closed | `main` | — | `—` |

**Hermes Pascuas** — 0 PR(s)

_(sin PRs de este autor en este repositorio)_

---

### `lms-catalog-db`

**Oscar Areiza** — 2 PR(s)

| # | Título | Estado | Base | Fusionado | Commit de merge |
|---|---|---|---|---|---|
| #3 | fix(anexo-a): align with the course's -db repository norm (Liquibase) | Merged | `develop` | 2026-10-03 | `4361bbe` |
| #4 | fix: migrate to the shared PostgreSQL instance (Anexo J) | Merged | `develop` | 2026-10-05 | `34e754a` |

**Luis Meneses** — 3 PR(s)

| # | Título | Estado | Base | Fusionado | Commit de merge |
|---|---|---|---|---|---|
| #2 | chore(migrations): migrate catalog-db schema | Merged | `develop` | 2026-09-24 | `03c45d4` |
| #5 | chore(qa): promote catalog-db Liquibase schema from develop | Merged | `qa` | 2026-10-06 | `f1f00c4` |
| #6 | release: 1.0.0 — catalog-db Liquibase schema | Merged | `main` | 2026-10-06 | `8b05827` |

**Hermes Pascuas** — 0 PR(s)

_(sin PRs de este autor en este repositorio)_

---

### `lms-membership-api`

**Oscar Areiza** — 3 PR(s)

| # | Título | Estado | Base | Fusionado | Commit de merge |
|---|---|---|---|---|---|
| #10 | fix(anexo-c): align with the course's -api repository norm | Merged | `develop` | 2026-09-26 | `552f7aa` |
| #11 | fix: RS256/HS256 auth, durable idempotency, schema/UUID fixes | Merged | `develop` | 2026-10-03 | `3697226` |
| #2 | chore: migrate membership-service entry point and tooling from lms-library | Closed | `develop` | — | `—` |

**Luis Meneses** — 0 PR(s)

_(sin PRs de este autor en este repositorio)_

**Hermes Pascuas** — 11 PR(s)

| # | Título | Estado | Base | Fusionado | Commit de merge |
|---|---|---|---|---|---|
| #3 | feat(membership): add student aggregate and domain ports | Merged | `develop` | 2026-09-26 | `3e1d854` |
| #4 | feat(membership): add create student use case | Merged | `develop` | 2026-09-26 | `3ebaafa` |
| #5 | feat(postgres): add connection pool and student repository | Merged | `develop` | 2026-09-26 | `81f4c30` |
| #6 | feat(http): add auth, correlation id and cors middleware | Merged | `develop` | 2026-09-26 | `b70fc95` |
| #7 | feat(membership): add update, deactivate and suspend use cases | Merged | `develop` | 2026-09-26 | `670534f` |
| #8 | feat(http): add student handlers | Merged | `develop` | 2026-09-26 | `26fa5e0` |
| #9 | feat(api): add service entry point | Merged | `develop` | 2026-09-26 | `8156ba0` |
| #12 | style(http): sort imports in student handler with gofmt | Merged | `develop` | 2026-10-04 | `5362f24` |
| #13 | chore(qa): promote HU-02 and HU-03 to qa | Merged | `qa` | 2026-10-07 | `a729d66` |
| #14 | chore(release): add HU-02 and HU-03 to release/2.0.1 | Merged | `release/2.0.1` | 2026-10-07 | `015730c` |
| #15 | release: v2.0.1 — membership (HU-02, HU-03) | Merged | `main` | 2026-10-07 | `4181abd` |

---

### `lms-membership-portal`

**Oscar Areiza** — 4 PR(s)

| # | Título | Estado | Base | Fusionado | Commit de merge |
|---|---|---|---|---|---|
| #5 | fix(anexo-h): align with the course's front-end repository norm | Merged | `develop` | 2026-10-03 | `1bb2527` |
| #10 | fix: remoteEntry.js path/CORS, apiClient types, UI polish | Merged | `develop` | 2026-10-12 | `062196e` |
| #11 | promote(qa): remoteEntry.js path/CORS, apiClient types, UI polish | Merged | `qa` | 2026-10-12 | `3c12435` |
| #12 | release: v2.0.2 — remoteEntry.js path/CORS, apiClient types, UI polish | Merged | `main` | 2026-10-12 | `2857966` |

**Luis Meneses** — 0 PR(s)

_(sin PRs de este autor en este repositorio)_

**Hermes Pascuas** — 7 PR(s)

| # | Título | Estado | Base | Fusionado | Commit de merge |
|---|---|---|---|---|---|
| #2 | chore(portal): add vite, react and typescript setup | Merged | `develop` | 2026-09-26 | `88714c9` |
| #3 | feat(students): add students list page | Merged | `develop` | 2026-09-26 | `378a93c` |
| #4 | feat(students): add student registration form page | Merged | `develop` | 2026-09-26 | `dd0752c` |
| #6 | feat(portal): wire routes and app entry point | Merged | `develop` | 2026-09-28 | `50a5420` |
| #7 | chore(qa): promote HU-02 and HU-03 to qa | Merged | `qa` | 2026-10-07 | `adde9bf` |
| #8 | chore(release): add HU-02 and HU-03 to release/2.0.1 | Merged | `release/2.0.1` | 2026-10-07 | `56d849f` |
| #9 | release: v2.0.1 — membership (HU-02, HU-03) | Merged | `main` | 2026-10-07 | `f62ed7f` |

---

### `lms-membership-db`

**Oscar Areiza** — 3 PR(s)

| # | Título | Estado | Base | Fusionado | Commit de merge |
|---|---|---|---|---|---|
| #3 | fix(anexo-a): align with the course's -db repository norm (Liquibase) | Merged | `develop` | 2026-09-26 | `679f485` |
| #4 | fix: add membership.idempotency_key table + supporting migrations | Merged | `develop` | 2026-10-03 | `af20708` |
| #5 | fix: migrate to the shared PostgreSQL instance (Anexo J) | Merged | `develop` | 2026-10-05 | `b7a389e` |

**Luis Meneses** — 0 PR(s)

_(sin PRs de este autor en este repositorio)_

**Hermes Pascuas** — 4 PR(s)

| # | Título | Estado | Base | Fusionado | Commit de merge |
|---|---|---|---|---|---|
| #2 | chore: add docker-compose for membership database | Merged | `develop` | 2026-09-24 | `c4202c2` |
| #6 | chore(qa): promote HU-02 and HU-03 to qa | Merged | `qa` | 2026-10-07 | `96ac312` |
| #7 | chore(release): add HU-02 and HU-03 to release/2.0.1 | Merged | `release/2.0.1` | 2026-10-07 | `906e0ce` |
| #8 | chore(release): v2.0.1 — membership (HU-02, HU-03) | Merged | `main` | 2026-10-07 | `33a74cd` |

---

### `lms-circulation-api`

**Oscar Areiza** — 10 PR(s)

| # | Título | Estado | Base | Fusionado | Commit de merge |
|---|---|---|---|---|---|
| #2 | chore: add circulation-service entry point from lms-library-sandbox | Merged | `develop` | 2026-09-26 | `f437030` |
| #3 | fix: RS256/HS256 auth, durable idempotency, 422 mapping, mongo-driver go.sum | Merged | `develop` | 2026-10-03 | `50d402e` |
| #6 | qa: promote Circulation API (HU-06/07/08, Mongo persistence) | Merged | `qa` | 2026-10-05 | `0602486` |
| #7 | test(stress): add stress-test script for /health, validated against qa | Merged | `develop` | 2026-10-05 | `394ad42` |
| #4 | fix: circulation_app's authSource moves from admin to loan_db (Anexo J) | Merged | `develop` | 2026-10-05 | `5586f3f` |
| #8 | release: 2.0.1 — Circulation API (HU-06, HU-07, HU-08) | Merged | `main` | 2026-10-05 | `54bbaa1` |
| #9 | feat(domain): lower LoanPeriodDays from 7 to 3 | Merged | `develop` | 2026-10-12 | `a95ce7e` |
| #10 | promote(qa): lower LoanPeriodDays from 7 to 3 | Merged | `qa` | 2026-10-12 | `e700bd2` |
| #11 | release: v2.0.2 — LoanPeriodDays = 3, authSource fix, stress test | Merged | `main` | 2026-10-12 | `4b6d64e` |
| #5 | Promote develop to qa (Circulation API) | Closed | `qa` | — | `—` |

**Luis Meneses** — 0 PR(s)

_(sin PRs de este autor en este repositorio)_

**Hermes Pascuas** — 0 PR(s)

_(sin PRs de este autor en este repositorio)_

---

### `lms-circulation-portal`

**Oscar Areiza** — 7 PR(s)

| # | Título | Estado | Base | Fusionado | Commit de merge |
|---|---|---|---|---|---|
| #7 | qa: promote HU-06/07/08 (loan registration, list+return, overdue report) | Merged | `qa` | 2026-10-05 | `8fbd7e7` |
| #8 | release: 2.0.1 — Circulation Portal (HU-06, HU-07, HU-08) | Merged | `main` | 2026-10-05 | `1659c62` |
| #9 | fix: SearchableSelect dropdown overlap under Module Federation, loan list readability | Merged | `develop` | 2026-10-12 | `e1d0056` |
| #10 | promote(qa): SearchableSelect dropdown overlap fix, loan list readability | Merged | `qa` | 2026-10-12 | `2afbd89` |
| #11 | release: v2.0.2 — SearchableSelect overlap fix, loan list readability | Merged | `main` | 2026-10-12 | `b806df1` |
| #6 | Promote develop to qa (Circulation Portal) | Closed | `qa` | — | `—` |
| #5 | feat: register loan form (HU-06), type-ahead student/book search | Closed | `develop` | — | `—` |

**Luis Meneses** — 0 PR(s)

_(sin PRs de este autor en este repositorio)_

**Hermes Pascuas** — 3 PR(s)

| # | Título | Estado | Base | Fusionado | Commit de merge |
|---|---|---|---|---|---|
| #2 | chore(portal): scaffold circulation portal as module federation remote | Merged | `develop` | 2026-10-01 | `9c2a6e8` |
| #3 | feat(loans): add loans list with return action and overdue report (HU-07, HU-08) | Merged | `develop` | 2026-10-03 | `cd0cff2` |
| #4 | feat(loans): add loan registration form with idempotency key (HU-06) | Merged | `develop` | 2026-10-04 | `84f4ba4` |

---

### `lms-circulation-db`

**Oscar Areiza** — 9 PR(s)

| # | Título | Estado | Base | Fusionado | Commit de merge |
|---|---|---|---|---|---|
| #2 | chore: add MongoDB compose service | Merged | `develop` | 2026-09-26 | `c0c63b1` |
| #3 | fix: add idempotency_keys collection + role grants | Merged | `develop` | 2026-10-03 | `d024870` |
| #4 | fix: migrate to the shared MongoDB instance (Anexo J) | Merged | `develop` | 2026-10-04 | `a7da442` |
| #6 | qa: promote Circulation DB Mongo migrations (Anexo J + idempotency_keys) | Merged | `qa` | 2026-10-05 | `a0be65c` |
| #7 | release: 2.0.1 — Circulation DB (Anexo J + idempotency_keys) | Merged | `main` | 2026-10-05 | `0e0c932` |
| #8 | fix(dcl): grant circulation_writer update on idempotency_keys | Merged | `develop` | 2026-10-12 | `5232e08` |
| #9 | promote(qa): grant circulation_writer update on idempotency_keys | Merged | `qa` | 2026-10-12 | `6c6aabf` |
| #10 | release: v2.0.2 — idempotency_keys update privilege, Anexo J backfill | Merged | `main` | 2026-10-12 | `d29da28` |
| #5 | Promote develop to qa (Circulation DB — Mongo migrations) | Closed | `qa` | — | `—` |

**Luis Meneses** — 0 PR(s)

_(sin PRs de este autor en este repositorio)_

**Hermes Pascuas** — 0 PR(s)

_(sin PRs de este autor en este repositorio)_

---

### `lms-infra-postgres`

**Oscar Areiza** — 3 PR(s)

| # | Título | Estado | Base | Fusionado | Commit de merge |
|---|---|---|---|---|---|
| #2 | feat: single PostgreSQL instance for the whole project (Anexo J) | Merged | `develop` | 2026-10-05 | `4466b17` |
| #3 | qa: promote single PostgreSQL instance (Anexo J) | Merged | `qa` | 2026-10-07 | `3e19856` |
| #4 | release: 2.0.1 — single PostgreSQL instance (Anexo J) | Merged | `main` | 2026-10-07 | `15fe583` |

**Luis Meneses** — 0 PR(s)

_(sin PRs de este autor en este repositorio)_

**Hermes Pascuas** — 0 PR(s)

_(sin PRs de este autor en este repositorio)_

---

### `lms-infra-mongo`

**Oscar Areiza** — 3 PR(s)

| # | Título | Estado | Base | Fusionado | Commit de merge |
|---|---|---|---|---|---|
| #1 | feat(infra): shared MongoDB instance for Circulation (Anexo J) | Merged | `develop` | 2026-10-06 | `fbda419` |
| #2 | qa: promote shared MongoDB instance for Circulation (Anexo J) | Merged | `qa` | 2026-10-07 | `e1debec` |
| #3 | release: 2.0.1 — shared MongoDB instance for Circulation (Anexo J) | Merged | `main` | 2026-10-07 | `f1a359e` |

**Luis Meneses** — 0 PR(s)

_(sin PRs de este autor en este repositorio)_

**Hermes Pascuas** — 0 PR(s)

_(sin PRs de este autor en este repositorio)_

---

## Totales generales por autor (12 repositorios)

| Autor | Total de PRs |
|---|---|
| Oscar Areiza | 62 |
| Luis Meneses | 17 |
| Hermes Pascuas | 28 |
