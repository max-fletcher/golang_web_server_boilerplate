-- The syntax for a sqlc query is like this->->name: {funcName} :{noOfRecordsToReturn}
-- After defining your schema, go to your project root(where sqlc.yaml is) and use "sqlc generate" to generate the functions

-- name: CreatePost :one
INSERT INTO posts (id, title, content, photo, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetPostById :one
SELECT * 
FROM posts 
WHERE id = $1
LIMIT 1;

-- name: GetPosts :many
SELECT * 
FROM posts 
WHERE
    ($1 = '' OR title ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: GetPostsCount :one
SELECT COUNT(*)
FROM posts 
WHERE
    ($1 = '' OR title ILIKE '%' || $1 || '%');

-- name: UpdatePost :one
UPDATE posts
SET title = $2,
    content = $3,
    photo = $4,
    user_id = $5,
    updated_at = $6
WHERE id = $1
RETURNING *;

-- name: DeletePost :execrows
DELETE FROM posts
WHERE id = $1
RETURNING *;