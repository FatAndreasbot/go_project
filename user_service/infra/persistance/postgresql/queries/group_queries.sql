
-- name: GetGroupByID :one
select
    g.name,
    g.id as "id",
    JSON_AGG(
        JSON_BUILD_OBJECT(
            'ID', p.id,
            'Name', p.name
        ) order by p.id
    ) as permissions
from
    groups g
    left join group_permissions gp on g.id = gp.group_id
    left join permissions p on p.id = gp.permission_id
where
    g.id = $1
group by g.name, g.id;

-- name: GetGroupList :many
select
    g.name,
    g.id as "id",
    JSON_AGG(
        JSON_BUILD_OBJECT(
            'ID', p.id,
            'Name', p.name
        ) order by p.id
    ) as permissions,
    g.id as created
from
    groups g
    left join group_permissions gp on g.id = gp.group_id
    left join permissions p on p.id = gp.permission_id
where g.id >= $1
group by g.name, g.id
limit $2;


-- name: UpdateGroup :exec
update groups g set
    "name" = $1
where g.id = $2;

-- name: AddPermissionsToGroup :exec
insert into group_permissions (group_id, permission_id) values ($1, $2);

-- name: RemovePermissionsFromGroup :exec
delete from group_permissions gp where gp.group_id = $1 and gp.permission_id = $2;

-- name: GetGroupPermissions :many
select
    p.name,
    p.id
from
    permissions p
    join group_permissions gp on gp.permission_id = p.id
where gp.group_id = $1;

-- name: CreateGroup :one
insert into groups (name) values ($1)
returning id;

-- name: DeleteGroup :exec
delete from groups where id = $1;
