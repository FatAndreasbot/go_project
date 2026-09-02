-- name: GetUserByUsername :one
select
    u.username,
    u.password_hash,
    u."id",
    g."name" as group_name,
    g."id" as group_uuid,
    JSON_AGG(
        JSON_BUILD_OBJECT(
            'ID', p.id,
            'Name', p.name
        ) order by p.id
    ) as permissions
from
    users u
    left join "groups" g on g.id = u.group_id
    left join group_permissions gp on g.id = gp.group_id
    left join permissions p on p.id = gp.permission_id
where
	u.username = $1
group by
    u.username,
    u.password_hash,
    u."id",
    g."name",
    g."id";

-- name: GetUserByID :one
select
    u.username,
    u.password_hash,
    u."id",
    g."name" as group_name,
    g."id" as group_uuid,
    JSON_AGG(
        JSON_BUILD_OBJECT(
            'ID', p.id,
            'Name', p.name
        ) order by p.id
    ) as permissions
from
    users u
    left join "groups" g on g.id = u.group_id
    left join group_permissions gp on g.id = gp.group_id
    left join permissions p on p.id = gp.permission_id
where
	u."id" = $1
group by
    u.username,
    u.password_hash,
    u."id",
    g."name",
    g."id";

-- name: GetUserList :many
select
    u.username,
    u.password_hash,
    u."id",
    g."name" as group_name,
    g."id" as group_uuid,
	JSON_AGG(
        JSON_BUILD_OBJECT(
            'ID', p.id,
            'Name', p.name
        ) order by p.id
    ) as permissions
from
    users u
    left join "groups" g on g.id = u.group_id
    left join group_permissions gp on g.id = gp.group_id
    left join permissions p on p.id = gp.permission_id
group by
	u.id,
    u.username,
    u.password_hash,
    u."id",
    g."name",
    g."id"
order by u.id
limit $1 offset $2;

-- name: UpdateUser :exec
update Users u set
    "username" = $1,
    "password_hash" = $2,
    group_id = $3
where
    u.id = $4;

-- name: CreateUser :one
insert into Users (id, username, password_hash, group_id) values ($1, $2, $3, $4)
returning id;

-- name: DeleteUser :exec
delete from users u where u.id = $1;
