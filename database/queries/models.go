// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package queries

import (
	"database/sql"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	Admin     bool
	CreatedAt int64
	UpdatedAt int64
	DeletedAt sql.NullInt64
}
