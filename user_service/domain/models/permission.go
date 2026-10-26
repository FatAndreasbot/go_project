package models

import "github.com/google/uuid"

type PermissionName string

// list of all permission names
const (
	Any PermissionName = "any"

	UserCreate PermissionName = "user_create"
	UserRead   PermissionName = "user_read"
	UserUpdate PermissionName = "user_update"
	UserDelete PermissionName = "user_delete"

	GroupCreate PermissionName = "group_create"
	GroupRead   PermissionName = "group_read"
	GroupUpdate PermissionName = "group_update"
	GroupDelete PermissionName = "group_delete"
)

type Permission struct {
	ID   uuid.UUID
	Name PermissionName
}
