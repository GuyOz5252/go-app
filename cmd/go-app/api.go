package main

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/GuyOz5252/go-app/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

type application struct {
	config        *Config
	logger        *slog.Logger
	db            *sql.DB
	tokenAuth     *jwtauth.JWTAuth
	healthHandler *handlers.HealthHandler
	userHandler   *handlers.UserHandler
}

func (app *application) mount() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Get("/", app.healthHandler.Check)

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
	})

	return mux
}

func (app *application) newHttpServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    app.config.Address,
		Handler: handler,
	}
}
