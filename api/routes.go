package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"
)

type handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *handler {
	return &handler{db}
}

func (h *handler) Shorten(w http.ResponseWriter, r *http.Request) {

	var body *shortBody
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	parsedUrl, err := url.Parse(body.LongURL)
	if err != nil {
		w.Write([]byte("invalid URl"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func (h *handler) Slug(w http.ResponseWriter, r *http.Request) {

}
