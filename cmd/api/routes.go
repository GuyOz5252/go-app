package main

import (
	"github.com/GuyOz5252/go-app/cmd/api/handlers"
	"github.com/go-chi/chi/v5"
)

func (app *application) addRoutes(mux *chi.Mux) {
	healthHandler := handlers.NewHealthHandler()
	mux.Get("/", healthHandler.Check)

	userHandler := handlers.NewUserHandler(app.userService)
	mux.Route("/users", func(r chi.Router) {
		r.Get("/{id}", userHandler.GetByID)
		r.Post("/", userHandler.Create)
	})
}
