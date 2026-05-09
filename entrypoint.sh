#!/bin/sh
set -e

echo "Running migrations..."
case "${DATABASE_DRIVER:-postgres}" in
  postgres|postgresql)
    GOOSE_DRIVER=postgres GOOSE_DBSTRING="$DATABASE_URL" GOOSE_MIGRATION_DIR=migrations goose up
    ;;
  sqlite)
    GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING="$DATABASE_URL" GOOSE_MIGRATION_DIR=sqlite_migrations goose up
    ;;
  *)
    echo "Unsupported DATABASE_DRIVER: $DATABASE_DRIVER"
    exit 1
    ;;
esac

echo "Starting app..."
exec ./server
