package project_init

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	postgresadapter "github.com/FatAndreasbot/go_project/user_service/adapters/outbound/postgres_adapter"
	"github.com/FatAndreasbot/go_project/user_service/infra/config"
)

func setupOutgoingAdapters() (*postgresadapter.UserPersistanceAdapter, *postgresadapter.GroupPersistanceAdapter, error) {
	connString := config.GetConfig().DBConnString
	conn, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, nil, err
	}

	err = conn.Ping()
	if err != nil {
		log.Default().Println(err)
		return nil, nil, errors.New("could not establish connection to database")
	}

	conn.SetMaxOpenConns(20)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxLifetime(30 * time.Minute)
	conn.SetConnMaxIdleTime(5 * time.Minute)

	userAdapter := postgresadapter.NewUserPersistanceAdapter(conn)
	groupAdapter := postgresadapter.NewGroupPersistanceAdapter(conn)

	return userAdapter, groupAdapter, nil
}
