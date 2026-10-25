package project_init

import (
	"database/sql"
	"log"
	"time"

	postgresrepository "github.com/FatAndreasbot/go_project/user_service/adapters/outbound/postgres_repository"
	sessionstorage "github.com/FatAndreasbot/go_project/user_service/adapters/outbound/session_storage"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/ports/outgoing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func SetupOutgoingAdapters(config *models.ConfigOptions) (
	outgoing.UserRepository,
	outgoing.GroupRepository,
	outgoing.PermissionRepository,

	outgoing.SessionHandler,

	error,
) {
	connString := config.GetDBConnString()

	conn, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if err := conn.Ping(); err != nil {
		log.Default().Println(err.Error())
		return nil, nil, nil, nil, err
	}

	conn.SetMaxOpenConns(20)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxIdleTime(5 * time.Minute)
	conn.SetConnMaxLifetime(30 * time.Minute)

	upa := postgresrepository.NewUserPersistanceAdapter(conn)
	gpa := postgresrepository.NewGroupPersistanceAdapter(conn)
	ppa := postgresrepository.NewPermissionPersistanceAdapter(conn)

	sessionHandler := sessionstorage.NewMapSessionStorate(config)

	return upa, gpa, ppa, sessionHandler, nil
}
