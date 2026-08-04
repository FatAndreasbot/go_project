package project_init

import (
	"database/sql"

	postgresadapter "github.com/FatAndreasbot/go_project/user_service/adapters/outbound/persistance/postgres_adapter"
	"github.com/FatAndreasbot/go_project/user_service/infra/config"
)

func setupOutgoingAdapters() (*postgresadapter.UserPersistanceAdapter, *postgresadapter.GroupPersistanceAdapter, error) {
	connString := config.GetDBConnString()
	conn, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, nil, err
	}

	userAdapter := postgresadapter.NewUserPersistanceAdapter(conn)
	groupAdapter := postgresadapter.NewGroupPersistanceAdapter(conn)

	return userAdapter, groupAdapter, nil
}
