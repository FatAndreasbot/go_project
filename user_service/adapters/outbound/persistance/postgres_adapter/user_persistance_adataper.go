package postgresadapter

import (
	"context"
	"encoding/json"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/infra/persistance/postgresql/sqlc_gen"
	"github.com/google/uuid"
)

func permissionDifference(old, new []*models.Permission) (toAdd, toRemove []*models.Permission) {
	toSet := func(list []*models.Permission) map[uuid.UUID]*models.Permission {
		set := make(map[uuid.UUID]*models.Permission, len(list))
		for _, permission := range list {
			set[permission.ID] = permission
		}
		return set
	}

	oldSet := toSet(old)
	newSet := toSet(new)

	for id, permission := range oldSet {
		_, ok := newSet[id]
		if !ok {
			toRemove = append(toRemove, permission)
		}
	}

	for id, permission := range newSet {
		_, ok := oldSet[id]
		if !ok {
			toAdd = append(toRemove, permission)
		}
	}
	return
}

// implements UserPersistancePort
type UserPersistanceAdapter struct {
	q *sqlc_gen.Queries
}

// UserPersistancePort
// GetUserByUsername(context.Context, string) (*models.User, error)
// GetUserByID(context.Context, uuid.UUID) (*models.User, error)

// GetUserList(context.Context, int, int) ([]*models.User, error)

// UpdateUser(context.Context, uuid.UUID, *models.User) error
// CreateUser(context.Context, *models.User) (uuid.UUID, error)
// DeleteUser(context.Context, uuid.UUID) error

func (adapter *UserPersistanceAdapter) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	row, err := adapter.q.GetUserByUsername(ctx, username)
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

func (adapter *UserPersistanceAdapter) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	row, err := adapter.q.GetUserByID(ctx, id)
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

func (adapter *UserPersistanceAdapter) GetUserList(ctx context.Context, limit, offset int) ([]*models.User, error) {
	rows, err := adapter.q.GetUserList(ctx, sqlc_gen.GetUserListParams{
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

func (adapter *UserPersistanceAdapter) UpdateUser(ctx context.Context, userID uuid.UUID, newUserData *models.User) error {
	return adapter.q.UpdateUser(ctx, sqlc_gen.UpdateUserParams{
		Username:     newUserData.Name,
		PasswordHash: newUserData.PasswordHash,
		GroupID:      newUserData.Group.ID,
	})
}

func (adapter *UserPersistanceAdapter) CreateUser(ctx context.Context, user *models.User) error {
	userID, err := adapter.q.CreateUser(ctx, sqlc_gen.CreateUserParams{
		Username:     user.Name,
		PasswordHash: user.PasswordHash,
		GroupID:      user.Group.ID,
	})
	if err != nil {
		return err
	}

	user.ID = userID
	return nil
}

func (adapter *UserPersistanceAdapter) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return adapter.q.DeleteUser(ctx, id)
}
