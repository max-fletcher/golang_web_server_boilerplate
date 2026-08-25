-- The syntax for a sqlc query is like this->->name: {funcName} :{noOfRecordsToReturn}
-- After defining your schema, go to your project root(where sqlc.yaml is) and use "sqlc generate" to generate the functions

-- name: CreatePost :one
INSERT INTO posts (id, title, content, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetPostById :one
SELECT * 
FROM posts 
WHERE id = $1
LIMIT 1;

-- name: GetPosts :one
SELECT * 
FROM posts 
ORDER BY created_at DESC;

-- name: UpdatePost :one
UPDATE posts
SET title = $2,
    content = $3,
    user_id = $4,
    updated_at = $5
WHERE id = $1
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1
RETURNING *;

-- name: SearchPostsByTitleStrSearch :many
SELECT * 
FROM posts 
WHERE title ILIKE '%' || $1 || '%'
ORDER BY created_at DESC;