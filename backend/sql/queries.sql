-- name: CreateUser :exec
INSERT INTO "user" ("id", "username", "password", "updatedAt", "role") VALUES ($1, $2, $3, $4, $5);

-- name: GetUserByUsername :one
SELECT * FROM "user" WHERE "username" = $1;

-- name: GetAllMessages :many
SELECT * FROM "message" LIMIT $1;

-- name: GetMessageById :one
SELECT * FROM "message" WHERE "id" = $1;

-- name: DeleteMessageById :exec
DELETE FROM "message" WHERE "id" = $1;

-- name: CreateMessage :one
INSERT INTO "message" ("id", "userId", "content", "imageUrl", "updatedAt") VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetMessagesByUserId :many
SELECT * FROM "message" WHERE "userId" = $1 LIMIT $2;

-- name: GetMessagesByUsername :many
SELECT * FROM "message" WHERE "userId" = (SELECT "id" FROM "user" WHERE "username" = $1) LIMIT $2;
