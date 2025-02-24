#!/bin/bash
set -a
source .env
set +a 

VERSION=latest

if [[ -n $1 ]]; then
    VERSION=$1
fi

export DATABASE_URL="postgres://$DATABASE_USER:$DATABASE_PASSWORD@$DATABASE_HOST:$DATABASE_PORT/$DATABASE_NAME?sslmode=disable"

migrate -path "$MIGRATION_FOLDER" -database "$DATABASE_URL" force "$VERSION"
