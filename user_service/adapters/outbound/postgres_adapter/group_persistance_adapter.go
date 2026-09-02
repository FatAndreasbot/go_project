package postgresadapter

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

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

const transactionKey = "sql_tx-3abaaa8c-ac90-46ed-8d9e-2e750c95b4cd"

func (adp *GroupPersistanceAdapter) GetGroupByID(ctx context.Context, id uuid.UUID) (group *models.Group, err error) {
	ctxTx := ctx.Value(transactionKey)
	var tx *sql.Tx
	if ctxTx == nil {
		tx, err = adp.conn.Begin()
		if err != nil {
			return
		}
		ctx = context.WithValue(ctx, transactionKey, tx)
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()
	} else {
		var ok bool
		tx, ok = ctxTx.(*sql.Tx)
		if !ok {
			err = errors.New("could not cast into *sql.Tx")
		}
	}

	row, err := adp.q.WithTx(tx).GetGroupByID(ctx, id)
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

func (adp *GroupPersistanceAdapter) GetGroupList(ctx context.Context, limit int, name_cursor string) (groups []*models.Group, err error) {
	ctxTx := ctx.Value(transactionKey)
	var tx *sql.Tx
	if ctxTx == nil {
		tx, err = adp.conn.Begin()
		if err != nil {
			return
		}
		ctx = context.WithValue(ctx, transactionKey, tx)
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()
	} else {
		var ok bool
		tx, ok = ctxTx.(*sql.Tx)
		if !ok {
			err = errors.New("could not cast into *sql.Tx")
		}
	}

	rows, err := adp.q.WithTx(tx).GetGroupList(ctx, sqlc_gen.GetGroupListParams{
		Limit: int32(limit),
		Name:  name_cursor,
	})
	if err != nil {
		return nil, err
	}

	groups = make([]*models.Group, 0, len(rows))

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
	ctxTx := ctx.Value(transactionKey)
	var tx *sql.Tx
	if ctxTx == nil {
		tx, err = adp.conn.Begin()
		if err != nil {
			return
		}
		ctx = context.WithValue(ctx, transactionKey, tx)
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()
	} else {
		var ok bool
		tx, ok = ctxTx.(*sql.Tx)
		if !ok {
			err = errors.New("could not cast into *sql.Tx")
		}
	}

	err = adp.q.WithTx(tx).UpdateGroup(ctx, sqlc_gen.UpdateGroupParams{
		Name: newGroupData.Name,
		ID:   oldGroupID,
	})
	if err != nil {
		return
	}

	oldGroupData, err := adp.GetGroupByID(ctx, oldGroupID)
	if err != nil {
		return
	}

	permissionsToAdd, peroissionsToRemove := PermissionDifference(oldGroupData.Permissions, newGroupData.Permissions)

	for _, permission := range permissionsToAdd {
		err = adp.q.WithTx(tx).AddPermissionsToGroup(ctx, sqlc_gen.AddPermissionsToGroupParams{
			GroupID:      oldGroupID,
			PermissionID: permission.ID,
		})
		if err != nil {
			return
		}
	}

	for _, permission := range peroissionsToRemove {
		adp.q.WithTx(tx).RemovePermissionsFromGroup(ctx, sqlc_gen.RemovePermissionsFromGroupParams{
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
	ctxTx := ctx.Value(transactionKey)
	var tx *sql.Tx
	if ctxTx == nil {
		tx, err = adp.conn.Begin()
		if err != nil {
			return
		}
		ctx = context.WithValue(ctx, transactionKey, tx)
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()
	} else {
		var ok bool
		tx, ok = ctxTx.(*sql.Tx)
		if !ok {
			err = errors.New("could not cast into *sql.Tx")
		}
	}

	newGroupID, err = adp.q.WithTx(tx).CreateGroup(ctx, groupData.Name)

	for _, permission := range groupData.Permissions {
		err = adp.q.WithTx(tx).AddPermissionsToGroup(ctx, sqlc_gen.AddPermissionsToGroupParams{
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

func (adp *GroupPersistanceAdapter) DeleteGroup(ctx context.Context, groupID uuid.UUID) (err error) {
	ctxTx := ctx.Value(transactionKey)
	var tx *sql.Tx
	if ctxTx == nil {
		tx, err = adp.conn.Begin()
		if err != nil {
			return
		}
		ctx = context.WithValue(ctx, transactionKey, tx)
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()
	} else {
		var ok bool
		tx, ok = ctxTx.(*sql.Tx)
		if !ok {
			err = errors.New("could not cast into *sql.Tx")
		}
	}

	return adp.q.WithTx(tx).DeleteGroup(ctx, groupID)
}
