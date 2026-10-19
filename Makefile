migrate-up:
	docker compose exec api sh -c 'goose -dir /app/db/migrations mysql "$$DB_USER:$$DB_PASSWORD@tcp($$DB_HOST:$$DB_PORT)/$$DB_NAME?parseTime=true" up'

migrate-down:
	docker compose exec api sh -c 'goose -dir /app/db/migrations mysql "$$DB_USER:$$DB_PASSWORD@tcp($$DB_HOST:$$DB_PORT)/$$DB_NAME?parseTime=true" down'
