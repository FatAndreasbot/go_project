package postgresrepository

import (
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
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
