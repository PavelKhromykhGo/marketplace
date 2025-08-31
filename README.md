если база на хосте:
docker run --rm -p 8080:8080 \
-e DATABASE_URL="host=host.docker.internal port=5432 user=postgres password=postgres dbname=marketplace sslmode=disable" \
-e MIGRATIONS_DIR="/app/migrations" \
marketplace

если база в контейнере:
docker network create mpnet
docker run -d --name postgres --network mpnet -e POSTGRES_DB=marketplace -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:16
docker run --rm --name app --network mpnet -p 8080:8080 \
-e DATABASE_URL="host=postgres port=5432 user=postgres password=postgres dbname=marketplace sslmode=disable" \
-e MIGRATIONS_DIR="/app/migrations" \
marketplace