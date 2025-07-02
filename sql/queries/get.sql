-- name: ListUsers :many
SELECT * FROM USERS;

-- name: AddUser :one
INSERT INTO USERS (username, hash)
VALUES ($1, $2)
RETURNING *;

-- name: GetUser :one
SELECT * FROM USERS
WHERE username = $1;

-- name: CreateUserPost :one
INSERT INTO POSTS (post, username)
VALUES ($1, $2)
RETURNING *;

-- name: CreateUserComment :one
INSERT INTO COMMENTS (comment, post_id, username)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetFeed :many
SELECT 
  POSTS.post_id,
  POSTS.post,
  POSTS.username AS username_post,
  COALESCE(array_agg(COMMENTS.comment) FILTER (WHERE COMMENTS.comment IS NOT NULL), '{}') AS comments,
  COALESCE(array_agg(COMMENTS.username) FILTER (WHERE COMMENTS.username IS NOT NULL), '{}') AS commenter_usernames
FROM POSTS
LEFT JOIN COMMENTS ON COMMENTS.post_id = POSTS.post_id
GROUP BY POSTS.post_id, POSTS.post, POSTS.username;