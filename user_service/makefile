
.PHONY: migrations

migrations:
	goose -dir \
	./infra/postgresql/migrations postgres \
	"postgres://postgres:postgres@localhost:5432/user-service?sslmode=disable" up

migrations_down:
	goose -dir \
	./infra/postgresql/migrations postgres \
	"postgres://postgres:postgres@localhost:5432/user-service?sslmode=disable" down
