package controllers

import "github.com/sdwalsh/mirango-go/models"
import "github.com/gorilla/securecookie"

// Env carries database access to controllers
type Env struct {
	DB   models.Datastore
	S    *securecookie.SecureCookie
	Hmac []byte
	Salt string
}
