package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sdwalsh/mirango-go/models"
	"golang.org/x/crypto/bcrypt"
)

// UserCustomClaim is the custom claim for user authentication contains a uuid.UUID and jwt.StandardClaims
type UserCustomClaim struct {
	UUID uuid.UUID `json:"uuid"`
	jwt.StandardClaims
}

//
// ContextKey solution to "should not use basic type string as key in context.WithValue"
// Provided by @matryer
// https://medium.com/@matryer/context-keys-in-go-5312346a868d
//

type contextKey string

func (c contextKey) String() string {
	return "mirango-go controller context key" + string(c)
}

var (
	contextSignedIn = contextKey("signed_in")
	contextUser     = contextKey("user")
	contextAdmin    = contextKey("admin")
)

// getUser is a useful function for taking the user claim from the jwt
func (env *Env) getUser(tokenString string, hmacSecret []byte) (uuid.UUID, error) {
	u := *new(uuid.UUID)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
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
		// Context defaults to false for signed in
		ctx := context.WithValue(r.Context(), contextSignedIn, false)
		ctx = context.WithValue(ctx, contextAdmin, false)
		cookie, err := r.Cookie("authentication")
		if err != nil {
			env.log(r, err)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		ID, err := uuid.Parse(cookie.Value)
		if err != nil {
			env.log(r, err)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		user, err := env.DB.GetUserByID(ID)
		if err != nil {
			env.log(r, err)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		// SignedIn is updated to true if it reaches this point and adds the user struct
		ctx = context.WithValue(r.Context(), contextSignedIn, true)
		ctx = context.WithValue(ctx, contextUser, user)
		ctx = context.WithValue(ctx, contextAdmin, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminOnly blocks all requests unless from a user with a role of ADMIN
// assumes user is stored in context (run UserCtx before running this middleware)
func (env *Env) AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, ok := ctx.Value(contextUser).(*models.User)
		// User was not found, invalid jwt, or not signed in
		if !ok {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		// If user is not an admin return unauthorized
		if user.Role != "ADMIN" {
			http.Error(w, http.StatusText(401), 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Login takes a user and password from a login form and checks it against
// the hashed password in the database sets jwt if accepted
func (env *Env) Login(w http.ResponseWriter, r *http.Request) {
	// Clean everything but the password (password is hashed)
	s := bluemonday.UGCPolicy()
	uname := s.Sanitize(r.FormValue("user"))
	password := r.FormValue("password")
	t := time.Now().Add(time.Hour)

	// Database call and bcrypt compare hashed passwords
	u, err := env.DB.GetUserByUname(uname)
	if err != nil {
		env.log(r, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = bcrypt.CompareHashAndPassword(u.Digest, []byte(password+env.Salt))
	if err != nil {
		env.log(r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Generate new jwt with uuid claim and add cookie
	token := jwt.NewWithClaims(jwt.SigningMethodES512, UserCustomClaim{
		u.ID,
		jwt.StandardClaims{
			ExpiresAt: t.Unix(),
		},
	})
	tokenString, err := token.SignedString(env.Hmac)
	cookie := http.Cookie{
		Name:     "authentication",
		Value:    tokenString,
		Expires:  t,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

// Logout will delete the authentication cookie
func (env *Env) Logout(w http.ResponseWriter, r *http.Request) {
	// Setting MaxAge to < 0 will delete cookie now
	cookie := http.Cookie{
		Name:     "authentication",
		MaxAge:   -1,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

// CreateAccount takes user, email, gpg, password, password2 from a form
// cleans any data that might show up on a page and returns 500 if the two passwords do not match
// 409 if there is an error in hashing and 200 if the user is added to the database successfully
func (env *Env) CreateAccount(w http.ResponseWriter, r *http.Request) {
	// Clean everything but the password (password is hashed)
	s := bluemonday.UGCPolicy()
	user := s.Sanitize(r.FormValue("user"))
	email := s.Sanitize(r.FormValue("email"))
	gpg := s.Sanitize(r.FormValue("gpg"))
	password := r.FormValue("password")
	password2 := r.FormValue("password2")

	// Are passwords the same?
	if password != password2 {
		w.WriteHeader(http.StatusUnauthorized)
	}

	// Generate the hash if there is an error in hashing we'll return 500
	digest, err := bcrypt.GenerateFromPassword([]byte(password+env.Salt), 10)
	if err != nil {
		env.log(r, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Ignore the user struct returned and report back 409 if error or 202 if ok
	_, err = env.DB.InsertUser(user, digest, "MEMBER", email, gpg)
	if err != nil {
		env.log(r, err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
