package main

import (
	"crypto/rand"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
	"github.com/gorilla/securecookie"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/sdwalsh/mirango-go/controllers"
	"github.com/sdwalsh/mirango-go/models"
)

// Configuration is the struct of all required environmental variables
type Configuration struct {
	Database string
	Hmac     string
	Salt     string
	Port     string
}

// Main sets up the server configuration and middleware and start the server
func main() {
	// Process environmental configuration
	var c Configuration
	err := envconfig.Process("mirango", &c)
	if err != nil {
		log.Fatal("Could not process environmental variables")
	}

	// Generate initial key for gorilla/csrf log.Fatal if key generation fails
	key := make([]byte, 32)
	_, err = rand.Read(key)
	if err != nil {
		log.Fatal("Key for CSRF protection could not be generated.")
	}

	// Setup gorilla/securecookie
	hashKey := securecookie.GenerateRandomKey(64)
	blockKey := securecookie.GenerateRandomKey(32)
	s := securecookie.New(hashKey, blockKey)

	post, err := sqlx.Connect("postgres", c.Database)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	data := new(models.DB)
	data.DB = post

	// Pass around Env to routes
	e := controllers.Env{
		DB:   data,
		S:    s,
		Hmac: []byte(c.Hmac),
		Salt: c.Salt,
	}

	// Create new chi router and add middleware
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Start server and add csrf middleware (32 bit key and chi router)
	http.ListenAndServe(c.Port, csrf.Protect(key)(r))
}
