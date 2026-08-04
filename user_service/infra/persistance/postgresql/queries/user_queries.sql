-- name: GetUserByUsername :one
select
    u.username,
    u.password_hash,
    u."uuid",
    g."name",
    g."uuid" as group_uuid,
    JSON_AGG(
        JSON_BUILD_OBJECT(
            p.uuid, p.name
        )
    ) as permissions
from
    users u
    join "groups" g on g.uuid = u.group_id
    join group_permissions gp on g.id = gp.group_id
    join permissions p on p.id = gp.permission_id
where
	u.username = $1
group by 
    u.username,
    u.password_hash,
    u."uuid",
    g."name",
    g."uuid";

-- name: GetUserByID :one
select
    u.username,
    u.password_hash,
    u."uuid",
    g."name",
    g."uuid" as group_uuid,
    JSON_AGG(
        JSON_BUILD_OBJECT(
            p.uuid, p.name
        )
    ) as permissions
from
    users u
    join "groups" g on g.id = u.group_id 
    join group_permissions gp on g.id = gp.group_id
    join permissions p on p.id = gp.permission_id
where
	u."uuid" = $1;

-- name: GetUserList :many
select
    u.username,
    u.password_hash,
    u."uuid",
    g."name",
    g."uuid" as group_uuid
from
    users u
    join "groups" g on g.id = u.group_id  
limit
    $1 offset $2;

-- name: UpdateUser :exec
update Users u set
    "username" = $1,
    "password_hash" = $2,
    group_id = $3
where
    u.uuid = $4;

-- name: CreateUser :one
insert into Users (username, password_hash, group_id) values ($1, $2, $3)
returning uuid;

-- name: DeleteUser :exec
delete from users u where u.uuid = $1;

