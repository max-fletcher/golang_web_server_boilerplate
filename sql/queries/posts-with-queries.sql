-- The syntax for a sqlc query is like this->->name: {funcName} :{noOfRecordsToReturn}
-- After defining your schema, go to your project root(where sqlc.yaml is) and use "sqlc generate" to generate the functions

-- name: GetPostWithUserById :one
SELECT 
  posts.id,
  posts.title,
  posts.content,
  posts.photo,
  posts.user_id,
  posts.created_at AS post_created_at,
  posts.updated_at AS post_updated_at,
  users.id AS user_id,
  users.name AS user_name,
  users.email AS user_email
FROM posts 
INNER JOIN users 
ON posts.user_id = users.id
WHERE posts.id = $1
LIMIT 1;

-- name: GetPostsWithUser :many
SELECT 
  posts.id,
  posts.title,
  posts.content,
  posts.photo,
  posts.user_id,
  posts.created_at AS post_created_at,
  posts.updated_at AS post_updated_at,
  users.id AS user_id,
  users.name AS user_name,
  users.email AS user_email
FROM posts 
INNER JOIN users 
ON posts.user_id = users.id
WHERE
    ($1 = '' OR posts.title ILIKE '%' || $1 || '%')
ORDER BY posts.created_at DESC
LIMIT $2
OFFSET $3;