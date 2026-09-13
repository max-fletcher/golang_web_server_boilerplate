-- The syntax for a sqlc query is like this->->name: {funcName} :{noOfRecordsToReturn}
-- After defining your schema, go to your project root(where sqlc.yaml is) and use "sqlc generate" to generate the functions

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetRefreshTokenByTokenHash :one
SELECT * 
FROM refresh_tokens 
WHERE token_hash = $1
LIMIT 1;

-- name: UpdateRefreshToken :one
UPDATE refresh_tokens
SET token_hash = $2,
    expires_at = $3,
    revoked_at = $4,
    updated_at = $5
WHERE user_id = $1
RETURNING *;

-- name: DeleteRefreshToken :execrows
DELETE FROM refresh_tokens
WHERE id = $1
RETURNING *;

-- name: DeleteRefreshTokenByUserId :execrows
DELETE FROM refresh_tokens
WHERE user_id = $1
RETURNING *;

-- name: DeleteAllRefreshTokens :exec
TRUNCATE TABLE refresh_tokens;