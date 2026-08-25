-- The syntax for a sqlc query is like this->->name: {funcName} :{noOfRecordsToReturn}
-- After defining your schema, go to your project root(where sqlc.yaml is) and use "sqlc generate" to generate the functions

-- name: CreateUser :one
INSERT INTO users (id, name, email, password, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserById :one
SELECT * 
FROM users 
WHERE id = $1
LIMIT 1;

-- name: GetUsers :many
SELECT * 
FROM users 
ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users
SET name = $2,
    email = $3,
    password = $4,
    updated_at = $5
WHERE id = $1
RETURNING *;

-- name: DeleteUser :execrows
DELETE FROM users
WHERE id = $1
RETURNING *;

-- name: GetUserByEmail :one
SELECT * 
FROM users 
WHERE email = $1
LIMIT 1;