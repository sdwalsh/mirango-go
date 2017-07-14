package models

import (
	"time"

	"github.com/google/uuid"
)

// Post struct based on posts table in database
type Post struct {
	ID          uuid.UUID `db:"id" json:"id"`
	UserID      uuid.UUID `db:"user_id" json:"user_id"`
	Title       string    `db:"title" json:"title"`
	Slug        string    `db:"slug" json:"slug"`
	SubTitle    string    `db:"sub_title" json:"sub_title"`
	Short       string    `db:"short" json:"short"`
	PostContent string    `db:"post_content" json:"post_content"`
	Digest      string    `db:"digest" json:"digest"`
	Published   bool      `db:"published" json:"published"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

// Image struct based on image table in database
type Image struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	URL       string    `db:"url" json:"url"`
	Medium    string    `db:"medium" json:"medium"`
	Small     string    `db:"small" json:"small"`
	Caption   string    `db:"caption" json:"caption"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Tag struct based on tag table in database
type Tag struct {
	ID   uuid.UUID `db:"id" json:"id"`
	Name string    `db:"name" json:"name"`
	Slug string    `db:"slug" json:"slug"`
}

////////////////////
// Post Functions //
////////////////////

// PublishedPosts returns all published posts in database
func (db *DB) PublishedPosts() (*[]Post, error) {
	p := new([]Post)
	sql := "SELECT * FROM posts WHERE published = true"
	err := db.Get(p, sql)
	return p, err
}

// UnpublishedPosts returns all unpublished posts in database
func (db *DB) UnpublishedPosts() (*[]Post, error) {
	p := new([]Post)
	sql := "SELECT * FROM posts WHERE published = false"
	err := db.Get(p, sql)
	return p, err
}

// AllPosts returns all posts in database
func (db *DB) AllPosts() (*[]Post, error) {
	p := new([]Post)
	sql := "SELECT * FROM posts"
	err := db.Get(p, sql)
	return p, err
}

// FindPost returns the post that matches the uuid
func (db *DB) FindPost(id uuid.UUID) (*Post, error) {
	p := new(Post)
	sql := "SELECT * FROM posts WHERE id = $1"
	err := db.Get(p, sql, id)
	return p, err
}

// FindPostsByUser returns a slice of posts created by the given user
func (db *DB) FindPostsByUser(user uuid.UUID) (*[]Post, error) {
	p := new([]Post)
	sql := "SELECT * FROM posts WHERE user_id = $1"
	err := db.Select(p, sql, user)
	return p, err
}

// InsertPost creates a post for the given user and returns the post
func (db *DB) InsertPost(user uuid.UUID, title string, slug string, subtitle string, short string, content string, digest string, published bool) (*Post, error) {
	p := new(Post)
	sql := "INSERT INTO posts (user_id, title, slug, subtitle, short, content, digest, published) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *"
	err := db.Get(p, sql, user, title, slug, subtitle, short, content, digest, published)
	return p, err
}

// UpdatePost updates a post in the database and returns the updated image
func (db *DB) UpdatePost(id uuid.UUID, user uuid.UUID, title string, slug string, subtitle string, short string, content string, digest string, published bool) (*Post, error) {
	p := new(Post)
	sql := "UPDATE posts SET (user_id, title, slug, subtitle, short, content, digest, published) VALUES ($2, $3, $4, $5, $6, $7, $8) WHERE id = $1 RETURNING *"
	err := db.Get(p, sql, id, user, title, slug, subtitle, short, content, digest, published)
	return p, err
}

// DeletePost deletes and returns the image from the database that matches the uuid
func (db *DB) DeletePost(id uuid.UUID, user uuid.UUID) (*Post, error) {
	p := new(Post)
	sql := "DELETE FROM posts WHERE id = $1 AND user_id = $2 RETURNING *"
	err := db.Get(p, sql, id, user)
	return p, err
}

/////////////////////
// Image Functions //
/////////////////////

// AllImages returns all images in the database
func (db *DB) AllImages() (*[]Image, error) {
	i := new([]Image)
	sql := "SELECT * FROM images"
	err := db.Get(i, sql)
	return i, err
}

// FindImage returns an image from the database for the given uuid
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
