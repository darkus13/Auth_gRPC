#!/bin/bash
source .env

sleep 2 && goose -dir "${MIGRATIONS_DIR}" postgres "${MIGRATION_DSN}" up -v