migrate-up:
	docker compose exec api goose \
		-dir /app/db/migrations \
		mysql "app:password@tcp(mysql:3306)/portfolio?parseTime=true" up

migrate-down:
	docker compose exec api goose \
		-dir /app/db/migrations \
		mysql "app:password@tcp(mysql:3306)/portfolio?parseTime=true" down
