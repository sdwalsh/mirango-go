package controllers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/gorilla/csrf"
	"github.com/gorilla/securecookie"
	"github.com/sdwalsh/mirango-go/models"
)

// Env carries database access to controllers
type Env struct {
	DB    models.Datastore
	S     *securecookie.SecureCookie
	Hmac  []byte
	Salt  string
	Sugar *zap.SugaredLogger
}

// Helper to log any errors
func (env *Env) log(r *http.Request, err error) {
	env.Sugar.Infow("error encountered during controller",
		"request ip:", r.RemoteAddr,
		"header:", r.Header,
		"method:", r.Method,
		"url:", r.URL,
		"error:", err,
	)
}

// Nonspecific routes go here

// Dashboard is a function that wraps calls commonly used on the homepage
func (env *Env) Dashboard(w http.ResponseWriter, r *http.Request) {
	p, err := env.DB.GetPosts(0, 10)
	if err != nil {
		env.log(r, err)
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
	}
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*p)
}
