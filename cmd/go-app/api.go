package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

func (app *application) mount() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Get("/", app.healthHandler.Check)
	mux.Get("/ws", app.websocketHandler.ServeWebSocket)

	mux.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(jwtauth.Verifier(app.tokenAuth))
				r.Use(jwtauth.Authenticator(app.tokenAuth))

				r.Get("/{id}", app.userHandler.GetById)
			})

			r.Post("/", app.userHandler.Register)
			r.Post("/login", app.userHandler.Login)
		})

		r.Route("/chats", func(r chi.Router) {
			r.Use(jwtauth.Verifier(app.tokenAuth))
			r.Use(jwtauth.Authenticator(app.tokenAuth))

			r.Post("/", app.chatHandler.Create)
			r.Get("/", app.chatHandler.List)
			r.Post("/{chatId}/messages", app.chatHandler.SendMessage)
		})
	})

	return mux
}

func (app *application) newHttpServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    app.config.Address,
		Handler: handler,
	}
}
