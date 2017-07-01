package models

import (
	"time"

	"github.com/google/uuid"
)

// User struct based on users table in database
type User struct {
	ID           uuid.UUID `db:"id"`
	Uname        string    `db:"uname"`
	Digest       []byte    `db:"digest"`
	Role         string    `db:"role"`
	Email        string    `db:"email"`
	GpgKey       string    `db:"email"`
	LastOnlineAt time.Time `db:"updated_at"`
	CreatedAt    time.Time `db:"created_at"`
}
