package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func mapRoutes(mux *chi.Mux) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
}
