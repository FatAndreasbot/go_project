package postgresadapter

import (
	"database/sql"

	"github.com/FatAndreasbot/go_project/user_service/infra/postgresql/sqlc_gen"
)

// implements UserPersistancePort
type UserPersistanceAdapter struct {
	conn *sql.DB
	q    *sqlc_gen.Queries
}

func NewUserPersistanceAdapter(conn *sql.DB) *UserPersistanceAdapter {
	return &UserPersistanceAdapter{
		conn: conn,
		q:    sqlc_gen.New(conn),
	}
}
