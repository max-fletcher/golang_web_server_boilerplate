-- The syntax for a sqlc query is like this->->name: {funcName} :{noOfRecordsToReturn}
-- After defining your schema, go to your project root(where sqlc.yaml is) and use "sqlc generate" to generate the functions

-- name: CreateUserRole :one
INSERT INTO user_roles (id, user_id, role_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserRoleById :one
SELECT * 
FROM user_roles 
WHERE id = $1
LIMIT 1;

-- name: GetUserRoleByUserId :one
SELECT u.id, u.name, u.email, u.created_at, u.updated_at, r.name
FROM user_roles AS ur
INNER JOIN users AS u ON u.id = ur.user_id
INNER JOIN roles AS r ON r.id = ur.role_id
WHERE ur.user_id = $1
LIMIT 1;

-- name: GetUserRoles :many
SELECT * 
FROM user_roles 
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: GetUserRolesCount :one
SELECT COUNT(*)
FROM user_roles;

-- name: UpdateUserRole :one
UPDATE user_roles
SET user_id = $1,
    role_id = $2,
    updated_at = $3
WHERE id = $4
RETURNING *;

-- name: DeleteUserRoleById :execrows
DELETE FROM user_roles
WHERE id = $1
RETURNING *;

-- name: DeleteUserRoleByUserIdAndRoleId :execrows
DELETE FROM user_roles
WHERE user_id = $1 AND role_id = $2
RETURNING *;