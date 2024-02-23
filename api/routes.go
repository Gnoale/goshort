package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"goshort/store"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
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
	_, err := url.ParseRequestURI(body.LongURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid URl"))
		return
	}
	// insert in DB
	// unique constraint on the url column
	id, err := h.q.CreateURL(r.Context(), body.LongURL)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(err.Error()))
		return
	}
	// encode the slug from our database ID
	slug, err := encode(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	slugUrl := url.URL{
		Scheme: "http",
		Host:   r.Host,
		Path:   "api/v1/" + slug,
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(slugUrl.String() + "\n"))
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
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		return
	}
	w.Header().Add("location", url)
	w.WriteHeader(http.StatusMovedPermanently)
}
