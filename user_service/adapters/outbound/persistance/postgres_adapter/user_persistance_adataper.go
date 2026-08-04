package postgresadapter

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/infra/persistance/postgresql/sqlc_gen"
	"github.com/google/uuid"
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

// UserPersistancePort
// GetUserByUsername(context.Context, string) (*models.User, error)
// GetUserByID(context.Context, uuid.UUID) (*models.User, error)

// GetUserList(context.Context, int, int) ([]*models.User, error)

// UpdateUser(context.Context, uuid.UUID, *models.User) error
// CreateUser(context.Context, *models.User) (uuid.UUID, error)
// DeleteUser(context.Context, uuid.UUID) error

func (adp *UserPersistanceAdapter) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	row, err := adp.q.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	permissions := make([]*models.Permission, 0)

	err = json.Unmarshal(row.Permissions, &permissions)
	if err != nil {
		return nil, err
	}

	user := models.User{
		ID:           row.Uuid,
		Name:         row.Username,
		PasswordHash: row.PasswordHash,
		Group: &models.Group{
			ID:          row.GroupUuid,
			Name:        row.GroupName,
			Permissions: permissions,
		},
	}

	return &user, nil
}

func (adp *UserPersistanceAdapter) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	row, err := adp.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	permissions := make([]*models.Permission, 0)
	err = json.Unmarshal(row.Permissions, &permissions)
	if err != nil {
		return nil, err
	}

	user := models.User{
		ID:           row.Uuid,
		Name:         row.Username,
		PasswordHash: row.PasswordHash,
		Group: &models.Group{
			ID:          row.GroupUuid,
			Name:        row.GroupName,
			Permissions: permissions,
		},
	}

	return &user, nil
}

func (adp *UserPersistanceAdapter) GetUserList(ctx context.Context, limit, offset int) ([]*models.User, error) {
	rows, err := adp.q.GetUserList(ctx, sqlc_gen.GetUserListParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	users := make([]*models.User, 0, len(rows))

	for _, row := range rows {
		var permissions []*models.Permission
		err = json.Unmarshal(row.Permissions, &permissions)

		if err != nil {
			return nil, err
		}

		users = append(users, &models.User{
			ID:           row.Uuid,
			Name:         row.Username,
			PasswordHash: row.PasswordHash,
			Group: &models.Group{
				ID:          row.GroupUuid,
				Name:        row.GroupName,
				Permissions: permissions,
			},
		})
	}

	return users, nil
}

func (adp *UserPersistanceAdapter) UpdateUser(ctx context.Context, userID uuid.UUID, newUserData *models.User) error {
	return adp.q.UpdateUser(ctx, sqlc_gen.UpdateUserParams{
		Username:     newUserData.Name,
		PasswordHash: newUserData.PasswordHash,
		GroupID:      newUserData.Group.ID,
	})
}

func (adp *UserPersistanceAdapter) CreateUser(ctx context.Context, user *models.User) (uuid.UUID, error) {
	userID, err := adp.q.CreateUser(ctx, sqlc_gen.CreateUserParams{
		Username:     user.Name,
		PasswordHash: user.PasswordHash,
		GroupID:      user.Group.ID,
	})
	if err != nil {
		return uuid.Nil, err
	}

	user.ID = userID
	return userID, nil
}

func (adp *UserPersistanceAdapter) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return adp.q.DeleteUser(ctx, id)
}
