package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func newServer() *http.Server {
	mux := chi.NewRouter()

	helloWorld(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return server
}
