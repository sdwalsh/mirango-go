package models

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

/*
Database / model access designed after Alex Edward's post on Go database organization
http://www.alexedwards.net/blog/organising-database-access
Code licensed under the MIT License

Copyright 2015 Alex Edwards

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
*/

// Datastore interface contains all of our functions for the PostgreSQL database
// this allows us to mock the database during tests!
type Datastore interface {
	// User Functions
	InsertUser(uname string, digest []byte, role string, email string, gpg string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByID(user uuid.UUID) (*User, error)
	GetUserByUname(uname string) (*User, error)
	// Post Functions
	PublishedPosts(start int, end int) (*[]Post, error)
	UnpublishedPosts() (*[]Post, error)
	GetPosts(start int, end int) (*[]Post, error)
	FindPost(id uuid.UUID) (*Post, error)
	FindPostsByUser(user uuid.UUID) (*[]Post, error)
	InsertPost(user uuid.UUID, title string, slug string, subtitle string, short string, content string, digest string, published bool) (*Post, error)
	UpdatePost(id uuid.UUID, user uuid.UUID, title string, slug string, subtitle string, short string, content string, digest string, published bool) (*Post, error)
	DeletePost(id uuid.UUID, user uuid.UUID) (*Post, error)
	// Image Functions
	AllImages() (*[]Image, error)
	FindImage(id uuid.UUID) (*Image, error)
	FindImagesByUser(user uuid.UUID) (*[]Image, error)
	InsertImage(user uuid.UUID, url string, medium string, small string, caption string) (*Image, error)
	DeleteImage(id uuid.UUID) (*Image, error)
}

// DB holds the database access method (allows us to mock the database)
type DB struct {
	*sqlx.DB
}

// CreateDB initializes the database connection and returns a new DB struct pointer
// databaseName is expected to be a PostgreSQL URL
func CreateDB(database string, connectionOptions string) (*DB, error) {
	db, err := sqlx.Connect("postgres", connectionOptions)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
