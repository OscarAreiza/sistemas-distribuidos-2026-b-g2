# ADR-C3-02: Cross-AZ active-standby for the LMS data stores with automated backups to OBS

## Status

Proposed

## Context

All LMS data lives in two shared instances: one PostgreSQL 16 instance with one schema per domain (catalog, membership, workflow `saga_runs`) and one MongoDB instance holding Circulation's `loan_db` (`10-week/hu-status/PERSISTENCE-SAGA-OUTBOX.md`, section 1). Both run as containers on the same Docker host (`lms-infra-postgres/compose.yml`, `lms-infra-mongo/compose.yml`), with no replica and no scheduled backup, so losing that host or its disk loses every loan, member and penalty record. Loans are the critical data: a lost loan means a book that the library cannot trace. The design note fixes the critical data store, the RTO and the RPO in section 4 (Disaster recovery), and lists the missing replica as the reliability weakness in section 2 (Five principles).

## Decision

We will move the PostgreSQL schemas to RDS for PostgreSQL in primary/standby mode with the standby in a second availability zone, and Circulation's `loan_db` to DDS (Document Database Service) as a three-node replica set spread across AZs. Both services fail over automatically on an AZ or instance failure, which meets the RTO in section 4. Automated full backups run daily and incremental/log backups run every 5 minutes, stored in OBS with a 7-day retention, which meets the RPO and lets us restore to a point in time after a logical error (for example, a bad migration). The per-domain schemas, least-privilege logins and Liquibase migrations (`ADR-010`) stay as they are; only the host changes.

## Alternatives considered

- Keep self-managed containers and add a nightly `pg_dump` / `mongodump` to OBS: it is cheap, but it is backup-and-restore only, so the RPO would be up to 24 hours and the RTO the time to rebuild a host and restore by hand, which is far above the numbers in section 4.
- Active-active across two regions: it gives the lowest RTO and RPO, but it needs multi-region write conflict handling for loans and doubles the cost, which is not justified for one university library with a single campus.

## Consequences

- Positive: an AZ failure no longer loses data or stops loans, failover needs no manual action, and point-in-time recovery protects us against human error as well as infrastructure failure.
- Negative: managed RDS and DDS instances cost more than containers on the existing host, the services must handle a short reconnect during failover (connection retries in the drivers), and the team loses direct root access to the database engines.
