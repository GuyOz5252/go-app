package main

import (
	"github.com/GuyOz5252/go-app/cmd/api/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func (app *application) addRoutes(mux *chi.Mux) {
	healthHandler := handlers.NewHealthHandler()
	mux.Get("/", healthHandler.Check)

	userHandler := handlers.NewUserHandler(app.userService, app.tokenAuth)
	mux.Route("/users", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(app.tokenAuth))
			r.Use(jwtauth.Authenticator(app.tokenAuth))

			r.Get("/{id}", userHandler.GetByID)
		})

		r.Post("/", userHandler.Create)
		r.Post("/login", userHandler.Login)
	})
}
