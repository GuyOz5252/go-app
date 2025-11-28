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

	mapRoutes(app, mux)

	server := &http.Server{
		Addr:    app.config.address,
		Handler: mux,
	}

	return server
}
