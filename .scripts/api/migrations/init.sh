#!/bin/bash
set -a
source .env
set +a 

migrate create -ext sql -dir "$MIGRATION_FOLDER" -seq $1
