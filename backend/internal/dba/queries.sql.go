// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: queries.sql

package dba

import (
	"context"
	"database/sql"
	"time"
)

const createMessage = `-- name: CreateMessage :one
INSERT INTO "message" ("id", "userId", "content", "imageUrl", "updatedAt") VALUES ($1, $2, $3, $4, $5) RETURNING id, "createdAt", "updatedAt", content, "imageUrl", "userId"
`

type CreateMessageParams struct {
	ID        string
	UserId    string
	Content   sql.NullString
	ImageUrl  sql.NullString
	UpdatedAt time.Time
}

func (q *Queries) CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error) {
	row := q.db.QueryRowContext(ctx, createMessage,
		arg.ID,
		arg.UserId,
		arg.Content,
		arg.ImageUrl,
		arg.UpdatedAt,
	)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Content,
		&i.ImageUrl,
		&i.UserId,
	)
	return i, err
}

const createUser = `-- name: CreateUser :exec
INSERT INTO "user" ("id", "username", "password", "updatedAt", "role") VALUES ($1, $2, $3, $4, $5)
`

type CreateUserParams struct {
	ID        string
	Username  string
	Password  string
	UpdatedAt time.Time
	Role      UserRole
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.ExecContext(ctx, createUser,
		arg.ID,
		arg.Username,
		arg.Password,
		arg.UpdatedAt,
		arg.Role,
	)
	return err
}

const deleteMessageById = `-- name: DeleteMessageById :exec
DELETE FROM "message" WHERE "id" = $1
`

func (q *Queries) DeleteMessageById(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteMessageById, id)
	return err
}

const getAllMessages = `-- name: GetAllMessages :many
SELECT id, "createdAt", "updatedAt", content, "imageUrl", "userId" FROM "message" LIMIT $1
`

func (q *Queries) GetAllMessages(ctx context.Context, limit int32) ([]Message, error) {
	rows, err := q.db.QueryContext(ctx, getAllMessages, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Content,
			&i.ImageUrl,
			&i.UserId,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMessageById = `-- name: GetMessageById :one
SELECT id, "createdAt", "updatedAt", content, "imageUrl", "userId" FROM "message" WHERE "id" = $1
`

func (q *Queries) GetMessageById(ctx context.Context, id string) (Message, error) {
	row := q.db.QueryRowContext(ctx, getMessageById, id)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Content,
		&i.ImageUrl,
		&i.UserId,
	)
	return i, err
}

const getMessagesByUserId = `-- name: GetMessagesByUserId :many
SELECT id, "createdAt", "updatedAt", content, "imageUrl", "userId" FROM "message" WHERE "userId" = $1 LIMIT $2
`

type GetMessagesByUserIdParams struct {
	UserId string
	Limit  int32
}

func (q *Queries) GetMessagesByUserId(ctx context.Context, arg GetMessagesByUserIdParams) ([]Message, error) {
	rows, err := q.db.QueryContext(ctx, getMessagesByUserId, arg.UserId, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Content,
			&i.ImageUrl,
			&i.UserId,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMessagesByUsername = `-- name: GetMessagesByUsername :many
SELECT id, "createdAt", "updatedAt", content, "imageUrl", "userId" FROM "message" WHERE "userId" = (SELECT "id" FROM "user" WHERE "username" = $1) LIMIT $2
`

type GetMessagesByUsernameParams struct {
	Username string
	Limit    int32
}

func (q *Queries) GetMessagesByUsername(ctx context.Context, arg GetMessagesByUsernameParams) ([]Message, error) {
	rows, err := q.db.QueryContext(ctx, getMessagesByUsername, arg.Username, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Content,
			&i.ImageUrl,
			&i.UserId,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, "createdAt", "updatedAt", role, username, password FROM "user" WHERE "username" = $1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Role,
		&i.Username,
		&i.Password,
	)
	return i, err
}
