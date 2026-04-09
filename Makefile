include .env
export

docker-up:
	@docker compose up -d db

docker-down:
	@docker compose down

migrate-create:
	@test -n "$(seq)" || (echo "Missing seq parameter"; exit 1)
	docker compose run --rm db-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-action:
	@test -n "$(action)" || (echo "Missing action parameter"; exit 1)
	docker compose run --rm db-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@db:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

run-server:
	go run ./cmd/api/main.go 