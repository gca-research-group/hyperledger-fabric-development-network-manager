POSTGRES_DATA_PATH=./.docker/api/volumes/postgres/data

if [[ ! -d "$POSTGRES_DATA_PATH" ]]; then
  mkdir -p "$POSTGRES_DATA_PATH"
fi

docker compose -f ./.docker/api/docker-compose.yml up hfndm_database --build -d

./.scripts/app/api/migrations/up.sh

go mod tidy

air
