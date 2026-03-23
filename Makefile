include .env
export

migrate-up:
	migrate -path migrations/ -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(MIGRATE_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" up

migrate-down:
	migrate -path migrations/ -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(MIGRATE_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" down

test:
	go test ./internal/service/... ./internal/handlers/...

coverage:
	go test ./internal/service/... ./internal/handlers/... -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out

coverage-html: coverage
	go tool cover -html=coverage.out -o coverage.html
	xdg-open coverage.html

clean-coverage:
	rm -f coverage.out coverage.html