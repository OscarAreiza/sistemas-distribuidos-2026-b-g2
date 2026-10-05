#!/bin/sh
# Runs once, only when the postgres-data volume is empty (official postgres
# image behavior — rules/3-Anexo-J, J.5.2 "La base se crea una sola vez").
# Idempotent (J.5.5): a new domain's user is appended here and the script is
# re-run by hand in an environment that already exists — every statement
# below tolerates being re-applied without error.
set -e

psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB" <<SQL
CREATE EXTENSION IF NOT EXISTS pgcrypto;

SELECT format('CREATE ROLE %I LOGIN PASSWORD %L', 'membership_app', '$MEMBERSHIP_APP_PASSWORD')
WHERE NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'membership_app')\gexec
ALTER ROLE membership_app SET search_path = membership;

SELECT format('CREATE ROLE %I LOGIN PASSWORD %L', 'catalog_app', '$CATALOG_APP_PASSWORD')
WHERE NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'catalog_app')\gexec
ALTER ROLE catalog_app SET search_path = catalog;
SQL
