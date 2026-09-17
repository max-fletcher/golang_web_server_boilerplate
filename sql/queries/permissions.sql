-- The syntax for a sqlc query is like this->->name: {funcName} :{noOfRecordsToReturn}
-- After defining your schema, go to your project root(where sqlc.yaml is) and use "sqlc generate" to generate the functions

-- name: CreatePermission :one
INSERT INTO permissions (id, name, module_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetPermissionsById :one
SELECT * 
FROM permissions 
WHERE id = $1
LIMIT 1;

-- name: GetPermissions :many
SELECT * 
FROM permissions 
WHERE
    ($1 = '' OR name ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: GetPermissionsCount :one
SELECT COUNT(*)
FROM permissions 
WHERE
    ($1 = '' OR name ILIKE '%' || $1 || '%');