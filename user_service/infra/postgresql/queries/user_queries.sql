-- name: GetUserDataByEmail :one
select
    u.id,
    u.email,
    u.password_hash,
    u.group_id,
    u.created_at,
    u.idempotency_key
from
    users u
where
	u.email = $1;

-- name: GetUserDataByID :one
select
    u.id,
    u.email,
    u.password_hash,
    u.group_id,
    u.created_at,
    u.idempotency_key
from
    users u
where
	u.id = $1;

-- name: GetUserList :many
select
    u.id,
    u.email,
    u.password_hash,
    u.group_id,
    u.created_at,
    u.idempotency_key
from
    users u
WHERE
    u.created_at >= $1
    and u.id <> $2
order by u.created_at
limit $3;

-- name: UpdateUser :one
update Users u set
    email = $1,
    password_hash = $2,
    group_id = $3
where
    u.id = $4
RETURNING u.created_at;

-- name: CreateUser :one
insert into Users (email, password_hash, group_id, idempotency_key) values ($1, $2, $3, $4)
returning id, created_at;

-- name: DeleteUser :exec
delete from users u where u.id = $1;

-- name: CheckUserWasCreated :one
select (count(*)!=0) as was_created from users u 
where u.idempotency_key = $1;
