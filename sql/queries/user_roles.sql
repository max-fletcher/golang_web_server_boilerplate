-- The syntax for a sqlc query is like this->->name: {funcName} :{noOfRecordsToReturn}
-- After defining your schema, go to your project root(where sqlc.yaml is) and use "sqlc generate" to generate the functions

-- name: CreateUserRole :one
INSERT INTO user_roles (id, user_id, role_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserRoleById :one
SELECT u.id, u.name, u.email, u.created_at, u.updated_at, r.name as role_name
FROM user_roles AS ur
INNER JOIN users AS u 
  ON u.id = ur.user_id
INNER JOIN roles AS r 
  ON r.id = ur.role_id
WHERE ur.id = $1;

-- name: GetUserRolesByUserID :many
SELECT u.id, u.name, u.email, u.created_at, u.updated_at, r.name as role_name
FROM users AS u
LEFT JOIN user_roles AS ur
  ON u.id = ur.user_id
LEFT JOIN roles AS r 
  ON r.id = ur.role_id
WHERE u.id = $1;

-- name: GetUserRoleByUserIDAndRoleID :one
SELECT
    u.id,
    u.name,
    u.email,
    u.created_at,
    u.updated_at,
    r.name AS role_name
FROM user_roles AS ur
INNER JOIN users AS u
    ON u.id = ur.user_id
INNER JOIN roles AS r
    ON r.id = ur.role_id
WHERE ur.user_id = $1
  AND ur.role_id = $2;

-- name: GetUserRoles :many
WITH paginated_users AS (
  SELECT u.id, u.name, u.email, u.created_at, u.updated_at
  FROM users AS u
  WHERE
    $1 = ''
    OR u.name ILIKE '%' || $1 || '%'
    OR EXISTS (
        SELECT 1
        FROM user_roles AS ur
        INNER JOIN roles AS r
          ON r.id = ur.role_id
        WHERE ur.user_id = u.id
          AND r.name ILIKE '%' || $1 || '%'
    )
  ORDER BY u.created_at DESC
  LIMIT $2
  OFFSET $3
)
SELECT u.id, u.name, u.email, u.created_at, u.updated_at, ur.role_id, r.name as role_name
FROM paginated_users AS u
LEFT JOIN user_roles AS ur 
  ON u.id = ur.user_id
LEFT JOIN roles AS r 
  ON r.id = ur.role_id
ORDER BY u.created_at DESC;

-- name: GetUserRolesCount :one
SELECT COUNT(*)
FROM users AS u
WHERE
  $1 = ''
  OR u.name ILIKE '%' || $1 || '%'
  OR EXISTS (
      SELECT 1
      FROM user_roles AS ur
      INNER JOIN roles AS r
        ON r.id = ur.role_id
      WHERE ur.user_id = u.id
        AND r.name ILIKE '%' || $1 || '%'
  );

-- name: DeleteUserRoleById :execrows
DELETE FROM user_roles
WHERE id = $1
RETURNING *;

-- name: DeleteUserRoleByUserIDAndRoleId :execrows
DELETE FROM user_roles
WHERE user_id = $1 AND role_id = $2
RETURNING *;