
.PHONY: migrations build

migrations:
	goose -dir \
	./infra/postgresql/migrations postgres \
	"postgres://postgres:postgres@localhost:5432/user-service?sslmode=disable" up

migrations_down:
	goose -dir \
	./infra/postgresql/migrations postgres \
	"postgres://postgres:postgres@localhost:5432/user-service?sslmode=disable" down

generate_queries:
	sqlc generate

clear_queries:
	rm -rf infra/postgresql/sqlc_gen/*

build:
	mkdir -p build
	go build -o build/user_service cmd/main.go