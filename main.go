package main

import (
	"crypto/rand"
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
	"github.com/gorilla/securecookie"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/sdwalsh/mirango-go/controllers"
	"github.com/sdwalsh/mirango-go/models"
)

// Specification is the struct of all required environmental variables
type Specification struct {
	Database string
	User     string
	Password string
	SSL      string
	Hmac     string
	Salt     string
	Port     string
}

// Main sets up the server configuration and middleware and start the server
func main() {
	// Process environmental configuration
	var c Specification
	err := envconfig.Process("mirango", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("%+v\n", c)

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

	// Database setup and ping
	databaseOptions := "user=" + c.User + " password=" + c.Password + " dbname=" + c.Database + " sslmode=" + c.SSL
	log.Print(databaseOptions)
	post, err := sqlx.Connect("postgres", databaseOptions)
	if err != nil {
		log.Fatal(err.Error())
	}
	data := new(models.DB)
	data.DB = post

	// Zap logger setup
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Pass around Env to routes
	e := controllers.Env{
		DB:    data,
		S:     s,
		Hmac:  []byte(c.Hmac),
		Salt:  c.Salt,
		Sugar: sugar,
	}

	// Create new chi router and add middleware
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Timeout(60 * time.Second))

	// Load the user into context if logged in
	r.Use(e.UserCtx)

	// Homepage
	r.Get("/", e.Dashboard)

	// Post Routes
	r.Get("/posts", e.GetPublishedPosts)
	r.Get("/posts/{postID}", e.GetPost)

	r.Post("/login", e.Login)
	r.Post("/logout", e.Logout)

	// User / Admin Routes

	r.Route("/admin", func(r chi.Router) {
		r.Use(e.UserCtx)
		r.Use(e.AdminOnly)

		r.Get("/posts/unpublished", e.GetUnpublishedPosts)
		r.Post("/posts", e.CreatePost)
		r.Put("/posts/{postID}", e.UpdatePost)
		r.Delete("/posts/{postID}", e.DeletePost)

		r.Post("/users", e.CreateAccount)
	})

	// Start server and add csrf middleware (32 bit key and chi router)
	err = http.ListenAndServe(c.Port, csrf.Protect(key)(r))
	if err != nil {
		log.Fatal("Cannot start server")
	}
}
