-- The syntax for a sqlc query is like this->->name: {funcName} :{noOfRecordsToReturn}
-- After defining your schema, go to your project root(where sqlc.yaml is) and use "sqlc generate" to generate the functions
-- The query:GetUsersWithRolesAndPermissions is using CTE
-- The query:UserHasPermission returns a boolean

-- name: CreateRolePermission :one
INSERT INTO role_permissions (id, role_id, permission_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetRolePermissions :many
SELECT
    rp.id,
    rp.role_id,
    rp.permission_id,
    m.name AS module_name,
    r.name AS role_name,
    p.name AS permission_name,
    rp.created_at as created_at,
    rp.updated_at as updated_at
FROM role_permissions AS rp
INNER JOIN roles AS r
    ON r.id = rp.role_id
INNER JOIN permissions AS p
    ON p.id = rp.permission_id
INNER JOIN modules AS m
    ON m.id = p.module_id
WHERE
    $1 = ''
    OR r.name ILIKE '%' || $1 || '%'
    OR m.name ILIKE '%' || $1 || '%'
    OR p.name ILIKE '%' || $1 || '%'
ORDER BY rp.created_at DESC
LIMIT $2
OFFSET $3;

-- name: GetRolePermissionsCount :one
SELECT COUNT(*)
FROM role_permissions AS rp
INNER JOIN roles AS r
    ON r.id = rp.role_id
INNER JOIN permissions AS p
    ON p.id = rp.permission_id
INNER JOIN modules AS m
    ON m.id = p.module_id
WHERE
    $1 = ''
    OR r.name ILIKE '%' || $1 || '%'
    OR m.name ILIKE '%' || $1 || '%'
    OR p.name ILIKE '%' || $1 || '%';

-- name: GetUsersWithRolesAndPermissions :many
WITH paginated_users AS (
    SELECT DISTINCT
        u.id,
        u.name,
        u.email,
        u.created_at,
        u.updated_at
    FROM users AS u
    INNER JOIN user_roles AS ur
        ON u.id = ur.user_id
    INNER JOIN roles AS r
        ON r.id = ur.role_id
    INNER JOIN role_permissions AS rp
        ON r.id = rp.role_id
    INNER JOIN permissions AS p
        ON p.id = rp.permission_id
    INNER JOIN modules AS m
        ON m.id = p.module_id
    WHERE
        $1 = ''
        OR u.name ILIKE '%' || $1 || '%'
        OR r.name ILIKE '%' || $1 || '%'
        OR m.name ILIKE '%' || $1 || '%'
        OR p.name ILIKE '%' || $1 || '%'
    ORDER BY u.created_at DESC
    LIMIT $2
    OFFSET $3
)
SELECT
    u.id,
    u.name,
    u.email,
    u.created_at,
    u.updated_at,
    r.name AS role_name,
    p.name AS permission_name,
    m.name AS module_name
FROM paginated_users AS u
INNER JOIN user_roles AS ur
    ON u.id = ur.user_id
INNER JOIN roles AS r
    ON r.id = ur.role_id
INNER JOIN role_permissions AS rp
    ON r.id = rp.role_id
INNER JOIN permissions AS p
    ON p.id = rp.permission_id
INNER JOIN modules AS m
    ON m.id = p.module_id
ORDER BY u.created_at DESC;

-- name: GetUsersRolePermissionCount :one
SELECT COUNT(DISTINCT u.id)
FROM users AS u
INNER JOIN user_roles AS ur
    ON u.id = ur.user_id
INNER JOIN roles AS r
    ON r.id = ur.role_id
INNER JOIN role_permissions AS rp
    ON r.id = rp.role_id
INNER JOIN permissions AS p
    ON p.id = rp.permission_id
INNER JOIN modules AS m
    ON m.id = p.module_id
WHERE
    $1 = ''
    OR u.name ILIKE '%' || $1 || '%'
    OR r.name ILIKE '%' || $1 || '%'
    OR m.name ILIKE '%' || $1 || '%'
    OR p.name ILIKE '%' || $1 || '%';

-- name: GetRolePermissionByID :one
SELECT rp.id, rp.role_id, rp.permission_id, r.name as role_name, p.name as permission_name, m.name as module_name
FROM role_permissions as rp
INNER JOIN roles AS r
  ON r.id = rp.role_id
INNER JOIN permissions AS p
  ON p.id = rp.permission_id
INNER JOIN modules AS m
  ON m.id = p.module_id
WHERE rp.id = $1;

-- name: GetRolePermissionByRoleIDAndPermissionID :one
SELECT rp.id, rp.role_id, rp.permission_id, r.name as role_name, p.name as permission_name, m.name as module_name
FROM role_permissions as rp
INNER JOIN roles AS r
  ON r.id = rp.role_id
INNER JOIN permissions AS p
  ON p.id = rp.permission_id
INNER JOIN modules AS m
  ON m.id = p.module_id
WHERE rp.role_id = $1 AND rp.permission_id = $2;

-- name: GetUserWithRolesAndPermissionsByUserID :many
SELECT
    u.id,
    u.name,
    u.email,
    u.created_at,
    u.updated_at,
    r.name AS role_name,
    p.name AS permission_name,
    m.name AS module_name
FROM users AS u
LEFT JOIN user_roles AS ur
    ON u.id = ur.user_id
LEFT JOIN roles AS r
    ON r.id = ur.role_id
LEFT JOIN role_permissions AS rp
    ON r.id = rp.role_id
LEFT JOIN permissions AS p
    ON p.id = rp.permission_id
LEFT JOIN modules AS m
    ON m.id = p.module_id
WHERE u.id = $1
ORDER BY r.name, m.name, p.name;

-- name: UserHasPermission :one
SELECT EXISTS (
    SELECT 1
    FROM user_roles ur
    INNER JOIN role_permissions rp
        ON rp.role_id = ur.role_id
    INNER JOIN permissions p
        ON p.id = rp.permission_id
    INNER JOIN modules m
        ON m.id = p.module_id
    WHERE ur.user_id = $1
      AND m.name = $2
      AND p.name = $3
);

-- name: DeleteRolePermissionByID :execrows
DELETE FROM role_permissions
WHERE id = $1
RETURNING *;

-- name: DeleteRolePermissionByRoleIDAndPermissionID :execrows
DELETE FROM role_permissions
WHERE role_id = $1 AND permission_id = $2
RETURNING *;