package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func newServer() *http.Server {
	mux := chi.NewRouter()

	// TODO: configure this to my logger
	mux.Use(middleware.Logger)

	mapRoutes(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return server
}
