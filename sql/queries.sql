-- name: GetUserByUsername :one
SELECT * FROM "user" WHERE "username" = $1;

-- name: GetAllMessages :many
SELECT * FROM "message" LIMIT $1;

-- name: GetMessagesByUser :many
SELECT * FROM "message" WHERE "userId" = $1 LIMIT $2;
