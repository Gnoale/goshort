package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"goshort/store"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	q store.Repository
}

func NewHandler(db *sql.DB) *handler {
	return &handler{store.NewRepository(db)}
}

// Shorten handle the url shortening for POST request
func (h *handler) Shorten(w http.ResponseWriter, r *http.Request) {
	body := shortBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// validate our url
	_, err := url.Parse(body.LongURL)
	if err != nil {
		w.Write([]byte("invalid URl"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// insert in DB
	ctx := r.Context()
	// unique constraint on the url column
	id, err := h.q.CreateURL(ctx, body.LongURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	// encode the slug from our database ID
	slug, err := encode(id)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	slugUrl := url.URL{
		Scheme: r.URL.Scheme,
		Host:   r.URL.Host,
		Path:   slug,
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(slugUrl.String()))
}

// Slug return and HTTP 301 redirect upon succerssfull slug retrieval from the database to a complete URl
func (h *handler) Slug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	id, err := decode(slug)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// fetch the original url from database
	url, err := h.q.GetURL(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Add("location", url)
	w.WriteHeader(http.StatusMovedPermanently)
}
