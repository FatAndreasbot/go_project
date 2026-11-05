
-- name: GetGroupDataByID :one
select
    g.name,
    g.id as "id",
    g.created_at
from
    groups g
where
    g.id = $1;

-- name: GetGroupDataByName :one
select
    g.name,
    g.id as "id",
    g.created_at
from
    groups g
where
    g.name = $1;


-- name: GetGroupList :many
select
    g.name,
    g.id as "id",
    g.created_at
from
    groups g
where g.created_at >= $1
and g.id <> $2
order by g.created_at
limit $3;


-- name: UpdateGroupName :exec
update groups g set
    "name" = $1
where g.id = $2;

-- name: AddPermissionsToGroup :exec
insert into group_permissions (group_id, permission_id) 
select 
    @group_id::uuid,
    unnest(@permission_ids::uuid[]);

-- name: RemovePermissionsFromGroup :exec
delete from 
    group_permissions gp 
where 
    gp.group_id = @group_id::uuid
    and gp.permission_id = ANY(@permission_ids::uuid[]);

-- name: GetGroupPermissions :many
select
    p.name,
    p.id
from
    permissions p
    join group_permissions gp on gp.permission_id = p.id
where gp.group_id = $1;

-- name: CreateGroup :one
insert into groups (name, idempotency_key) values ($1, $2)
returning id, created_at;

-- name: DeleteGroup :exec
delete from groups where id = $1;

-- name: CheckGroupWasCreated :one
select (count(*)!=0) as was_created from groups g 
where g.idempotency_key = $1;