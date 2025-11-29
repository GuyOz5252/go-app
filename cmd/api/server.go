package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) newServer() *http.Server {
	mux := chi.NewRouter()

	// TODO: configure this to my logger
	mux.Use(middleware.Logger)

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("healthy"))
	})

	mux.Route("/api", func(r chi.Router) {
		r.Mount("/users", mountUserRoutes(app))
	})

	server := &http.Server{
		Addr:    app.config.Address,
		Handler: mux,
	}

	return server
}
