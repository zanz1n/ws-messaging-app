-- name: CreateUser :exec
INSERT INTO "user" ("id", "username", "password", "updatedAt", "role") VALUES ($1, $2, $3, $4, $5);

-- name: GetUserByUsername :one
SELECT * FROM "user" WHERE "username" = $1;

-- name: GetAllMessages :many
SELECT * FROM "message" LIMIT $1;

-- name: GetMessagesByUser :many
SELECT * FROM "message" WHERE "userId" = $1 LIMIT $2;
