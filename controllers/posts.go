package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sdwalsh/mirango-go/models"
)

// CreatePost takes form data and inserts a post into the database
// expects the user and role to be in the request context.
func (env *Env) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Grab the context to get the user
	ctx := r.Context()
	// Clean everything
	s := bluemonday.UGCPolicy()
	user := ctx.Value(contextUser).(*models.User)
	title := s.Sanitize(r.FormValue("title"))
	slug := s.Sanitize(r.FormValue("slug"))
	subtitle := s.Sanitize(r.FormValue("subtitle"))
	short := s.Sanitize(r.FormValue("short"))
	content := s.Sanitize(r.FormValue("content"))
	digest := s.Sanitize(r.FormValue("digest"))
	// published must be parsed into a bool
	published, err := strconv.ParseBool(s.Sanitize(r.FormValue("published")))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	p, err := env.DB.InsertPost(user.ID, title, slug, subtitle, short, content, digest, published)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	// Send out created post
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

// GetPost if ID matches a post return a json post. If the post is unpublished
// check if user is an admin otherwise return 404
func (env *Env) GetPost(w http.ResponseWriter, r *http.Request) {
}

// GetAllPosts is an admin only function that returns all posts published or not
func (env *Env) GetAllPosts(w http.ResponseWriter, r *http.Request) {
}

// GetAllPublishedPosts returns all published / public posts
func (env *Env) GetAllPublishedPosts(w http.ResponseWriter, r *http.Request) {
}

// UpdatePost takes form data and a post ID to update stored information
func (env *Env) UpdatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s := bluemonday.UGCPolicy()
	user := ctx.Value(contextUser).(*models.User)
	// post ID needs to be sanitized and parsed - return if error parsing
	id, err := uuid.Parse(s.Sanitize(r.FormValue("id")))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	title := s.Sanitize(r.FormValue("title"))
	slug := s.Sanitize(r.FormValue("slug"))
	subtitle := s.Sanitize(r.FormValue("subtitle"))
	short := s.Sanitize(r.FormValue("short"))
	content := s.Sanitize(r.FormValue("content"))
	digest := s.Sanitize(r.FormValue("digest"))
	// published must be parsed into a bool
	published, err := strconv.ParseBool(s.Sanitize(r.FormValue("published")))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	p, err := env.DB.UpdatePost(id, user.ID, title, slug, subtitle, short, content, digest, published)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	// Send out updated post
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

// DeletePost removes a post from the database if the user is the owner
func (env *Env) DeletePost(w http.ResponseWriter, r *http.Request) {
	// Grab the context to get the user
	ctx := r.Context()
	user := ctx.Value(contextUser).(*models.User)
	id, err := uuid.Parse(chi.URLParam(r, "postID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	// Ignored returned post since it was deleted
	_, err = env.DB.DeletePost(id, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	// Send out deleted post
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
