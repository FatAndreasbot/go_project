
-- name: GetGroupByID :one
select 
    g.name,
    g.uuid as "id",
    JSON_AGG(
        JSON_BUILD_OBJECT(
            p.uuid, p.name
        )
    )
from 
    groups g
    join group_permissions gp on g.id = gp.group_id
    join permissions p on p.id = gp.permission_id
where 
    g.uuid = $1
group by g.name, g.uuid;

-- name: GetGroupList :many
select 
    g.name,
    g.uuid as "id",
    JSON_AGG(
        JSON_BUILD_OBJECT(
            p.uuid, p.name
        )
    )
from 
    groups g
    join group_permissions gp on g.id = gp.group_id
    join permissions p on p.id = gp.permission_id
where 
    g.uuid = $1
group by g.name, g.uuid;

-- name: UpdateGroup :exec
update groups g set 
    "name" = $1;

-- name: AddPermissionsToGroup :exec
insert into group_permissions (group_id, permission_id) values ($1, $2);

-- name: RemovePermissionsFromGroup :exec
delete from group_permissions gp where gp.group_id = $1 and gp.permission_id = $2;

-- name: GetGroupPermissions :many
select
    p.name,
    p.uuid
from 
    permissions p
    join group_permissions gp on gp.permission_id = p.id
where gp.group_id = $1;

-- name: CreateGroup :one
insert into groups (name) values ($1)
returning uuid;

-- name: DeleteGroup :exec
delete from groups where uuid = $1;