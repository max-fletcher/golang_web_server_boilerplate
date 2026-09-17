-- The syntax for a sqlc query is like this->->name: {funcName} :{noOfRecordsToReturn}
-- After defining your schema, go to your project root(where sqlc.yaml is) and use "sqlc generate" to generate the functions

-- name: CreateModule :one
INSERT INTO modules (id, name, created_at, updated_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetModuleById :one
SELECT * 
FROM modules 
WHERE id = $1
LIMIT 1;

-- name: GetModules :many
SELECT * 
FROM modules 
WHERE
    ($1 = '' OR name ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: GetModulesCount :one
SELECT COUNT(*)
FROM modules 
WHERE
    ($1 = '' OR name ILIKE '%' || $1 || '%');