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

// Image struct based on image table in database
type Image struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	URL       string    `db:"url"`
	Medium    string    `db:"medium"`
	Small     string    `db:"small"`
	Caption   string    `db:"caption"`
	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt time.Time `db:"created_at"`
}

// Tag struct based on tag table in database
type Tag struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
	Slug string    `db:"slug"`
}

// FindPost ...
func (db *DB) FindPost(id uuid.UUID) (*Post, error) {
	p := new(Post)
	sql := "SELECT * FROM posts WHERE id = $1"
	err := db.Get(p, sql, id)
	return p, err
}

// FindPostsByUser ...
func (db *DB) FindPostsByUser(user uuid.UUID) (*[]Post, error) {
	p := new([]Post)
	sql := "SELECT * FROM posts WHERE user_id = $1"
	err := db.Select(p, sql, user)
	return p, err
}

// InsertPost ...
func (db *DB) InsertPost(user uuid.UUID, title string, slug string, subtitle string, short string, content string, digest string, published bool) (*Post, error) {
	p := new(Post)
	sql := "INSERT INTO posts (user_id, title, slug, subtitle, short, content, digest, published) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *"
	err := db.Get(p, sql, user, title, slug, subtitle, short, content, digest, published)
	return p, err
}

// UpdatePost ...
func (db *DB) UpdatePost(id uuid.UUID, user uuid.UUID, title string, slug string, subtitle string, short string, content string, digest string, published bool) (*Post, error) {
	p := new(Post)
	sql := "UPDATE posts SET (user_id, title, slug, subtitle, short, content, digest, published) VALUES ($2, $3, $4, $5, $6, $7, $8) WHERE id = $1 RETURNING *"
	err := db.Get(p, sql, id, user, title, slug, subtitle, short, content, digest, published)
	return p, err
}

// DeletePost ...
func (db *DB) DeletePost(id uuid.UUID) (*Post, error) {
	p := new(Post)
	sql := "DELETE FROM posts WHERE id = $1 RETURNING *"
	err := db.Get(p, sql, id)
	return p, err
}

// FindImage 
func (db *DB) FindImage(id uuid.UUID) (*Image, error) {
	i := new(Image)
	sql := "SELECT * FROM images WHERE id = $1"
	err := db.Get(i, sql, id)
	return i, err
}

// FindImagesByUser returns an slice of images from the database for a given user
func (db *DB) FindImagesByUser(user uuid.UUID) (*[]Image, error) {
	i := new([]Image)
	sql := "SELECT * FROM images WHERE user_id = $1"
	err := db.Select(i, sql, user)
	return i, err
}

// InsertImage attempts to insert an image into the database provided three sizes of the image (original, medium, and small) return an error
// if it cannot be added or user reference invalid
func (db *DB) InsertImage(user uuid.UUID, url string, medium string, small string, caption string) (*Image, error) {
	i := new(Image)
	sql := "INSERT INTO images (user_id, url, medium, small, caption) VALUES ($1, $2, $3, $4, $5) RETURNING *"
	err := db.Get(i, sql, user)
	return i, err
}

//DeleteImage takes an id of an image and if exists deletes from database returns error if not found
func (db *DB) DeleteImage(id uuid.UUID) (*Image, error) {
	i := new(Image)
	sql := "DELETE FROM images WHERE id = $1"
	err := db.Get(i, sql, id)
	return i, err
}
