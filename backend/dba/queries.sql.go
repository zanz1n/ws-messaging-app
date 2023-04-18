package dba

import (
	"context"
	"database/sql"
)

const createMessage = `-- name: CreateMessage :one
INSERT INTO "message" ("id", "userId", "content", "imageUrl", "updatedAt", "createdAt") VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, "createdAt", "updatedAt", content, "imageUrl", "userId"
`

type CreateMessageParams struct {
	ID        string         `json:"id"`
	UserId    string         `json:"userId"`
	Content   sql.NullString `json:"content"`
	ImageUrl  sql.NullString `json:"imageUrl"`
	UpdatedAt int64          `json:"updatedAt"`
	CreatedAt int64          `json:"createdAt"`
}

func (q *Queries) CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error) {
	row := q.db.QueryRowContext(ctx, createMessage,
		arg.ID,
		arg.UserId,
		arg.Content,
		arg.ImageUrl,
		arg.UpdatedAt,
		arg.CreatedAt,
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
INSERT INTO "user" ("id", "username", "password", "role", "updatedAt", "createdAt") VALUES ($1, $2, $3, $4, $5, $6)
`

type CreateUserParams struct {
	ID        string   `json:"id"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Role      UserRole `json:"role"`
	UpdatedAt int64    `json:"updatedAt"`
	CreatedAt int64    `json:"createdAt"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.ExecContext(ctx, createUser,
		arg.ID,
		arg.Username,
		arg.Password,
		arg.Role,
		arg.UpdatedAt,
		arg.CreatedAt,
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
	items := []Message{}
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
	UserId string `json:"userId"`
	Limit  int32  `json:"limit"`
}

func (q *Queries) GetMessagesByUserId(ctx context.Context, arg GetMessagesByUserIdParams) ([]Message, error) {
	rows, err := q.db.QueryContext(ctx, getMessagesByUserId, arg.UserId, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Message{}
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
	Username string `json:"username"`
	Limit    int32  `json:"limit"`
}

func (q *Queries) GetMessagesByUsername(ctx context.Context, arg GetMessagesByUsernameParams) ([]Message, error) {
	rows, err := q.db.QueryContext(ctx, getMessagesByUsername, arg.Username, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Message{}
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

const getMessagesWithOffset = `-- name: GetMessagesWithOffset :many
SELECT id, "createdAt", "updatedAt", content, "imageUrl", "userId" FROM "message" WHERE "createdAt" < $1 ORDER BY "createdAt" DESC LIMIT $2
`

type GetMessagesWithOffsetParams struct {
	CreatedAt int64 `json:"createdAt"`
	Limit     int32 `json:"limit"`
}

func (q *Queries) GetMessagesWithOffset(ctx context.Context, arg GetMessagesWithOffsetParams) ([]Message, error) {
	rows, err := q.db.QueryContext(ctx, getMessagesWithOffset, arg.CreatedAt, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Message{}
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

const getUserById = `-- name: GetUserById :one
SELECT id, "createdAt", "updatedAt", role, username, password FROM "user" WHERE "id" = $1
`

func (q *Queries) GetUserById(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
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

const getUserByIdSafe = `-- name: GetUserByIdSafe :one
SELECT "id", "createdAt", "updatedAt", "role", "username" FROM "user" WHERE "id" = $1
`

type GetUserByIdSafeRow struct {
	ID        string   `json:"id"`
	CreatedAt int64    `json:"createdAt"`
	UpdatedAt int64    `json:"updatedAt"`
	Role      UserRole `json:"role"`
	Username  string   `json:"username"`
}

func (q *Queries) GetUserByIdSafe(ctx context.Context, id string) (GetUserByIdSafeRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByIdSafe, id)
	var i GetUserByIdSafeRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Role,
		&i.Username,
	)
	return i, err
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

const getUserByUsernameSafe = `-- name: GetUserByUsernameSafe :one
SELECT "id", "createdAt", "updatedAt", "role", "username" FROM "user" WHERE "id" = $1
`

type GetUserByUsernameSafeRow struct {
	ID        string   `json:"id"`
	CreatedAt int64    `json:"createdAt"`
	UpdatedAt int64    `json:"updatedAt"`
	Role      UserRole `json:"role"`
	Username  string   `json:"username"`
}

func (q *Queries) GetUserByUsernameSafe(ctx context.Context, id string) (GetUserByUsernameSafeRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsernameSafe, id)
	var i GetUserByUsernameSafeRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Role,
		&i.Username,
	)
	return i, err
}
