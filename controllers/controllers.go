package controllers

import "github.com/sdwalsh/mirango-go/models"

// Env carries database access to controllers
type Env struct {
	DB   models.Datastore
	hmac []byte
}
