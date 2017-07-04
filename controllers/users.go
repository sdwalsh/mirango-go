package controllers

import (
	"context"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// Login, Logout, Sign-up

func getUser(tokenString string, hmacSecret []byte) (uuid.UUID, error) {
	u := *new(uuid.UUID)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})
	if err != nil {
		return u, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		u, err = uuid.Parse(claims["uuid"].(string))
		return u, err
	}
	return u, err
}

// UserCtx loads the user into the context if it exists
func (env *Env) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("authentication")
		if err != nil {
			next.ServeHTTP(w, r)
		}
		ID, err := uuid.Parse(cookie.Value)
		if err != nil {
			next.ServeHTTP(w, r)
		}
		user, err := env.DB.GetUserByID(ID)
		if err != nil {
			next.ServeHTTP(w, r)
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (env *Env) Login(w http.ResponseWriter, r *http.Request) {

}

func (env *Env) Logout(w http.ResponseWriter, r *http.Request) {

}

func (env *Env) CreateAccount(w http.ResponseWriter, r *http.Request) {

}
