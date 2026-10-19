package postgresrepository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/infra/postgresql/sqlc_gen"
	"github.com/google/uuid"
)

func PermissionDifference(old, new []*models.Permission) (toAdd, toRemove []*models.Permission) {
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
			toAdd = append(toAdd, permission)
		}
	}
	return
}

type GroupPersistanceAdapter struct {
	conn *sql.DB
	q    *sqlc_gen.Queries
}

func NewGroupPersistanceAdapter(conn *sql.DB) *GroupPersistanceAdapter {
	return &GroupPersistanceAdapter{
		conn: conn,
		q:    sqlc_gen.New(conn),
	}
}

func (adp *GroupPersistanceAdapter) CreateGroup(ctx context.Context, groupName string, permissions []*models.Permission, idempotencyKey uuid.UUID) (*models.Group, error) {
	tx, err := adp.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	groupData, err := adp.q.WithTx(tx).CreateGroup(ctx, groupName)
	if err != nil {
		return nil, errors.Join(err, tx.Rollback())
	}

	newGroup := models.Group{
		MainGroupData: models.MainGroupData{
			ID:        groupData.ID,
			Name:      groupName,
			CreatedAt: groupData.CreatedAt,
		},
		Permissions: make([]*models.Permission, 0),
	}

	return &newGroup, tx.Commit()
}

func (adp *GroupPersistanceAdapter) GetGroupByID(ctx context.Context, groupID uuid.UUID) (*models.MainGroupData, error) {
	groupData, err := adp.q.GetGroupDataByID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return &models.MainGroupData{
		ID:        groupData.ID,
		Name:      groupData.Name,
		CreatedAt: groupData.CreatedAt,
	}, nil
}

func (adp *GroupPersistanceAdapter) GetGroups(ctx context.Context, cursor time.Time, count int) ([]*models.MainGroupData, error) {
	groupData, err := adp.q.GetGroupList(ctx, sqlc_gen.GetGroupListParams{
		CreatedAt: cursor,
		Limit:     int32(count),
	})
	if err != nil {
		return nil, err
	}

	groups := make([]*models.MainGroupData, 0, len(groupData))

	for _, group := range groupData {
		groups = append(groups, &models.MainGroupData{
			ID:        group.ID,
			Name:      group.Name,
			CreatedAt: group.CreatedAt,
		})
	}

	return groups, nil
}

func (adp *GroupPersistanceAdapter) UpdateGroup(ctx context.Context, groupData *models.MainGroupData) error {
	tx, err := adp.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = adp.q.WithTx(tx).UpdateGroupName(ctx, sqlc_gen.UpdateGroupNameParams{
		Name: groupData.Name,
		ID:   groupData.ID,
	})
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}

	return tx.Commit()
}

func (adp *GroupPersistanceAdapter) DeleteGroup(ctx context.Context, groupID uuid.UUID) error {
	tx, err := adp.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = adp.q.WithTx(tx).DeleteGroup(ctx, groupID)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}

	return tx.Commit()
}

func (adp *GroupPersistanceAdapter) GetGroupPermissions(ctx context.Context, groupID uuid.UUID) ([]*models.Permission, error) {
	permissionsData, err := adp.q.GetGroupPermissions(ctx, groupID)
	if err != nil {
		return nil, err
	}

	permissions := make([]*models.Permission, 0, len(permissionsData))
	for _, permission := range permissionsData {
		permissions = append(permissions, &models.Permission{
			ID:   permission.ID,
			Name: models.PermissionName(permission.Name),
		})
	}

	return permissions, nil
}

func (adp *GroupPersistanceAdapter) AddPermission(ctx context.Context, groupID, permissionID uuid.UUID) error {
	tx, err := adp.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = adp.q.WithTx(tx).AddPermissionToGroup(ctx, sqlc_gen.AddPermissionToGroupParams{
		GroupID:      groupID,
		PermissionID: permissionID,
	})
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}

	return tx.Commit()
}

func (adp *GroupPersistanceAdapter) RemovePermission(ctx context.Context, groupID, permissionID uuid.UUID) error {
	tx, err := adp.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = adp.q.WithTx(tx).RemovePermissionFromGroup(ctx, sqlc_gen.RemovePermissionFromGroupParams{
		GroupID:      groupID,
		PermissionID: permissionID,
	})
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}

	return tx.Commit()
}

func (adp *GroupPersistanceAdapter) CheckIdempotencyKey(ctx context.Context, key uuid.UUID) error {
	tx, err := adp.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	wasCreated, err := adp.q.WithTx(tx).CheckUserWasCreated(ctx, key)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if wasCreated {
		return errors.New("grpup was created")
	}
	return nil
}
