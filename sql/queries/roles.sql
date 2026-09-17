-- The syntax for a sqlc query is like this->->name: {funcName} :{noOfRecordsToReturn}
-- After defining your schema, go to your project root(where sqlc.yaml is) and use "sqlc generate" to generate the functions

-- name: CreateRole :one
INSERT INTO roles (id, name, created_at, updated_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetRoleById :one
SELECT * 
FROM roles 
WHERE id = $1
LIMIT 1;

-- name: GetRoles :many
SELECT * 
FROM roles 
WHERE
    ($1 = '' OR name ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: GetRolesCount :one
SELECT COUNT(*)
FROM roles 
WHERE
    ($1 = '' OR name ILIKE '%' || $1 || '%');

-- name: UpdateRole :one
UPDATE roles
SET name = $2,
    updated_at = $3
WHERE id = $1
RETURNING *;

-- name: DeleteRole :execrows
DELETE FROM roles
WHERE id = $1
RETURNING *;