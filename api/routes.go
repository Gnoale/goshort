package api

import (
	"database/sql"
	"encoding/json"
	"goshort/store"
	"net/http"
	"net/url"
)

type handler struct {
	q store.Repository
}

func NewHandler(db *sql.DB) *handler {
	return &handler{store.NewRepository(db)}
}

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
	// TODO: manage the error and if the value already exist return http 409 duplicate
	id, err := h.q.CreateURL(ctx, body.LongURL)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slugUrl := url.URL{
		Scheme: r.URL.Scheme,
		Host:   r.URL.Host,
		Path:   encode(id),
	}
	w.Write([]byte(slugUrl.String()))
}

func (h *handler) Slug(w http.ResponseWriter, r *http.Request) {

}
