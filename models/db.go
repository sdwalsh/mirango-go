package models

import (
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
}

// DB ...
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
