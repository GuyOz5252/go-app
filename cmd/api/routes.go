package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func mapRoutes(app *application, mux *chi.Mux) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("healthy"))
	})

	mux.Route("/api", func(r chi.Router) {
		r.Mount("/users", mountUserRoutes(app))
	})
}
