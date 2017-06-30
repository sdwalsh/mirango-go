package models

import (
	"time"

	"github.com/google/uuid"
)

// Post struct based on posts table in database
type Post struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	Title       string    `db:"title"`
	Slug        string    `db:"slug"`
	SubTitle    string    `db:"sub_title"`
	Short       string    `db:"short"`
	PostContent string    `db:"post_content"`
	Digest      string    `db:"digest"`
	Published   bool      `db:"published"`
	UpdatedAt   time.Time `db:"updated_at"`
	CreatedAt   time.Time `db:"created_at"`
}
