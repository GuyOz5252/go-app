package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func mountUserRoutes(app *application) http.Handler {
	mux := chi.NewRouter()

	mux.Post("/", func(w http.ResponseWriter, r *http.Request) {})

	return mux
}
