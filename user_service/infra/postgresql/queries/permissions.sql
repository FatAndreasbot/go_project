-- name: GetPermissionByID :one
select 
	p.name,
	p.id 
from 
	permissions p
where 
	p.id = $1;


-- name: GetAllPermissions :many
select 
	p.name,
	p.id 
from 
	permissions p;
