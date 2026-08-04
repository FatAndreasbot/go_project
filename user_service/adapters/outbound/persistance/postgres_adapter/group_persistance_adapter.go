package postgresadapter

import (
	"context"
	"database/sql"
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

type GroupPersistanceAdapter struct {
	q    *sqlc_gen.Queries
	conn *sql.DB
}

func NewGroupPersistanceAdapter(conn *sql.DB) *GroupPersistanceAdapter {
	return &GroupPersistanceAdapter{
		conn: conn,
		q:    sqlc_gen.New(conn),
	}
}

// GroupPersistancePort interface
// GetGroupByID(context.Context, uuid.UUID) (*models.Group, error)

// GetGroupList(ctx context.Context, limit, offset int) ([]*models.Group, error)

// UpdateGroup(context.Context, uuid.UUID, *models.Group) error
// CreateGroup(context.Context, *models.Group) (uuid.UUID, error)
// DeleteGroup(context.Context, uuid.UUID) error

func (adp *GroupPersistanceAdapter) GetGroupByID(ctx context.Context, id uuid.UUID) (*models.Group, error) {
	row, err := adp.q.GetGroupByID(ctx, id)
	if err != nil {
		return nil, err
	}

	permissions := make([]*models.Permission, 0)

	err = json.Unmarshal(row.Permissions, &permissions)
	if err != nil {
		return nil, err
	}

	return &models.Group{
		ID:          row.ID,
		Name:        row.Name,
		Permissions: permissions,
	}, nil
}

func (adp *GroupPersistanceAdapter) GetGroupList(ctx context.Context, limit, offset int) ([]*models.Group, error) {
	rows, err := adp.q.GetGroupList(ctx, sqlc_gen.GetGroupListParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	groups := make([]*models.Group, 0, len(rows))

	for _, row := range rows {
		var permissions []*models.Permission
		err = json.Unmarshal(row.Permissions, &permissions)

		if err != nil {
			return nil, err
		}

		groups = append(groups, &models.Group{
			ID:          row.ID,
			Name:        row.Name,
			Permissions: permissions,
		})
	}

	return groups, nil
}

func (adp *GroupPersistanceAdapter) UpdateGroup(ctx context.Context, oldGroupID uuid.UUID, newGroupData *models.Group) (err error) {
	tx, err := adp.conn.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = adp.q.UpdateGroup(ctx, newGroupData.Name)
	if err != nil {
		return
	}

	oldGroupData, err := adp.GetGroupByID(ctx, oldGroupID)
	if err != nil {
		return
	}

	permissionsToAdd, peroissionsToRemove := permissionDifference(oldGroupData.Permissions, newGroupData.Permissions)

	for _, permission := range permissionsToAdd {
		err = adp.q.AddPermissionsToGroup(ctx, sqlc_gen.AddPermissionsToGroupParams{
			GroupID:      oldGroupID,
			PermissionID: permission.ID,
		})
		if err != nil {
			return
		}
	}

	for _, permission := range peroissionsToRemove {
		adp.q.RemovePermissionsFromGroup(ctx, sqlc_gen.RemovePermissionsFromGroupParams{
			GroupID:      oldGroupID,
			PermissionID: permission.ID,
		})
		if err != nil {
			return
		}
	}

	err = tx.Commit()
	return
}

func (adp *GroupPersistanceAdapter) CreateGroup(ctx context.Context, groupData *models.Group) (newGroupID uuid.UUID, err error) {
	tx, err := adp.conn.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	newGroupID, err = adp.q.CreateGroup(ctx, groupData.Name)

	for _, permission := range groupData.Permissions {
		err = adp.q.AddPermissionsToGroup(ctx, sqlc_gen.AddPermissionsToGroupParams{
			GroupID:      newGroupID,
			PermissionID: permission.ID,
		})
		if err != nil {
			return
		}
	}

	err = tx.Commit()

	return
}

func (adp *GroupPersistanceAdapter) DeleteGroup(ctx context.Context, groupID uuid.UUID) error {
	return adp.q.DeleteGroup(ctx, groupID)
}
