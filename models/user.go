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

// InsertUser ...
func (db *DB) InsertUser(uname string, digest []byte, role string, email string, gpg string) (*User, error) {
	var u User
	sql := "INSERT INTO users (uname, digest, role, email, gpg) VALUES ($1, $2, $3, $4, $5) RETURNING *"
	err := db.Get(&u, sql, uname, digest, role, email, gpg)
	return &u, err
}

// GetUserByEmail ...
func (db *DB) GetUserByEmail(email string) (*User, error) {
	var u User
	sql := "SELECT * FROM users WHERE email = $1"
	err := db.Get(&u, sql, email)
	return &u, err
}

// GetUserByID ...
func (db *DB) GetUserByID(user uuid.UUID) (*User, error) {
	var u User
	sql := "SELECT * FROM users WHERE id = $1"
	err := db.Get(&u, sql, user)
	return &u, err
}
