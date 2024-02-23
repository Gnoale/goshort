package main

import (
	"database/sql"
	"goshort/api"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	os.Remove("./foo.db")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	b, err := os.ReadFile("./store/schema.sql")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(string(b))
	if err != nil {
		panic(err)
	}

	h := api.NewHandler(db)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/shorten", h.Shorten)
		r.Get("/{slug}", h.Slug)

	})
	http.ListenAndServe(":8000", r)
}
