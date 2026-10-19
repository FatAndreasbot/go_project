package postgresrepository

import (
	"context"
	"database/sql"
	"log"

	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/infra/postgresql/sqlc_gen"
	"github.com/google/uuid"
)

type PermissionPersistanceAdapter struct {
	conn *sql.DB
	q    *sqlc_gen.Queries
}

func (adp *PermissionPersistanceAdapter) GetPermissionByID(ctx context.Context, permissionID uuid.UUID) (*models.Permission, error) {
	permissionData, err := adp.q.GetPermissionByID(ctx, permissionID)
	if err != nil {
		log.Default().Println(err)
		return nil, err
	}

	return &models.Permission{
		ID:   permissionData.ID,
		Name: models.PermissionName(permissionData.Name),
	}, nil
}

func (adp *PermissionPersistanceAdapter) GetPermissions(ctx context.Context) ([]*models.Permission, error) {
	permissionData, err := adp.q.GetAllPermissions(ctx)
	if err != nil {
		log.Default().Println(err)
		return nil, err
	}

	permissions := make([]*models.Permission, 0, len(permissionData))
	for _, permissionRecord := range permissionData {
		permissions = append(permissions, &models.Permission{
			ID:   permissionRecord.ID,
			Name: models.PermissionName(permissionRecord.Name),
		})
	}

	return permissions, nil
}
