-- name: CreateUser :exec
INSERT INTO "user" ("id", "username", "password", "role", "updatedAt", "createdAt") VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetUserByUsername :one
SELECT * FROM "user" WHERE "username" = $1;

-- name: GetUserByUsernameSafe :one
SELECT "id", "createdAt", "updatedAt", "role", "username" FROM "user" WHERE "id" = $1;

-- name: GetAllMessages :many
SELECT * FROM "message" LIMIT $1;

-- name: GetMessageById :one
SELECT * FROM "message" WHERE "id" = $1;

-- name: DeleteMessageById :exec
DELETE FROM "message" WHERE "id" = $1;

-- name: CreateMessage :one
INSERT INTO "message" ("id", "userId", "content", "imageUrl", "updatedAt", "createdAt") VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetMessagesByUserId :many
SELECT * FROM "message" WHERE "userId" = $1 LIMIT $2;

-- name: GetMessagesByUsername :many
SELECT * FROM "message" WHERE "userId" = (SELECT "id" FROM "user" WHERE "username" = $1) LIMIT $2;

-- name: GetMessagesWithOffset :many
SELECT * FROM "message" WHERE "createdAt" < $1 ORDER BY "createdAt" DESC LIMIT $2;

-- name: GetUserById :one
SELECT * FROM "user" WHERE "id" = $1;

-- name: GetUserByIdSafe :one
SELECT "id", "createdAt", "updatedAt", "role", "username" FROM "user" WHERE "id" = $1;
