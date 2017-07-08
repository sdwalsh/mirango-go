package main

import (
	"crypto/rand"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	DatabaseURL string
	Hmac        string
	Salt        string
}

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

	// Create new chi router and add middleware
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Start server and add csrf middleware (32 bit key and chi router)
	http.ListenAndServe(":3333", csrf.Protect(key)(r))
}
