package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func helloWorld(mux *chi.Mux) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
}
