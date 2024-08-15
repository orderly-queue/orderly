// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package queries

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users(id, name, email, password)
VALUES ($1, $2, $3, $4)
RETURNING id, name, email, password, admin, created_at, updated_at, deleted_at
`

type CreateUserParams struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Password string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (*User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Password,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Admin,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return &i, err
}

const deleteUser = `-- name: DeleteUser :exec
UPDATE users
SET deleted_at = EXTRACT(epoch FROM NOW())
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, name, email, password, admin, created_at, updated_at, deleted_at
FROM users
WHERE email = $1
LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Admin,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return &i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, name, email, password, admin, created_at, updated_at, deleted_at
FROM users
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetUserById(ctx context.Context, id uuid.UUID) (*User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Admin,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return &i, err
}

const makeAdmin = `-- name: MakeAdmin :one
UPDATE users
SET admin = true, updated_at = $2
WHERE id = $1
RETURNING id, name, email, password, admin, created_at, updated_at, deleted_at
`

type MakeAdminParams struct {
	ID        uuid.UUID
	UpdatedAt int64
}

func (q *Queries) MakeAdmin(ctx context.Context, arg MakeAdminParams) (*User, error) {
	row := q.db.QueryRowContext(ctx, makeAdmin, arg.ID, arg.UpdatedAt)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Admin,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return &i, err
}

const removeAdmin = `-- name: RemoveAdmin :one
UPDATE users
SET admin = false, updated_at = $2
WHERE id = $1
RETURNING id, name, email, password, admin, created_at, updated_at, deleted_at
`

type RemoveAdminParams struct {
	ID        uuid.UUID
	UpdatedAt int64
}

func (q *Queries) RemoveAdmin(ctx context.Context, arg RemoveAdminParams) (*User, error) {
	row := q.db.QueryRowContext(ctx, removeAdmin, arg.ID, arg.UpdatedAt)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Password,
		&i.Admin,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return &i, err
}
