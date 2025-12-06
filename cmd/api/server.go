package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) newServer() *http.Server {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	app.addRoutes(mux)

	server := &http.Server{
		Addr:    app.config.Address,
		Handler: mux,
	}

	return server
}
