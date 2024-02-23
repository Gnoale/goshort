package main

import (
	"database/sql"
	"goshort/api"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	h := api.NewHandler(db)

	r.Route("api/v1", func(r chi.Router) {
		r.Post("/shorten", h.Shorten)
		r.Get("/{slug}", h.Slug)

	})
}
