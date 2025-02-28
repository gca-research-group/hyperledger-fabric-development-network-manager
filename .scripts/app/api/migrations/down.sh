#!/bin/bash
set -a
source .env
set +a 

export DATABASE_URL="postgres://$DATABASE_USER:$DATABASE_PASSWORD@$DATABASE_HOST:$DATABASE_PORT/$DATABASE_NAME?sslmode=disable"

migrate -database "$DATABASE_URL" -path "$MIGRATION_FOLDER" down
